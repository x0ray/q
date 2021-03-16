// Package qs - q scripting language
package qs

import (
	"crypto/rand"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"strings"
	"time"

	"github.com/x0ray/q/ar"
)

// Integer limit values.
const (
	MaxInt = int(^uint(0) >> 1)
	MinInt = int(-MaxInt - 1)
)

var startedAt time.Time

func init() {
	startedAt = time.Now()
}

func getIntField(L *LState, tb *LOAList, key string, v int) int {
	ret := tb.RawGetString(key)
	if ln, ok := ret.(LNumber); ok {
		return int(ln)
	}
	return v
}

func getBoolField(L *LState, tb *LOAList, key string, v bool) bool {
	ret := tb.RawGetString(key)
	if lb, ok := ret.(LBool); ok {
		return bool(lb)
	}
	return v
}

func osChdir(L *LState) int {
	err := os.Chdir(L.CheckString(1))
	if err != nil {
		L.Push(LNil)
		L.Push(LString(err.Error()))
		return 2
	} else {
		L.Push(LTrue)
		return 1
	}
}

func osClearenv(L *LState) int {
	os.Clearenv()
	return 1
}

func osClock(L *LState) int {
	L.Push(LNumber(float64(time.Now().Sub(startedAt)) / float64(time.Second)))
	return 1
}

func osDiffTime(L *LState) int {
	L.Push(LNumber(L.CheckInt64(1) - L.CheckInt64(2)))
	return 1
}

func osExecute(L *LState) int {
	var procAttr os.ProcAttr
	procAttr.Files = []*os.File{os.Stdin, os.Stdout, os.Stderr}
	cmd, args := popenArgs(L.CheckString(1))
	args = append([]string{cmd}, args...)
	process, err := os.StartProcess(cmd, args, &procAttr)
	if err != nil {
		L.Push(LNumber(1))
		return 1
	}

	ps, err := process.Wait()
	if err != nil || !ps.Success() {
		L.Push(LNumber(1))
		return 1
	}
	L.Push(LNumber(0))
	return 1
}

func osExist(L *LState) int {
	fn := L.CheckString(1)
	f, err := os.Open(fn)
	if err == nil {
		f.Close()
		L.Push(LTrue)
		return 1
	}
	L.Push(LFalse)
	return 1
}

// osEmbedded - returns true if the script is running embedded in an application
func osExit(L *LState) int {
	if !QsEmbedded { // exit only if script is NOT embedded
		L.Close()
		os.Exit(L.OptInt(1, 0))
	}
	return 1
}

func osDate(L *LState) int {
	t := time.Now()
	cfmt := "%c"
	if L.GetTop() >= 1 {
		cfmt = L.CheckString(1)
		if strings.HasPrefix(cfmt, "!") {
			t = time.Now().UTC()
			cfmt = strings.TrimLeft(cfmt, "!")
		}
		if L.GetTop() >= 2 {
			t = time.Unix(L.CheckInt64(2), 0)
		}
		if strings.HasPrefix(cfmt, "*t") {
			ret := L.NewOAList()
			ret.RawSetString("year", LNumber(t.Year()))
			ret.RawSetString("month", LNumber(t.Month()))
			ret.RawSetString("day", LNumber(t.Day()))
			ret.RawSetString("hour", LNumber(t.Hour()))
			ret.RawSetString("min", LNumber(t.Minute()))
			ret.RawSetString("sec", LNumber(t.Second()))
			ret.RawSetString("weekday", LNumber(t.Weekday()))
			// TODO yday & dst
			ret.RawSetString("yearday", LNumber(0))
			ret.RawSetString("isdst", LFalse)
			L.Push(ret)
			return 1
		}
	}
	L.Push(LString(strftime(t, cfmt)))
	return 1
}

// osFiles - create list of files into an OA list
func osFiles(L *LState) int {
	baseDir := L.CheckString(1)
	nargs := L.GetTop()
	if nargs != 1 {
		L.RaiseError("wrong number of arguments")
	}
	var err error
	var lst *LOAList
	lst = L.NewOAList()

	err = filepath.Walk(baseDir, func(path string, f os.FileInfo, err error) error {
		if err == nil {
			ls := newLOAList(0, 0)
			lst.RawSetString(path, LValue(ls))
			ls.RawSetString("mode", LString(f.Mode().String()))
			if !f.IsDir() {
				ls.RawSetString("name", LString(f.Name()))
				ls.RawSetString("size", LNumber(f.Size())) // in bytes
				ls.RawSetString("modtime", LString(f.ModTime().String()))
			}
		}
		return nil
	})
	if err != nil {
		L.RaiseError(err.Error())
	}

	L.Push(lst)
	return 1
}

func osGetEnv(L *LState) int {
	v := os.Getenv(L.CheckString(1))
	if len(v) == 0 {
		L.Push(LNil)
	} else {
		L.Push(LString(v))
	}
	return 1
}

func osGeteuid(L *LState) int {
	euid := os.Geteuid()
	L.Push(LNumber(euid))
	return 1
}

func osArgStr(L *LState) int {
	argstr := ""
	argv := ""
	for _, v := range scrArgs {
		argv = v
		if strings.ContainsAny(v, ` "'`) { // need quotes
			if !strings.ContainsAny(v, `\"`) {
				argv = strings.Replace(v, `"`, `\"`, -1)
			}
			if !strings.ContainsAny(v, `\'`) {
				argv = strings.Replace(v, `'`, `\'`, -1)
			}
			argv = `"` + argv + `"`
		}
		if argstr == "" {
			argstr = argv
		} else {
			argstr = argstr + " " + argv
		}
	}
	L.Push(LString(argstr))
	return 1
}

func osArgList(L *LState) int {
	inp := L.CheckString(1)
	args := new(ar.Args)
	err := args.ParseArg(inp)
	if err != nil {
		L.Push(LNil)
		return 1
	}
	sl := args.GetList()
	ret := L.NewOAList()
	for i, v := range sl {
		ret.RawSetInt(i+1, LString(v))
	}
	L.Push(ret)
	return 1
}

func osArgOpts(L *LState) int {
	inp := L.CheckString(1)
	args := new(ar.Args)
	err := args.ParseArg(inp)
	if err != nil {
		L.Push(LNil)
		return 1
	}
	opmap := args.GetMap()
	ret := L.NewOAList()
	for k, v := range opmap {
		ret.RawSetString(k, LString(v))
	}
	L.Push(ret)
	return 1
}

func osGetpid(L *LState) int {
	pid := os.Getpid()
	L.Push(LNumber(pid))
	return 1
}

func osGetppid(L *LState) int {
	ppid := os.Getppid()
	L.Push(LNumber(ppid))
	return 1
}

func osGethome(L *LState) int {
	usr, err := user.Current()
	if err != nil {
		L.Push(LNil)
		L.Push(LString(err.Error()))
		return 2
	} else {
		L.Push(LString(usr.Name))
	}
	return 1
}

func osGetuid(L *LState) int {
	uid := os.Getuid()
	L.Push(LNumber(uid))
	return 1
}

func osGetuser(L *LState) int {
	usr, err := user.Current()
	if err != nil {
		L.Push(LNil)
		L.Push(LString(err.Error()))
		return 2
	} else {
		L.Push(LString(usr.Username))
	}
	return 1
}

func osGetwd(L *LState) int {
	dir, err := os.Getwd()
	if err != nil {
		L.Push(LNil)
		L.Push(LString(err.Error()))
		return 2
	} else {
		L.Push(LString(dir))
	}
	return 1
}

func osHostname(L *LState) int {
	hn, err := os.Hostname()
	if err != nil {
		L.Push(LNil)
		L.Push(LString(err.Error()))
		return 2
	} else {
		L.Push(LString(hn))
	}
	return 1
}

func osRemove(L *LState) int {
	err := os.Remove(L.CheckString(1))
	if err != nil {
		L.Push(LNil)
		L.Push(LString(err.Error()))
		return 2
	} else {
		L.Push(LTrue)
		return 1
	}
}

func osRename(L *LState) int {
	err := os.Rename(L.CheckString(1), L.CheckString(2))
	if err != nil {
		L.Push(LNil)
		L.Push(LString(err.Error()))
		return 2
	} else {
		L.Push(LTrue)
		return 1
	}
}

func osSetLocale(L *LState) int {
	// setlocale is not supported
	L.Push(LFalse)
	return 1
}

func osSetEnv(L *LState) int {
	err := os.Setenv(L.CheckString(1), L.CheckString(2))
	if err != nil {
		L.Push(LNil)
		L.Push(LString(err.Error()))
		return 2
	} else {
		L.Push(LTrue)
		return 1
	}
}

func osSleep(L *LState) int {
	td := L.CheckNumber(1)
	sd := int64(td * 1000000000)
	time.Sleep(time.Duration(sd))
	return 1
}

func osTime(L *LState) int {
	if L.GetTop() == 0 {
		L.Push(LNumber(time.Now().Unix()))
	} else {
		tbl := L.CheckOAList(1)
		sec := getIntField(L, tbl, "sec", 0)
		min := getIntField(L, tbl, "min", 0)
		hour := getIntField(L, tbl, "hour", 12)
		day := getIntField(L, tbl, "day", -1)
		month := getIntField(L, tbl, "month", -1)
		year := getIntField(L, tbl, "year", -1)
		isdst := getBoolField(L, tbl, "isdst", false)
		t := time.Date(year, time.Month(month), day, hour, min, sec, 0, time.Local)
		// TODO dst
		if false {
			print(isdst)
		}
		L.Push(LNumber(t.Unix()))
	}
	return 1
}

func osTmpname(L *LState) int {
	file, err := ioutil.TempFile("", "")
	if err != nil {
		L.RaiseError("unable to generate a unique filename")
	}
	file.Close()
	os.Remove(file.Name()) // ignore errors
	L.Push(LString(file.Name()))
	return 1
}

func osUnsetenv(L *LState) int {
	err := os.Unsetenv(L.CheckString(1))
	if err != nil {
		L.Push(LNil)
		L.Push(LString(err.Error()))
		return 2
	} else {
		L.Push(LTrue)
		return 1
	}
}

func osUuidGen(L *LState) int {
	uuid, err := Uuidgenr()
	if err != nil {
		L.Push(LNil)
		L.Push(LString(err.Error()))
		return 2
	} else {
		L.Push(LString(uuid))
	}
	return 1
}

func osUuidGenFmt(L *LState) int {
	uuid, err := Uuidgenrf()
	if err != nil {
		L.Push(LNil)
		L.Push(LString(err.Error()))
		return 2
	} else {
		L.Push(LString(uuid))
	}
	return 1
}

// osEmbedded - returns true if the script is running embedded in an application
//   or false if the script is running within the q process
func osEmbedded(L *LState) int {
	L.Push(LBool(QsEmbedded))
	return 1
}

//------------------------------------------------------------------------------
// Support functions
//------------------------------------------------------------------------------

// Uuidgenr - generate 16 byte unformatted UUID based on pseudo random number
func Uuidgenr() ([]byte, error) {
	u := make([]byte, 16)
	_, err := rand.Read(u)
	if err == nil {
		// As part of the uuid spec, if you generate a uuid from random it
		// must contain a "4" as the 13th character and a "8", "9", "a",
		// or "b" in the 17th character
		// this makes sure that the 13th character is "4"
		u[6] = (u[6] | 0x40) & 0x4F
		// this make sure that the 17th is "8", "9", "a", or "b"
		u[8] = (u[8] | 0x80) & 0xBF
		return u, nil
	} else {
		return []byte{0}, err
	}
}

// Uuidgenrf - generate 16 byte unformatted UUID based on pseudo random number
func Uuidgenrf() (string, error) {
	u := make([]byte, 16)
	_, err := rand.Read(u)
	if err == nil {
		// As part of the uuid spec, if you generate a uuid from random it
		// must contain a "4" as the 13th character and a "8", "9", "a",
		// or "b" in the 17th character
		// this makes sure that the 13th character is "4"
		u[6] = (u[6] | 0x40) & 0x4F
		// this make sure that the 17th is "8", "9", "a", or "b"
		u[8] = (u[8] | 0x80) & 0xBF
		uuid := fmt.Sprintf("%X-%X-%X-%X-%X", u[0:4], u[4:6], u[6:8], u[8:10], u[10:])
		return uuid, nil
	} else {
		return "", err
	}
}

// fileSystemType converts file system type into a readable string
func fileSystemType(t int64) string {
	fstmap := map[int64]string{
		0xadf5:     "ADFS",      /* Advanced Disc Filing System based on Acorn Winchester Filing System */
		0xadff:     "AFFS",      /* Andrew File System (AFS) is a distributed file system */
		0x62646576: "BDEVFS",    /* Device file system on Unix-like operating systems, used for presenting device files */
		0x42465331: "BEFS",      /* BeOS File System (BeFS) driver */
		0x1badface: "BFS",       /* SCO UnixWare BFS filesystem for Linux */
		0x42494e4d: "BINFMTFS",  /* */
		0x9123683e: "BTRFS",     /* B-Tree filesystem */
		0x27e0eb:   "CGROUP",    /* */
		0xff534d42: "CIFS",      /* Common Internet File System */
		0x73757245: "CODA",      /* Coda distributed file system, experimental, developed by CMU */
		0x012ff7b7: "COH",       /* */
		0x28cd3d45: "CRAMFS",    /* cram a filesystem onto a small ROM */
		0x64626720: "DEBUGFS",   /* debugging virtual file system */
		0x1373:     "DEVFS",     /* */
		0x1cd1:     "DEVPTS",    /* */
		0xde5e81e4: "EFIVARFS",  /* */
		0x00414a53: "EFS",       /* SGI EFS, Extent File System (Irix <0.6) */
		0x137d:     "EXT",       /* Linux extended file system */
		0xef51:     "EXT2",      /* Linux second extended file system, by Rémy Card 1993, 32 TB max */
		0xef53:     "EXT2",      /* */
		0xef54:     "EXT3",      /* Linux third extended file system, allows journaling, by Stephen Tweedie 2001, 32 TB max */
		0xef55:     "EXT4",      /* Linux fourth extended file system, allows journaling (on/off), multiblock allocation, delayed allocation, 1024 PB  max */
		0x65735546: "FUSE",      /* Filesystem in Userspace, non-privileged user file systems runs file system code in user space, FUSE module provides a "bridge" to kernel interfaces */
		0xbad1dea:  "FUTEXFS",   /* Fast-Userspace-muTEX file system, not allowing two threads to access a shared resource at the same time */
		0x4244:     "HFS",       /* Macintosh HFS Filesystem */
		0x00c0ffee: "HOSTFS",    /* */
		0xf995e849: "HPFS",      /* High Performance Filesys (OS/2's HPFS) */
		0x958458f6: "HUGETLBFS", /* huge pages file system, RAM-based virtual/pseudo filesystem */
		0x9660:     "ISOFS",     /* CD/DVD filesystem (ISO-9660 / ECMA-119) */
		0x72b6:     "JFFS2",     /* The Journalling Flash File System, v2 */
		0x3153464a: "JFS",       /*  */
		0x137f:     "MINIX",     /* orig. minix */
		0x138f:     "MINIX1_30", /* 30 char minix */
		0x2468:     "MINIX2",    /* minix V2 */
		0x2478:     "MINIX2_30", /* minix V2, 30 char names */
		0x4d5a:     "MINIX3",    /* minix V3 fs, 60 char names */
		0x19800202: "MQUEUE",    /* message queue file system */
		0x4d44:     "MSDOS",     /* MS-DOS filesystem support */
		0x564c:     "NCP",       /* Netware NCP network protocol */
		0x6969:     "NFS",       /* Networks Filesystem */
		0x3434:     "NILFS",     /* log-structured file system, supports versioning of entire file system, continuous snapshotting, allows user restore of files mistakenly destroyed */
		0x5346544e: "NTFS",      /* NTFS 1.2/3.x driver by Anton Altaparmakov */
		0x7461636f: "OCFS2",     /* OCFS2 1.3.3 */
		0x9fa1:     "OPENPROM",  /* */
		0x50495045: "PIPEFS",    /* PipeFS virtual filesystem, mounted inside kernel, NOT mounted under "/", mounted on "pipe:", making PipeFS its own root */
		0x9fa0:     "PROC",      /* */
		0x6165676c: "PSTOREFS",  /* */
		0x002f:     "QNX4",      /* QNX (OS) Filesystem */
		0x68191122: "QNX6",      /* QNX (OS) Filesystem */
		0x858458f6: "RAMFS",     /* */
		0x52654973: "REISERFS",  /* */
		0x7275:     "ROMFS",     /* ROM filesystem. See genromfs */
		0xf97cff8c: "SELINUX",   /* */
		0x43415d53: "SMACK",     /* Simplified Mandatory Access Control Kernel file system */
		0x517b:     "SMB",       /* Server Message Block, aka Common Internet File System aka Samba */
		0x534f434b: "SOCKFS",    /* pseudofilesystem used by the socket interface to handle file operations on sockets */
		0x73717368: "SQUASHFS",  /* compressed read-only file system */
		0x62656572: "SYSFS",     /* in-memory filesystem, fs hierarchy based on internal organization of kernel data structures */
		0x012ff7b6: "SYSV2",     /* */
		0x012ff7b5: "SYSV4",     /* */
		0x01021994: "TMPFS",     /* temporary file storage file system */
		0x15013346: "UDF",       /* Universal Disk Format Filesystem */
		0x00011954: "UFS",       /* Unix File System aka Berkeley Fast File System, FFS or BSD Fast File System */
		0x9fa2:     "USBDEVICE", /* */
		0x01021997: "V9FS",      /* */
		0xa501fcf5: "VXFS",      /* VERITAS File System, extent-based, primary filesystem of HP-UX */
		0xabba1974: "XENFS",     /* high-performance alternative to NFS */
		0x012ff7b4: "XENIX",     /* XENIX(R) OS filesystem */
		0x58465342: "XFS",       /* high-performance 64-bit journaling file system created by SGI */
		0x012fd16d: "XIAFS"}     /* */

	fs, exists := fstmap[t]
	if exists {
		return fs
	}
	return "Unknown file system"
}
