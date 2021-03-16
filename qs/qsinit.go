// Package qs - q scripting language
package qs

const (
	// BaseLibName is here for consistency; the base procs have no namespace/library.
	BaseLibName = ""
	// LoadLibName is here for consistency; the loading system has no namespace/library.
	LoadLibName = "package"

	// TabLibName is the name of the list Library.
	// TabLibName = "l"

	// IoLibName is the name of the io Library.
	IoLibName = "i"

	// OsLibName is the name of the os Library.
	// OsLibName = "o"

	// StringLibName is the name of the string Library.
	// StringLibName = "s"

	// MathLibName is the name of the math Library.
	// MathLibName = "m"

	// DebugLibName is the name of the debug Library.
	// DebugLibName = "d"

	// ChannelLibName is the name of the channel Library.
	ChannelLibName = "c"

	// CoroutineLibName is the name of the coroutine Library.
	CoroutineLibName = "g"

	// EmiLibName is the name of the EMI Library.
	// EmiLibName = "e"
)

type oaLib struct {
	libName string
	libFunc LGProc
}

var oaLibs = []oaLib{
	oaLib{LoadLibName, OpenPackage},
	oaLib{BaseLibName, OpenBase},
	// oaLib{TabLibName, OpenOAList},
	oaLib{IoLibName, OpenIo},
	// oaLib{OsLibName, OpenOs},
	// oaLib{StringLibName, OpenString},
	// oaLib{MathLibName, OpenMath},
	// oaLib{DebugLibName, OpenDebug},
	oaLib{ChannelLibName, OpenChannel},
	oaLib{CoroutineLibName, OpenCoroutine},
	// oaLib{EmiLibName, OpenEmi},
}

// OpenLibs loads the built-in libraries. It is equivalent to running OpenLoad,
// then OpenBase, then iterating over the other OpenXXX procs in any order.
func (ls *LState) OpenLibs() {
	// NB: Map iteration order in Go is deliberately randomised, so must open Load/Base
	// prior to iterating.
	for _, lib := range oaLibs {
		ls.Push(ls.NewProc(lib.libFunc))
		ls.Push(LString(lib.libName))
		ls.Call(1, 0)
	}
}
