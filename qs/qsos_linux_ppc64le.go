// Package qs - q scripting language
package qs

import (
	"syscall"
	"time"
)

func osStat(L *LState) int {
	fn := L.CheckString(1)
	var b syscall.Stat_t
	err := syscall.Stat(fn, &b)
	if err != nil {
		L.Push(LNil)
		return 1
	}
	// get stat data into list
	var t time.Time
	ret := L.NewOAList()
	ret.RawSetString("name", LString(fn))
	ret.RawSetString("dev", LNumber(b.Dev))
	ret.RawSetString("ino", LNumber(b.Ino))
	ret.RawSetString("nlink", LNumber(b.Nlink))
	ret.RawSetString("mode", LNumber(b.Mode))
	ret.RawSetString("uid", LNumber(b.Uid))
	ret.RawSetString("gid", LNumber(b.Gid))
	ret.RawSetString("rdev", LNumber(b.Rdev))
	ret.RawSetString("size", LNumber(b.Size))
	ret.RawSetString("blksize", LNumber(b.Blksize))
	ret.RawSetString("blocks", LNumber(b.Blocks))
	t = time.Unix(b.Atim.Sec, b.Atim.Nsec)
	ret.RawSetString("atim", LString(strftime(t, "%c")))
	t = time.Unix(b.Mtim.Sec, b.Mtim.Nsec)
	ret.RawSetString("mtim", LString(strftime(t, "%c")))
	t = time.Unix(b.Ctim.Sec, b.Ctim.Nsec)
	ret.RawSetString("ctim", LString(strftime(t, "%c")))
	var mode string
	switch b.Mode & syscall.S_IFMT {
	case syscall.S_IFREG:
		mode = "Regular file"
	case syscall.S_IFDIR:
		mode = "Directory"
	case syscall.S_IFCHR:
		mode = "Character device"
	case syscall.S_IFBLK:
		mode = "Block device"
	case syscall.S_IFLNK:
		mode = "Symbolic link"
	case syscall.S_IFIFO:
		mode = "FIFO or pipe"
	case syscall.S_IFSOCK:
		mode = "Socket"
	default:
		mode = "Unknown"
	}
	ret.RawSetString("type", LString(mode))
	L.Push(ret)
	return 1
}

func osStatfs(L *LState) int {
	path := L.CheckString(1)
	var b syscall.Statfs_t
	err := syscall.Statfs(path, &b)
	if err != nil {
		L.Push(LNil)
		return 1
	}
	// get statfs data into list
	ret := L.NewOAList()
	ret.RawSetString("name", LString(path))
	ret.RawSetString("type", LNumber(b.Type))
	ret.RawSetString("bsize", LNumber(b.Bsize))
	ret.RawSetString("blocks", LNumber(b.Blocks))
	ret.RawSetString("bfree", LNumber(b.Bfree))
	ret.RawSetString("bavail", LNumber(b.Bavail))
	ret.RawSetString("files", LNumber(b.Files))
	ret.RawSetString("ffree", LNumber(b.Ffree))
	ret.RawSetString("fsid1", LNumber(b.Fsid.X__val[0]))
	ret.RawSetString("fsid2", LNumber(b.Fsid.X__val[1]))
	ret.RawSetString("namelen", LNumber(b.Namelen))
	ret.RawSetString("frsize", LNumber(b.Frsize))
	ret.RawSetString("typename", LString(fileSystemType(b.Type)))

	// add calculated fields
	var (
		percentFreeBytes float64
		percentFreeNodes float64
		percentUsedBytes float64
		percentUsedNodes float64
	)
	usedBlocks := b.Blocks - b.Bavail // b.Bavail was b.Bfree
	ret.RawSetString("usedblocks", LNumber(usedBlocks))
	totalBytes := b.Blocks * uint64(b.Bsize)
	ret.RawSetString("totalbytes", LNumber(totalBytes))
	freeBytes := b.Bfree * uint64(b.Bsize)
	ret.RawSetString("freebytes", LNumber(freeBytes))
	usedBytes := totalBytes - freeBytes
	ret.RawSetString("usedbytes", LNumber(usedBytes))

	if totalBytes != 0 {
		percentFreeBytes = (float64(freeBytes) * 100) / float64(totalBytes)
	} else {
		percentFreeBytes = 0
	}
	ret.RawSetString("percentfreebytes", LNumber(percentFreeBytes))
	if b.Files != 0 {
		percentFreeNodes = (float64(b.Ffree) * 100) / float64(b.Files)
	} else {
		percentFreeNodes = 0
	}
	ret.RawSetString("percentfreenodes", LNumber(percentFreeNodes))
	percentUsedBytes = 100.0 - percentFreeBytes
	percentUsedNodes = 100.0 - percentFreeNodes
	ret.RawSetString("percentusedbytes", LNumber(percentUsedBytes))
	ret.RawSetString("percentusednodes", LNumber(percentUsedNodes))

	L.Push(ret)
	return 1
}
