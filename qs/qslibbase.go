// Package qs - q scripting language
package qs

import (
	"flag"
	"fmt"
	"io"
	"math"
	"os"
	"runtime"
	"strconv"
	"strings"
)

func OpenBase(L *LState) int {
	global := L.Get(GlobalsIndex).(*LOAList)
	L.SetGlobal("_G", global)
	L.SetGlobal("_VERSION", LString(PGM+" "+VER))
	basemod := L.RegisterModule("_G", baseFuncs)
	global.RawSetString("ipairs", L.NewClosure(baseIpairs, L.NewProc(ipairsaux)))
	global.RawSetString("pairs", L.NewClosure(basePairs, L.NewProc(pairsaux)))

	// add gmatch for strings
	var mod *LOAList
	mod = basemod.(*LOAList)
	gmatch := L.NewClosure(strGmatch, L.NewProc(strGmatchIter))
	mod.RawSetString("gmatch", gmatch)
	mod.RawSetString("gfind", gmatch)

	// add constants for system
	mod.RawSetString("os", LString(runtime.GOOS))
	mod.RawSetString("arch", LString(runtime.GOARCH))

	// add constants for math
	mod.RawSetString("pi", LNumber(math.Pi))
	mod.RawSetString("e", LNumber(math.E))
	mod.RawSetString("phi", LNumber(math.Phi))
	mod.RawSetString("sqrt2", LNumber(math.Sqrt2))
	mod.RawSetString("sqrte", LNumber(math.SqrtE))
	mod.RawSetString("huge", LNumber(math.MaxFloat64))
	mod.RawSetString("small", LNumber(math.SmallestNonzeroFloat64))

	L.Push(basemod)
	return 1
}

var baseFuncs = map[string]LGProc{
	// base procs
	"assert":         baseAssert,
	"bye":            baseTerm,
	"collectgarbage": baseCollectGarbage,
	"error":          baseError,
	"getfenv":        baseGetFEnv,
	"getmetalist":    baseGetMetalist,
	"help":           baseHelp,
	"keys":           baseKeys,
	"load":           baseLoad,
	"loadfile":       baseLoadFile,
	"loadstring":     baseLoadString,
	"logd":           baseLogDebug,
	"loge":           baseLogError,
	"logi":           baseLogInfo,
	"logw":           baseLogWarning,
	"next":           baseNext,
	"pcall":          basePCall,
	"_printregs":     base_PrintRegs,
	"put":            basePut,
	"quit":           baseTerm,
	"rawequal":       baseRawEqual,
	"rawget":         baseRawGet,
	"rawset":         baseRawSet,
	"run":            baseRun,
	"select":         baseSelect,
	"setfenv":        baseSetFEnv,
	"setmetalist":    baseSetMetalist,
	"stop":           baseTerm,
	"tonumber":       baseToNumber,
	"tostring":       baseToString,
	"type":           baseType,
	"unpack":         baseUnpack,
	"xpcall":         baseXPCall,
	// loadlib
	"module":  loModule,
	"require": loRequire,
	// string procs
	"after":        strAfter,
	"before":       strBefore,
	"byte":         strByte,
	"char":         strChar,
	"contains":     strContains,
	"containsany":  strContainsAny,
	"count":        strCount,
	"decodebase64": strDecodeBase64,
	"dump":         strDump,
	"encodebase64": strEncodeBase64,
	"escapexml":    strEscapeXmlData,
	"find":         strFind,
	"format":       strFormat,
	"gsub":         strGsub,
	"hasprefix":    strHasPrefix,
	"hassuffix":    strHasSuffix,
	"index":        strIndex,
	"indexany":     strIndexAny,
	"isname":       strIsName,
	"isxmltagname": strIsXmlTagName,
	"lastindex":    strLastIndex,
	"lastindexany": strLastIndexAny,
	"length":       strLen,
	"len":          strLen,
	"lower":        strLower,
	"makexmltag":   strMakeXmlTagName,
	"match":        strMatch,
	"prxmatch":     strPrxMatch,
	"prxchange":    strPrxChange,
	"rep":          strRep,
	"replace":      strReplace,
	"reverse":      strReverse,
	"scan":         strScan,
	"scanall":      strScanAll,
	"sub":          strSub,
	"substr":       strSubstr,
	"trim":         strTrim,
	"trimleft":     strTrimLeft,
	"trimprefix":   strTrimPrefix,
	"trimright":    strTrimRight,
	"trimspace":    strTrimSpace,
	"trimsuffix":   strTrimSuffix,
	"title":        strTitle,
	"unescapexml":  strUnEscapeXmlData,
	"upper":        strUpper,
	// math procs
	"abs":        mathAbs,
	"acos":       mathAcos,
	"asin":       mathAsin,
	"atan":       mathAtan,
	"atan2":      mathAtan2,
	"ceil":       mathCeil,
	"cos":        mathCos,
	"cosh":       mathCosh,
	"deg":        mathDeg,
	"exp":        mathExp,
	"fact":       mathFact,
	"fib":        mathFib,
	"floor":      mathFloor,
	"fmod":       mathFmod,
	"frexp":      mathFrexp,
	"ldexp":      mathLdexp,
	"log":        mathLog,
	"log10":      mathLog10,
	"max":        mathMax,
	"mean":       mathMean,
	"median":     mathMedian,
	"min":        mathMin,
	"mod":        mathMod,
	"mode":       mathMode,
	"modf":       mathModf,
	"pow":        mathPow,
	"rad":        mathRad,
	"random":     mathRandom,
	"randomseed": mathRandomseed,
	"range":      mathRange,
	"rms":        mathRms,
	"sin":        mathSin,
	"sinh":       mathSinh,
	"sqrt":       mathSqrt,
	"stddev":     mathStdDev,
	"sum":        mathSum,
	"tan":        mathTan,
	"tanh":       mathTanh,
	"variance":   mathVariance,
	// os procs
	"argstr":     osArgStr,
	"arglist":    osArgList,
	"argopts":    osArgOpts,
	"chdir":      osChdir,
	"clearenv":   osClearenv,
	"clock":      osClock,
	"date":       osDate,
	"difftime":   osDiffTime,
	"embedded":   osEmbedded,
	"execute":    osExecute,
	"exist":      osExist,
	"exit":       osExit,
	"files":      osFiles,
	"getenv":     osGetEnv,
	"geteuid":    osGeteuid,
	"getpid":     osGetpid,
	"getppid":    osGetppid,
	"gethome":    osGethome,
	"getuid":     osGetuid,
	"getuser":    osGetuser,
	"getwd":      osGetwd,
	"hostname":   osHostname,
	"remove":     osRemove,
	"rename":     osRename,
	"setenv":     osSetEnv,
	"setlocale":  osSetLocale,
	"sleep":      osSleep,
	"stat":       osStat,
	"statfs":     osStatfs,
	"time":       osTime,
	"tmpname":    osTmpname,
	"unsetenv":   osUnsetenv,
	"uuidgen":    osUuidGen,
	"uuidgenfmt": osUuidGenFmt,
	// debug procs
	"dbggetfenv":     debugGetFEnv,
	"dbggetinfo":     debugGetInfo,
	"dbggetlocal":    debugGetLocal,
	"dbggetmetalist": debugGetMetalist,
	"dbggetupvalue":  debugGetUpvalue,
	"dbgsetfenv":     debugSetFEnv,
	"dbgsetlocal":    debugSetLocal,
	"dbgsetmetalist": debugSetMetalist,
	"dbgsetupvalue":  debugSetUpvalue,
	"dbgtraceback":   debugTraceback,
	// list procs
	"dumpl":      listDump,
	"getn":       listGetN,
	"concat":     listConcat,
	"insert":     listInsert,
	"maxn":       listMaxN,
	"erase":      listErase,
	"marshal":    listMarshal,
	"marshalxml": listMarshalXml,
	"unmarshal":  listUnMarshal,
	"sort":       listSort,
}

func baseAssert(L *LState) int {
	if !L.ToBool(1) {
		L.RaiseError(L.OptString(2, "assertion failed!"))
		return 0
	}
	return L.GetTop()
}

func baseCollectGarbage(L *LState) int {
	runtime.GC()
	return 0
}

func baseRun(L *LState) int {
	src := L.ToString(1)
	top := L.GetTop()
	fn, err := L.LoadFile(src)
	if err != nil {
		L.Push(LString(err.Error()))
		L.Panic(L)
	}
	L.Push(fn)
	L.Call(0, MultRet)
	return L.GetTop() - top
}

func baseError(L *LState) int {
	obj := L.CheckAny(1)
	level := L.OptInt(2, 1)
	L.Error(obj, level)
	return 0
}

func baseGetFEnv(L *LState) int {
	var value LValue
	if L.GetTop() == 0 {
		value = LNumber(1)
	} else {
		value = L.Get(1)
	}

	if fn, ok := value.(*LProc); ok {
		if !fn.IsG {
			L.Push(fn.Env)
		} else {
			L.Push(L.G.Global)
		}
		return 1
	}

	if number, ok := value.(LNumber); ok {
		level := int(float64(number))
		if level <= 0 {
			L.Push(L.Env)
		} else {
			cf := L.currentFrame
			for i := 0; i < level && cf != nil; i++ {
				cf = cf.Parent
			}
			if cf == nil || cf.Fn.IsG {
				L.Push(L.G.Global)
			} else {
				L.Push(cf.Fn.Env)
			}
		}
		return 1
	}

	L.Push(L.G.Global)
	return 1
}

func baseGetMetalist(L *LState) int {
	L.Push(L.GetMetalist(L.CheckAny(1)))
	return 1
}

func baseHelp(L *LState) int {
	fmt.Printf("Help for %s Ver %s\n\n %s", PGM, VER, HelpText) // print main command blurb
	flag.VisitAll(PrtHelpFlag)
	fmt.Printf("%s\n", HelpTextVerbose)
	return 0
}

func ipairsaux(L *LState) int {
	lst := L.CheckOAList(1)
	i := L.CheckInt(2)
	i++
	v := lst.RawGetInt(i)
	if v == LNil {
		return 0
	} else {
		L.Pop(1)
		L.Push(LNumber(i))
		L.Push(LNumber(i))
		L.Push(v)
		return 2
	}
}

func baseIpairs(L *LState) int {
	lst := L.CheckOAList(1)
	L.Push(L.Get(UpvalueIndex(1)))
	L.Push(lst)
	L.Push(LNumber(0))
	return 3
}

func baseKeys(L *LState) int {
	lst := L.CheckOAList(1)
	if lst.strdict == nil {
		L.Push(LNumber(0))
		return 1
	}
	ln := len(lst.strdict)
	L.Push(LNumber(ln))
	return 1
}

func loadaux(L *LState, reader io.Reader, segmentname string) int {
	if fn, err := L.Load(reader, segmentname); err != nil {
		L.Push(LNil)
		L.Push(LString(err.Error()))
		return 2
	} else {
		L.Push(fn)
		return 1
	}
}

func baseLoad(L *LState) int {
	fn := L.CheckProc(1)
	segmentname := L.OptString(2, "?")
	top := L.GetTop()
	buf := []string{}
	for {
		L.SetTop(top)
		L.Push(fn)
		L.Call(0, 1)
		ret := L.reg.Pop()
		if ret == LNil {
			break
		} else if LVCanConvToString(ret) {
			str := ret.String()
			if len(str) > 0 {
				buf = append(buf, string(str))
			} else {
				break
			}
		} else {
			L.Push(LNil)
			L.Push(LString("reader proc must return a string"))
			return 2
		}
	}
	return loadaux(L, strings.NewReader(strings.Join(buf, "")), segmentname)
}

func baseLoadFile(L *LState) int {
	var reader io.Reader
	var segmentname string
	var err error
	if L.GetTop() < 1 {
		reader = os.Stdin
		segmentname = "<stdin>"
	} else {
		segmentname = L.CheckString(1)
		reader, err = os.Open(segmentname)
		if err != nil {
			L.Push(LNil)
			L.Push(LString(fmt.Sprintf("can not open file: %v", segmentname)))
			return 2
		}
		defer reader.(*os.File).Close()
	}
	return loadaux(L, reader, segmentname)
}

func baseLoadString(L *LState) int {
	return loadaux(L, strings.NewReader(L.CheckString(1)), L.OptString(2, "<string>"))
}

func baseLogInfo(L *LState) int {
	top := L.GetTop()
	s := ""
	for i := 1; i <= top; i++ {
		s = s + L.ToStringMeta(L.Get(i)).String()
	}
	log.Info().Msgf("%s", s)
	return 0
}

func baseLogWarning(L *LState) int {
	top := L.GetTop()
	s := ""
	for i := 1; i <= top; i++ {
		s = s + L.ToStringMeta(L.Get(i)).String()
	}
	log.Warn().Msgf("%s", s)
	return 0
}

func baseLogError(L *LState) int {
	top := L.GetTop()
	s := ""
	for i := 1; i <= top; i++ {
		s = s + L.ToStringMeta(L.Get(i)).String()
	}
	log.Error().Msgf("%s", s)
	return 0
}

func baseLogDebug(L *LState) int {
	top := L.GetTop()
	s := ""
	for i := 1; i <= top; i++ {
		s = s + L.ToStringMeta(L.Get(i)).String()
	}
	if debug {
		log.Debug().Msgf("%s", s)
	}
	return 0
}

func baseNext(L *LState) int {
	lst := L.CheckOAList(1)
	index := LNil
	if L.GetTop() >= 2 {
		index = L.Get(2)
	}
	key, value := lst.Next(index)
	if key == LNil {
		L.Push(LNil)
		return 1
	}
	L.Push(key)
	L.Push(value)
	return 2
}

func pairsaux(L *LState) int {
	lst := L.CheckOAList(1)
	key, value := lst.Next(L.Get(2))
	if key == LNil {
		return 0
	} else {
		L.Pop(1)
		L.Push(key)
		L.Push(key)
		L.Push(value)
		return 2
	}
}

func basePairs(L *LState) int {
	lst := L.CheckOAList(1)
	L.Push(L.Get(UpvalueIndex(1)))
	L.Push(lst)
	L.Push(LNil)
	return 3
}

func basePCall(L *LState) int {
	L.CheckProc(1)
	nargs := L.GetTop() - 1
	if err := L.PCall(nargs, MultRet, nil); err != nil {
		L.Push(LFalse)
		if aerr, ok := err.(*ApiError); ok {
			L.Push(aerr.Object)
		} else {
			L.Push(LString(err.Error()))
		}
		return 2
	} else {
		L.Insert(LTrue, 1)
		return L.GetTop()
	}
}

func base_PrintRegs(L *LState) int {
	L.printReg()
	return 0
}

func basePut(L *LState) int {
	top := L.GetTop()
	for i := 1; i <= top; i++ {
		fmt.Print(L.ToStringMeta(L.Get(i)).String())
		if i != top {
			fmt.Print("\t")
		}
	}
	fmt.Println("")
	return 0
}

func baseRawEqual(L *LState) int {
	if L.CheckAny(1) == L.CheckAny(2) {
		L.Push(LTrue)
	} else {
		L.Push(LFalse)
	}
	return 1
}

func baseRawGet(L *LState) int {
	L.Push(L.RawGet(L.CheckOAList(1), L.CheckAny(2)))
	return 1
}

func baseRawSet(L *LState) int {
	L.RawSet(L.CheckOAList(1), L.CheckAny(2), L.CheckAny(3))
	return 0
}

func baseSelect(L *LState) int {
	L.CheckTypes(1, LTNumber, LTString)
	switch lv := L.Get(1).(type) {
	case LNumber:
		idx := int(lv)
		num := L.reg.Top() - L.indexToReg(int(lv)) - 1
		if idx < 0 {
			num++
		}
		return num
	case LString:
		if string(lv) != "#" {
			L.ArgError(1, "invalid string '"+string(lv)+"'")
		}
		L.Push(LNumber(L.GetTop() - 1))
		return 1
	}
	return 0
}

func baseSetFEnv(L *LState) int {
	var value LValue
	if L.GetTop() == 0 {
		value = LNumber(1)
	} else {
		value = L.Get(1)
	}
	env := L.CheckOAList(2)

	if fn, ok := value.(*LProc); ok {
		if fn.IsG {
			L.RaiseError("can not change the environment of given object")
		} else {
			fn.Env = env
			L.Push(fn)
			return 1
		}
	}

	if number, ok := value.(LNumber); ok {
		level := int(float64(number))
		if level <= 0 {
			L.Env = env
			return 0
		}

		cf := L.currentFrame
		for i := 0; i < level && cf != nil; i++ {
			cf = cf.Parent
		}
		if cf == nil || cf.Fn.IsG {
			L.RaiseError("can not change the environment of given object")
		} else {
			cf.Fn.Env = env
			L.Push(cf.Fn)
			return 1
		}
	}

	L.RaiseError("can not change the environment of given object")
	return 0
}

func baseSetMetalist(L *LState) int {
	L.CheckTypes(2, LTNil, LTOAList)
	obj := L.Get(1)
	if obj == LNil {
		L.RaiseError("can not set metalist to a nil object.")
	}
	mt := L.Get(2)
	if m := L.metalist(obj, true); m != LNil {
		if lst, ok := m.(*LOAList); ok && lst.RawGetString("__metalist") != LNil {
			L.RaiseError("can not change a protected metalist")
		}
	}
	L.SetMetalist(obj, mt)
	L.SetTop(1)
	return 1
}

func baseTerm(L *LState) int {
	L.Close()
	os.Exit(L.OptInt(1, 0))
	return 1
}

func baseToNumber(L *LState) int {
	base := L.OptInt(2, 10)
	switch lv := L.CheckAny(1).(type) {
	case LNumber:
		L.Push(lv)
	case LString:
		str := strings.Trim(string(lv), " \n\t")
		if strings.Index(str, ".") > -1 {
			if v, err := strconv.ParseFloat(str, LNumberBit); err != nil {
				L.Push(LNil)
			} else {
				L.Push(LNumber(v))
			}
		} else {
			if v, err := strconv.ParseInt(str, base, LNumberBit); err != nil {
				L.Push(LNil)
			} else {
				L.Push(LNumber(v))
			}
		}
	default:
		L.Push(LNil)
	}
	return 1
}

func baseToString(L *LState) int {
	v1 := L.CheckAny(1)
	L.Push(L.ToStringMeta(v1))
	return 1
}

func baseType(L *LState) int {
	L.Push(LString(L.CheckAny(1).Type().String()))
	return 1
}

func baseUnpack(L *LState) int {
	lst := L.CheckOAList(1)
	start := L.OptInt(2, 1)
	end := L.OptInt(3, lst.Len())
	for i := start; i <= end; i++ {
		L.Push(lst.RawGetInt(i))
	}
	ret := end - start + 1
	if ret < 0 {
		return 0
	}
	return ret
}

func baseXPCall(L *LState) int {
	fn := L.CheckProc(1)
	errfunc := L.CheckProc(2)

	top := L.GetTop()
	L.Push(fn)
	if err := L.PCall(0, MultRet, errfunc); err != nil {
		L.Push(LFalse)
		if aerr, ok := err.(*ApiError); ok {
			L.Push(aerr.Object)
		} else {
			L.Push(LString(err.Error()))
		}
		return 2
	} else {
		L.Insert(LTrue, top+1)
		return L.GetTop() - top
	}
}

func loModule(L *LState) int {
	name := L.CheckString(1)
	loaded := L.GetField(L.Get(RegistryIndex), "_LOADED")
	lst := L.GetField(loaded, name)
	if _, ok := lst.(*LOAList); !ok {
		lst = L.FindOAList(L.Get(GlobalsIndex).(*LOAList), name, 1)
		if lst == LNil {
			L.RaiseError("name conflict for module: %v", name)
		}
		L.SetField(loaded, name, lst)
	}
	if L.GetField(lst, "_NAME") == LNil {
		L.SetField(lst, "_M", lst)
		L.SetField(lst, "_NAME", LString(name))
		names := strings.Split(name, ".")
		pname := ""
		if len(names) > 1 {
			pname = strings.Join(names[:len(names)-1], ".") + "."
		}
		L.SetField(lst, "_PACKAGE", LString(pname))
	}

	caller := L.currentFrame.Parent
	if caller == nil {
		L.RaiseError("no calling stack.")
	} else if caller.Fn.IsG {
		L.RaiseError("module() can not be called from GProcs.")
	}
	L.SetFEnv(caller.Fn, lst)

	top := L.GetTop()
	for i := 2; i <= top; i++ {
		L.Push(L.Get(i))
		L.Push(lst)
		L.Call(1, 0)
	}
	L.Push(lst)
	return 1
}

var loopdetection = &LUserData{}

func loRequire(L *LState) int {
	name := L.CheckString(1)
	loaded := L.GetField(L.Get(RegistryIndex), "_LOADED")
	lv := L.GetField(loaded, name)
	if LVAsBool(lv) {
		if lv == loopdetection {
			L.RaiseError("loop or previous error loading module: %s", name)
		}
		L.Push(lv)
		return 1
	}
	loaders, ok := L.GetField(L.Get(RegistryIndex), "_LOADERS").(*LOAList)
	if !ok {
		L.RaiseError("package.loaders must be a list")
	}
	messages := []string{}
	var modasfunc LValue
	for i := 1; ; i++ {
		loader := L.RawGetInt(loaders, i)
		if loader == LNil {
			L.RaiseError("module %s not found:\n\t%s, ", name, strings.Join(messages, "\n\t"))
		}
		L.Push(loader)
		L.Push(LString(name))
		L.Call(1, 1)
		ret := L.reg.Pop()
		switch retv := ret.(type) {
		case *LProc:
			modasfunc = retv
			goto loopbreak
		case LString:
			messages = append(messages, string(retv))
		}
	}
loopbreak:
	L.SetField(loaded, name, loopdetection)
	L.Push(modasfunc)
	L.Push(LString(name))
	L.Call(1, 1)
	ret := L.reg.Pop()
	modv := L.GetField(loaded, name)
	if ret != LNil && modv == loopdetection {
		L.SetField(loaded, name, ret)
		L.Push(ret)
	} else if modv == loopdetection {
		L.SetField(loaded, name, LTrue)
		L.Push(LTrue)
	} else {
		L.Push(modv)
	}
	return 1
}
