// Package qs - q scripting language
package qs

import (
	"fmt"
	"strings"
)

const (
	VarArgHasArg   uint8 = 1
	VarArgIsVarArg uint8 = 2
	VarArgNeedsArg uint8 = 4
)

type DbgLocalInfo struct {
	Name    string
	StartPc int
	EndPc   int
}

type DbgCall struct {
	Name string
	Pc   int
}

type ProcProto struct {
	SourceName       string
	LineDefined      int
	LastLineDefined  int
	NumUpvalues      uint8
	NumParameters    uint8
	IsVarArg         uint8
	NumUsedRegisters uint8
	Code             []uint32
	Constants        []LValue
	ProcPrototypes   []*ProcProto

	DbgSourcePositions []int
	DbgLocals          []*DbgLocalInfo
	DbgCalls           []DbgCall
	DbgUpvalues        []string

	stringConstants []string
}

type Upvalue struct {
	next   *Upvalue
	reg    *registry
	index  int
	value  LValue
	closed bool
}

func (uv *Upvalue) Value() LValue {
	if uv.closed || uv.reg == nil {
		return uv.value
	}
	return uv.reg.array[uv.index]
}

func (uv *Upvalue) SetValue(value LValue) {
	if uv.IsClosed() {
		uv.value = value
	} else {
		uv.reg.Set(uv.index, value)
	}
}

func (uv *Upvalue) Close() {
	value := uv.Value()
	uv.closed = true
	uv.value = value
}

func (uv *Upvalue) IsClosed() bool {
	return uv.closed || uv.reg == nil
}

func UpvalueIndex(i int) int {
	return GlobalsIndex - i
}

func newProcProto(name string) *ProcProto {
	return &ProcProto{
		SourceName:       name,
		LineDefined:      0,
		LastLineDefined:  0,
		NumUpvalues:      0,
		NumParameters:    0,
		IsVarArg:         0,
		NumUsedRegisters: 2,
		Code:             make([]uint32, 0, 128),
		Constants:        make([]LValue, 0, 32),
		ProcPrototypes:   make([]*ProcProto, 0, 16),

		DbgSourcePositions: make([]int, 0, 128),
		DbgLocals:          make([]*DbgLocalInfo, 0, 16),
		DbgCalls:           make([]DbgCall, 0, 128),
		DbgUpvalues:        make([]string, 0, 16),

		stringConstants: make([]string, 0, 32),
	}
}

func (fp *ProcProto) String() string {
	return fp.str(1, 0)
}

func (fp *ProcProto) str(level int, count int) string {
	indent := strings.Repeat("  ", level-1)
	buf := []string{}
	buf = append(buf, fmt.Sprintf("| %v proc [%v] defn: %v\n",
		indent, count, level))
	buf = append(buf, fmt.Sprintf("| %v u-values:%v parms:%v stacks %v \n",
		indent, fp.NumUpvalues, fp.NumParameters, fp.NumUsedRegisters))
	for reg, linfo := range fp.DbgLocals {
		buf = append(buf, fmt.Sprintf("| %v.local %v : %v\n", indent, linfo.Name, reg))
	}
	for reg, upvalue := range fp.DbgUpvalues {
		buf = append(buf, fmt.Sprintf("| %v.u-value %v : %v\n", indent, upvalue, reg))
	}
	for reg, conzt := range fp.Constants {
		buf = append(buf, fmt.Sprintf("| %v.const %v : %v\n", indent, conzt.String(), reg))
	}
	buf = append(buf, "|\n")

	protono := 0
	for no, code := range fp.Code {
		inst := opGetOpCode(code)
		if inst == OP_CLOSURE {
			buf = append(buf, "|\n")
			buf = append(buf, fp.ProcPrototypes[protono].str(level+1, protono))
			buf = append(buf, "|\n")
			protono++
		}
		buf = append(buf, fmt.Sprintf("| %v[%03d] %v (line:%v)\n",
			indent, no+1, opToString(code), fp.DbgSourcePositions[no]))

	}
	buf = append(buf, fmt.Sprintf("| %v end proc\n", indent))
	return strings.Join(buf, "")
}

func newLProcL(proto *ProcProto, env *LOAList, nupvalue int) *LProc {
	return &LProc{
		IsG: false,
		Env: env,

		Proto:    proto,
		GProc:    nil,
		Upvalues: make([]*Upvalue, nupvalue),
	}
}

func newLProcG(gfunc LGProc, env *LOAList, nupvalue int) *LProc {
	return &LProc{
		IsG: true,
		Env: env,

		Proto:    nil,
		GProc:    gfunc,
		Upvalues: make([]*Upvalue, nupvalue),
	}
}

func (fn *LProc) LocalName(regno, pc int) (string, bool) {
	if fn.IsG {
		return "", false
	}
	p := fn.Proto
	for i := 0; i < len(p.DbgLocals) && p.DbgLocals[i].StartPc < pc; i++ {
		if pc < p.DbgLocals[i].EndPc {
			regno--
			if regno == 0 {
				return p.DbgLocals[i].Name, true
			}
		}
	}
	return "", false
}
