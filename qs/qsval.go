// Package qs - q scripting language
package qs

import (
	"fmt"
	"os"
)

type LValueType int

const (
	LTNil LValueType = iota
	LTBool
	LTNumber
	LTString
	LTProc
	LTUserData
	LTThread
	LTOAList
	LTChannel
)

var lValueNames = [9]string{"nil", "bool", "num", "str", "proc", "data", "thread", "list", "chan"}

func (vt LValueType) String() string {
	return lValueNames[int(vt)]
}

type LValue interface {
	String() string
	Type() LValueType
	// to reduce `runtime.assertI2T2` costs, this method should be used instead of the type assertion in heavy paths(typically inside the VM).
	assertFloat64() (float64, bool)
	// to reduce `runtime.assertI2T2` costs, this method should be used instead of the type assertion in heavy paths(typically inside the VM).
	assertString() (string, bool)
	// to reduce `runtime.assertI2T2` costs, this method should be used instead of the type assertion in heavy paths(typically inside the VM).
	assertProc() (*LProc, bool)
}

// LVIsFalse returns true if a given LValue is a nil or false otherwise false.
func LVIsFalse(v LValue) bool { return v == LNil || v == LFalse }

// LVIsFalse returns false if a given LValue is a nil or false otherwise true.
func LVAsBool(v LValue) bool { return v != LNil && v != LFalse }

// LVAsString returns string representation of a given LValue
// if the LValue is a string or number, otherwise an empty string.
func LVAsString(v LValue) string {
	switch sn := v.(type) {
	case LString, LNumber:
		return sn.String()
	default:
		return ""
	}
}

// LVCanConvToString returns true if a given LValue is a string or number
// otherwise false.
func LVCanConvToString(v LValue) bool {
	switch v.(type) {
	case LString, LNumber:
		return true
	default:
		return false
	}
}

// LVAsNumber tries to convert a given LValue to a number.
func LVAsNumber(v LValue) LNumber {
	switch lv := v.(type) {
	case LNumber:
		return lv
	case LString:
		if num, err := parseNumber(string(lv)); err == nil {
			return num
		}
	}
	return LNumber(0)
}

type LNilType struct{}

func (nl *LNilType) String() string                 { return "nil" }
func (nl *LNilType) Type() LValueType               { return LTNil }
func (nl *LNilType) assertFloat64() (float64, bool) { return 0, false }
func (nl *LNilType) assertString() (string, bool)   { return "", false }
func (nl *LNilType) assertProc() (*LProc, bool)     { return nil, false }

var LNil = LValue(&LNilType{})

type LBool bool

func (bl LBool) String() string {
	if bool(bl) {
		return "true"
	}
	return "false"
}
func (bl LBool) Type() LValueType               { return LTBool }
func (bl LBool) assertFloat64() (float64, bool) { return 0, false }
func (bl LBool) assertString() (string, bool)   { return "", false }
func (bl LBool) assertProc() (*LProc, bool)     { return nil, false }

var LTrue = LBool(true)
var LFalse = LBool(false)

type LString string

func (st LString) String() string                 { return string(st) }
func (st LString) Type() LValueType               { return LTString }
func (st LString) assertFloat64() (float64, bool) { return 0, false }
func (st LString) assertString() (string, bool)   { return string(st), true }
func (st LString) assertProc() (*LProc, bool)     { return nil, false }

// fmt.Formatter interface
func (st LString) Format(f fmt.State, c rune) {
	switch c {
	case 'd', 'i':
		if nm, err := parseNumber(string(st)); err != nil {
			defaultFormat(nm, f, 'd')
		} else {
			defaultFormat(string(st), f, 's')
		}
	default:
		defaultFormat(string(st), f, c)
	}
}

func (nm LNumber) String() string {
	if isInteger(nm) {
		return fmt.Sprint(int64(nm))
	}
	return fmt.Sprint(float64(nm))
}

func (nm LNumber) Type() LValueType               { return LTNumber }
func (nm LNumber) assertFloat64() (float64, bool) { return float64(nm), true }
func (nm LNumber) assertString() (string, bool)   { return "", false }
func (nm LNumber) assertProc() (*LProc, bool)     { return nil, false }

// fmt.Formatter interface
func (nm LNumber) Format(f fmt.State, c rune) {
	switch c {
	case 'q', 's':
		defaultFormat(nm.String(), f, c)
	case 'b', 'c', 'd', 'o', 'x', 'X', 'U':
		defaultFormat(int64(nm), f, c)
	case 'e', 'E', 'f', 'F', 'g', 'G':
		defaultFormat(float64(nm), f, c)
	case 'i':
		defaultFormat(int64(nm), f, 'd')
	default:
		if isInteger(nm) {
			defaultFormat(int64(nm), f, c)
		} else {
			defaultFormat(float64(nm), f, c)
		}
	}
}

type LOAList struct {
	Metalist LValue

	array   []LValue
	dict    map[LValue]LValue
	strdict map[string]LValue
	keys    []LValue
	k2i     map[LValue]int
}

func (lst *LOAList) String() string                 { return fmt.Sprintf("list: %p", lst) }
func (lst *LOAList) Type() LValueType               { return LTOAList }
func (lst *LOAList) assertFloat64() (float64, bool) { return 0, false }
func (lst *LOAList) assertString() (string, bool)   { return "", false }
func (lst *LOAList) assertProc() (*LProc, bool)     { return nil, false }

type LProc struct {
	IsG      bool
	Env      *LOAList
	Proto    *ProcProto
	GProc    LGProc
	Upvalues []*Upvalue
}
type LGProc func(*LState) int

func (fn *LProc) String() string                 { return fmt.Sprintf("proc: %p", fn) }
func (fn *LProc) Type() LValueType               { return LTProc }
func (fn *LProc) assertFloat64() (float64, bool) { return 0, false }
func (fn *LProc) assertString() (string, bool)   { return "", false }
func (fn *LProc) assertProc() (*LProc, bool)     { return fn, true }

type Global struct {
	MainThread    *LState
	CurrentThread *LState
	Registry      *LOAList
	Global        *LOAList

	builtinMts map[int]LValue
	tempFiles  []*os.File
	gccount    int32
}

type LState struct {
	G       *Global
	Parent  *LState
	Env     *LOAList
	Panic   func(*LState)
	Dead    bool
	Options Options

	stop         int32
	reg          *registry
	stack        *callFrameStack
	alloc        *allocator
	currentFrame *callFrame
	wrapped      bool
	uvcache      *Upvalue
	hasErrorFunc bool
}

func (ls *LState) String() string                 { return fmt.Sprintf("thread: %p", ls) }
func (ls *LState) Type() LValueType               { return LTThread }
func (ls *LState) assertFloat64() (float64, bool) { return 0, false }
func (ls *LState) assertString() (string, bool)   { return "", false }
func (ls *LState) assertProc() (*LProc, bool)     { return nil, false }

type LUserData struct {
	Value    interface{}
	Env      *LOAList
	Metalist LValue
}

func (ud *LUserData) String() string                 { return fmt.Sprintf("userdata: %p", ud) }
func (ud *LUserData) Type() LValueType               { return LTUserData }
func (ud *LUserData) assertFloat64() (float64, bool) { return 0, false }
func (ud *LUserData) assertString() (string, bool)   { return "", false }
func (ud *LUserData) assertProc() (*LProc, bool)     { return nil, false }

type LChannel chan LValue

func (ch LChannel) String() string                 { return fmt.Sprintf("channel: %p", ch) }
func (ch LChannel) Type() LValueType               { return LTChannel }
func (ch LChannel) assertFloat64() (float64, bool) { return 0, false }
func (ch LChannel) assertString() (string, bool)   { return "", false }
func (ch LChannel) assertProc() (*LProc, bool)     { return nil, false }
