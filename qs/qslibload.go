// Package qs - q scripting language
package qs

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

var loLoaders = []LGProc{loLoaderPreload, loLoaderOa}

func loGetPath(env string, defpath string) string {
	path := os.Getenv(env)
	if len(path) == 0 {
		path = defpath
	}
	path = strings.Replace(path, ";;", ";"+defpath+";", -1)
	if os.PathSeparator != '/' {
		dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
		if err != nil {
			panic(err)
		}
		path = strings.Replace(path, "!", dir, -1)
	}
	return path
}

func loFindFile(L *LState, name, pname string) (string, string) {
	name = strings.Replace(name, ".", string(os.PathSeparator), -1)
	lv := L.GetField(L.GetField(L.Get(EnvironIndex), "package"), pname)
	path, ok := lv.(LString)
	if !ok {
		L.RaiseError("package.%s must be a string", pname)
	}
	messages := []string{}
	for _, pattern := range strings.Split(string(path), ";") {
		oapath := strings.Replace(pattern, "?", name, -1)
		if _, err := os.Stat(oapath); err == nil {
			return oapath, ""
		} else {
			messages = append(messages, err.Error())
		}
	}
	return "", strings.Join(messages, "\n\t")
}

func OpenPackage(L *LState) int {
	packagemod := L.RegisterModule(LoadLibName, loFuncs)

	L.SetField(packagemod, "preload", L.NewOAList())

	loaders := L.CreateOAList(len(loLoaders), 0)
	for i, loader := range loLoaders {
		L.RawSetInt(loaders, i+1, L.NewProc(loader))
	}
	L.SetField(packagemod, "loaders", loaders)
	L.SetField(L.Get(RegistryIndex), "_LOADERS", loaders)

	loaded := L.NewOAList()
	L.SetField(packagemod, "loaded", loaded)
	L.SetField(L.Get(RegistryIndex), "_LOADED", loaded)

	L.SetField(packagemod, "path", LString(loGetPath(QsPath, QsPathDefault)))
	L.SetField(packagemod, "cpath", LString(""))

	L.Push(packagemod)
	return 1
}

var loFuncs = map[string]LGProc{
	"loadlib": loLoadLib,
	"seeall":  loSeeAll,
}

func loLoaderPreload(L *LState) int {
	name := L.CheckString(1)
	preload := L.GetField(L.GetField(L.Get(EnvironIndex), "package"), "preload")
	if _, ok := preload.(*LOAList); !ok {
		L.RaiseError("package.preload must be a list")
	}
	lv := L.GetField(preload, name)
	if lv == LNil {
		L.Push(LString(fmt.Sprintf("no field package.preload['%s']", name)))
		return 1
	}
	L.Push(lv)
	return 1
}

func loLoaderOa(L *LState) int {
	name := L.CheckString(1)
	path, msg := loFindFile(L, name, "path")
	if len(path) == 0 {
		L.Push(LString(msg))
		return 1
	}
	fn, err1 := L.LoadFile(path)
	if err1 != nil {
		L.RaiseError(err1.Error())
	}
	L.Push(fn)
	return 1
}

func loLoadLib(L *LState) int {
	L.RaiseError("loadlib is not supported")
	return 0
}

func loSeeAll(L *LState) int {
	mod := L.CheckOAList(1)
	mt := L.GetMetalist(mod)
	if mt == LNil {
		mt = L.CreateOAList(0, 1)
		L.SetMetalist(mod, mt)
	}
	L.SetField(mt, "__index", L.Get(GlobalsIndex))
	return 0
}
