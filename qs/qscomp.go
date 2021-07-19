// Package qs - q scripting language
package qs

import (
	"fmt"
	"math"
	"os"
	"reflect"

	"github.com/x0ray/q/qs/qsa"
)

const maxRegisters = 200

type expContextType int

const (
	ecGlobal expContextType = iota
	ecUpvalue
	ecLocal
	ecOAList
	ecVararg
	ecMethod
	ecNone
)

const regNotDefined = opMaxArgsA + 1
const labelNoJump = 0

type expcontext struct {
	ctype expContextType
	reg   int
	// varargopt >= 0: wants varargopt+1 results, i.e  a = func()
	// varargopt = -1: ignore results             i.e  func()
	// varargopt = -2: receive all results        i.e  a = {func()}
	varargopt int
}

type assigncontext struct {
	ec       *expcontext
	keyrk    int
	valuerk  int
	keyks    bool
	needmove bool
}

type lblabels struct {
	t int
	f int
	e int
	b bool
}

type constLValueExpr struct {
	qsa.ExprBase

	Value LValue
}

var _ecnone0 = &expcontext{ecNone, regNotDefined, 0}
var _ecnonem1 = &expcontext{ecNone, regNotDefined, -1}
var _ecnonem2 = &expcontext{ecNone, regNotDefined, -2}
var ecfuncdef = &expcontext{ecMethod, regNotDefined, 0}

func ecupdate(ec *expcontext, ctype expContextType, reg, varargopt int) {
	ec.ctype = ctype
	ec.reg = reg
	ec.varargopt = varargopt
}

func ecnone(varargopt int) *expcontext {
	switch varargopt {
	case 0:
		return _ecnone0
	case -1:
		return _ecnonem1
	case -2:
		return _ecnonem2
	}
	return &expcontext{ecNone, regNotDefined, varargopt}
}

func sline(pos qsa.PositionHolder) int {
	return pos.Line()
}

func eline(pos qsa.PositionHolder) int {
	return pos.LastLine()
}

func savereg(ec *expcontext, reg int) int {
	if ec.ctype != ecLocal || ec.reg == regNotDefined {
		return reg
	}
	return ec.reg
}

func raiseCompileError(context *funcContext, line int, format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	panic(&CompileError{Context: context, Line: line, Message: msg})
}

func isVarArgReturnExpr(expr qsa.Expr) bool {
	switch ex := expr.(type) {
	case *qsa.FuncCallExpr:
		return !ex.AdjustRet
	case *qsa.Comma3Expr:
		return true
	}
	return false
}

func lnumberValue(expr qsa.Expr) (LNumber, bool) {
	if ex, ok := expr.(*qsa.NumberExpr); ok {
		lv, err := parseNumber(ex.Value)
		if err != nil {
			lv = LNumber(math.NaN())
		}
		return lv, true
	} else if ex, ok := expr.(*constLValueExpr); ok {
		return ex.Value.(LNumber), true
	}
	return 0, false
}

type CompileError struct {
	Context *funcContext
	Line    int
	Message string
}

func (e *CompileError) Error() string {
	return fmt.Sprintf("line (%v) compile error for %v - %v", e.Line, e.Context.Proto.SourceName, e.Message)
}

type codeStore struct {
	codes []uint32
	lines []int
	pc    int
}

func (cd *codeStore) Add(inst uint32, line int) {
	if l := len(cd.codes); l <= 0 || cd.pc == l {
		cd.codes = append(cd.codes, inst)
		cd.lines = append(cd.lines, line)
	} else {
		cd.codes[cd.pc] = inst
		cd.lines[cd.pc] = line
	}
	cd.pc++
}

func (cd *codeStore) AddABC(op int, a int, b int, c int, line int) {
	cd.Add(opCreateABC(op, a, b, c), line)
}

func (cd *codeStore) AddABx(op int, a int, bx int, line int) {
	cd.Add(opCreateABx(op, a, bx), line)
}

func (cd *codeStore) AddASbx(op int, a int, sbx int, line int) {
	cd.Add(opCreateASbx(op, a, sbx), line)
}

func (cd *codeStore) PropagateKMV(top int, save *int, reg *int, inc int) {
	lastinst := cd.Last()
	if opGetArgA(lastinst) >= top {
		switch opGetOpCode(lastinst) {
		case OP_LOADK:
			cindex := opGetArgBx(lastinst)
			if cindex <= opMaxIndexRk {
				cd.Pop()
				*save = opRkAsk(cindex)
				return
			}
		case OP_MOVE:
			cd.Pop()
			*save = opGetArgB(lastinst)
			return
		}
	}
	*save = *reg
	*reg = *reg + inc
}

func (cd *codeStore) PropagateMV(top int, save *int, reg *int, inc int) {
	lastinst := cd.Last()
	if opGetArgA(lastinst) >= top {
		switch opGetOpCode(lastinst) {
		case OP_MOVE:
			cd.Pop()
			*save = opGetArgB(lastinst)
			return
		}
	}
	*save = *reg
	*reg = *reg + inc
}

func (cd *codeStore) SetOpCode(pc int, v int) {
	opSetOpCode(&cd.codes[pc], v)
}

func (cd *codeStore) SetA(pc int, v int) {
	opSetArgA(&cd.codes[pc], v)
}

func (cd *codeStore) SetB(pc int, v int) {
	opSetArgB(&cd.codes[pc], v)
}

func (cd *codeStore) SetC(pc int, v int) {
	opSetArgC(&cd.codes[pc], v)
}

func (cd *codeStore) SetBx(pc int, v int) {
	opSetArgBx(&cd.codes[pc], v)
}

func (cd *codeStore) SetSbx(pc int, v int) {
	opSetArgSbx(&cd.codes[pc], v)
}

func (cd *codeStore) At(pc int) uint32 {
	return cd.codes[pc]
}

func (cd *codeStore) List() []uint32 {
	return cd.codes[:cd.pc]
}

func (cd *codeStore) PosList() []int {
	return cd.lines[:cd.pc]
}

func (cd *codeStore) LastPC() int {
	return cd.pc - 1
}

func (cd *codeStore) Last() uint32 {
	if cd.pc == 0 {
		return opInvalidInstruction
	}
	return cd.codes[cd.pc-1]
}

func (cd *codeStore) Pop() {
	cd.pc--
}

type varNamePoolValue struct {
	Index int
	Name  string
}

type varNamePool struct {
	names  []string
	offset int
}

func newVarNamePool(offset int) *varNamePool {
	return &varNamePool{make([]string, 0, 16), offset}
}

func (vp *varNamePool) Names() []string {
	return vp.names
}

func (vp *varNamePool) List() []varNamePoolValue {
	result := make([]varNamePoolValue, len(vp.names), len(vp.names))
	for i, name := range vp.names {
		result[i].Index = i + vp.offset
		result[i].Name = name
	}
	return result
}

func (vp *varNamePool) LastIndex() int {
	return vp.offset + len(vp.names)
}

func (vp *varNamePool) Find(name string) int {
	for i := len(vp.names) - 1; i >= 0; i-- {
		if vp.names[i] == name {
			return i + vp.offset
		}
	}
	return -1
}

func (vp *varNamePool) RegisterUnique(name string) int {
	index := vp.Find(name)
	if index < 0 {
		return vp.Register(name)
	}
	return index
}

func (vp *varNamePool) Register(name string) int {
	vp.names = append(vp.names, name)
	return len(vp.names) - 1 + vp.offset
}

type codeBlock struct {
	LocalVars  *varNamePool
	BreakLabel int
	Parent     *codeBlock
	RefUpvalue bool
	LineStart  int
	LastLine   int
}

func newCodeBlock(localvars *varNamePool, blabel int, parent *codeBlock, pos qsa.PositionHolder) *codeBlock {
	bl := &codeBlock{localvars, blabel, parent, false, 0, 0}
	if pos != nil {
		bl.LineStart = pos.Line()
		bl.LastLine = pos.LastLine()
	}
	return bl
}

type funcContext struct {
	Proto    *ProcProto
	Code     *codeStore
	Parent   *funcContext
	Upvalues *varNamePool
	Block    *codeBlock
	Blocks   []*codeBlock
	regTop   int
	labelId  int
	labelPc  map[int]int
}

func newFuncContext(sourcename string, parent *funcContext) *funcContext {
	fc := &funcContext{
		Proto:    newProcProto(sourcename),
		Code:     &codeStore{make([]uint32, 0, 1024), make([]int, 0, 1024), 0},
		Parent:   parent,
		Upvalues: newVarNamePool(0),
		Block:    newCodeBlock(newVarNamePool(0), labelNoJump, nil, nil),
		regTop:   0,
		labelId:  1,
		labelPc:  map[int]int{},
	}
	fc.Blocks = []*codeBlock{fc.Block}
	return fc
}

func (fc *funcContext) NewLabel() int {
	ret := fc.labelId
	fc.labelId++
	return ret
}

func (fc *funcContext) SetLabelPc(label int, pc int) {
	fc.labelPc[label] = pc
}

func (fc *funcContext) GetLabelPc(label int) int {
	return fc.labelPc[label]
}

func (fc *funcContext) ConstIndex(value LValue) int {
	ctype := value.Type()
	for i, lv := range fc.Proto.Constants {
		if lv.Type() == ctype && lv == value {
			return i
		}
	}
	fc.Proto.Constants = append(fc.Proto.Constants, value)
	v := len(fc.Proto.Constants) - 1
	if v > opMaxArgBx {
		raiseCompileError(fc, fc.Proto.LineDefined, "number of constants exceeds limit")
	}
	return v
}

func (fc *funcContext) RegisterLocalVar(name string) int {
	ret := fc.Block.LocalVars.Register(name)
	fc.Proto.DbgLocals = append(fc.Proto.DbgLocals, &DbgLocalInfo{Name: name, StartPc: fc.Code.LastPC() + 1})
	fc.SetRegTop(fc.RegTop() + 1)
	return ret
}

func (fc *funcContext) FindLocalVarAndBlock(name string) (int, *codeBlock) {
	for block := fc.Block; block != nil; block = block.Parent {
		if index := block.LocalVars.Find(name); index > -1 {
			return index, block
		}
	}
	return -1, nil
}

func (fc *funcContext) FindLocalVar(name string) int {
	idx, _ := fc.FindLocalVarAndBlock(name)
	return idx
}

func (fc *funcContext) LocalVars() []varNamePoolValue {
	result := make([]varNamePoolValue, 0, 32)
	for _, block := range fc.Blocks {
		result = append(result, block.LocalVars.List()...)
	}
	return result
}

func (fc *funcContext) EnterBlock(blabel int, pos qsa.PositionHolder) {
	fc.Block = newCodeBlock(newVarNamePool(fc.RegTop()), blabel, fc.Block, pos)
	fc.Blocks = append(fc.Blocks, fc.Block)
}

func (fc *funcContext) CloseUpvalues() int {
	n := -1
	if fc.Block.RefUpvalue {
		n = fc.Block.Parent.LocalVars.LastIndex()
		fc.Code.AddABC(OP_CLOSE, n, 0, 0, fc.Block.LastLine)
	}
	return n
}

func (fc *funcContext) LeaveBlock() int {
	closed := fc.CloseUpvalues()
	fc.EndScope()
	fc.Block = fc.Block.Parent
	fc.SetRegTop(fc.Block.LocalVars.LastIndex())
	return closed
}

func (fc *funcContext) EndScope() {
	for _, vr := range fc.Block.LocalVars.List() {
		fc.Proto.DbgLocals[vr.Index].EndPc = fc.Code.LastPC()
	}
}

func (fc *funcContext) SetRegTop(top int) {
	if top > maxRegisters {
		raiseCompileError(fc, fc.Proto.LineDefined, "number of local variables too large")
	}
	fc.regTop = top
}

func (fc *funcContext) RegTop() int {
	return fc.regTop
}

func compileSeg(context *funcContext, segment []qsa.Stmt) {
	for _, stmt := range segment {
		compileStmt(context, stmt)
	}
}

func compileBlock(context *funcContext, segment []qsa.Stmt) {
	if len(segment) == 0 {
		return
	}
	ph := &qsa.Node{}
	ph.SetLine(sline(segment[0]))
	ph.SetLastLine(eline(segment[len(segment)-1]))
	context.EnterBlock(labelNoJump, ph)
	for _, stmt := range segment {
		compileStmt(context, stmt)
	}
	context.LeaveBlock()
}

func compileStmt(context *funcContext, stmt qsa.Stmt) {
	switch st := stmt.(type) {
	case *qsa.AssignStmt:
		compileAssignStmt(context, st)
	case *qsa.LocalAssignStmt:
		compileLocalAssignStmt(context, st)
	case *qsa.FuncCallStmt:
		compileFuncCallExpr(context, context.RegTop(), st.Expr.(*qsa.FuncCallExpr), ecnone(-1))
	case *qsa.DoBlockStmt:
		context.EnterBlock(labelNoJump, st)
		compileSeg(context, st.Stmts)
		context.LeaveBlock()
	case *qsa.WhileStmt:
		compileWhileStmt(context, st)
	case *qsa.RepeatStmt:
		compileRepeatStmt(context, st)
	case *qsa.FuncDefStmt:
		compileFuncDefStmt(context, st)
	case *qsa.ReturnStmt:
		compileReturnStmt(context, st)
	case *qsa.IfStmt:
		compileIfStmt(context, st)
	case *qsa.BreakStmt:
		compileBreakStmt(context, st)
	case *qsa.NumberForStmt:
		compileNumberForStmt(context, st)
	case *qsa.GenericForStmt:
		compileGenericForStmt(context, st)
	}
}

func compileAssignStmtLeft(context *funcContext, stmt *qsa.AssignStmt) (int, []*assigncontext) {
	reg := context.RegTop()
	acs := make([]*assigncontext, 0, len(stmt.Lhs))
	for i, lhs := range stmt.Lhs {
		islast := i == len(stmt.Lhs)-1
		switch st := lhs.(type) {
		case *qsa.IdentExpr:
			identtype := getIdentRefType(context, context, st)
			ec := &expcontext{identtype, regNotDefined, 0}
			switch identtype {
			case ecGlobal:
				context.ConstIndex(LString(st.Value))
			case ecUpvalue:
				context.Upvalues.RegisterUnique(st.Value)
			case ecLocal:
				if islast {
					ec.reg = context.FindLocalVar(st.Value)
				}
			}
			acs = append(acs, &assigncontext{ec, 0, 0, false, false})
		case *qsa.AttrGetExpr:
			ac := &assigncontext{&expcontext{ecOAList, regNotDefined, 0}, 0, 0, false, false}
			compileExprWithKMVPropagation(context, st.Object, &reg, &ac.ec.reg)
			compileExprWithKMVPropagation(context, st.Key, &reg, &ac.keyrk)
			if _, ok := st.Key.(*qsa.StringExpr); ok {
				ac.keyks = true
			}
			acs = append(acs, ac)

		default:
			log.Error().Msgf("Left side expression [%#v] invalid", st)
			os.Exit(RCERROR)
		}
	}
	return reg, acs
}

func compileAssignStmtRight(context *funcContext, stmt *qsa.AssignStmt, reg int, acs []*assigncontext) (int, []*assigncontext) {
	lennames := len(stmt.Lhs)
	lenexprs := len(stmt.Rhs)
	names_assigned := 0

	for names_assigned < lennames {
		ac := acs[names_assigned]
		ec := ac.ec
		var expr qsa.Expr = nil
		if names_assigned >= lenexprs {
			expr = &qsa.NilExpr{}
			expr.SetLine(sline(stmt.Lhs[names_assigned]))
			expr.SetLastLine(eline(stmt.Lhs[names_assigned]))
		} else if isVarArgReturnExpr(stmt.Rhs[names_assigned]) && (lenexprs-names_assigned-1) <= 0 {
			varargopt := lennames - names_assigned - 1
			regstart := reg
			reginc := compileExpr(context, reg, stmt.Rhs[names_assigned], ecnone(varargopt))
			reg += reginc
			for i := names_assigned; i < names_assigned+int(reginc); i++ {
				acs[i].needmove = true
				if acs[i].ec.ctype == ecOAList {
					acs[i].valuerk = regstart + (i - names_assigned)
				}
			}
			names_assigned = lennames
			continue
		}

		if expr == nil {
			expr = stmt.Rhs[names_assigned]
		}
		idx := reg
		reginc := compileExpr(context, reg, expr, ec)
		if ec.ctype == ecOAList {
			if _, ok := expr.(*qsa.LogicalOpExpr); !ok {
				context.Code.PropagateKMV(context.RegTop(), &ac.valuerk, &reg, reginc)
			} else {
				ac.valuerk = idx
				reg += reginc
			}
		} else {
			ac.needmove = reginc != 0
			reg += reginc
		}
		names_assigned += 1
	}

	rightreg := reg - 1

	// extra right exprs
	for i := names_assigned; i < lenexprs; i++ {
		varargopt := -1
		if i != lenexprs-1 {
			varargopt = 0
		}
		reg += compileExpr(context, reg, stmt.Rhs[i], ecnone(varargopt))
	}
	return rightreg, acs
}

func compileAssignStmt(context *funcContext, stmt *qsa.AssignStmt) {
	code := context.Code
	lennames := len(stmt.Lhs)
	reg, acs := compileAssignStmtLeft(context, stmt)
	reg, acs = compileAssignStmtRight(context, stmt, reg, acs)

	for i := lennames - 1; i >= 0; i-- {
		ex := stmt.Lhs[i]
		switch acs[i].ec.ctype {
		case ecLocal:
			if acs[i].needmove {
				code.AddABC(OP_MOVE, context.FindLocalVar(ex.(*qsa.IdentExpr).Value), reg, 0, sline(ex))
				reg -= 1
			}
		case ecGlobal:
			code.AddABx(OP_SETGLOBAL, reg, context.ConstIndex(LString(ex.(*qsa.IdentExpr).Value)), sline(ex))
			reg -= 1
		case ecUpvalue:
			code.AddABC(OP_SETUPVAL, reg, context.Upvalues.RegisterUnique(ex.(*qsa.IdentExpr).Value), 0, sline(ex))
			reg -= 1
		case ecOAList:
			opcode := OP_SETTABLE
			if acs[i].keyks {
				opcode = OP_SETTABLEKS
			}
			code.AddABC(opcode, acs[i].ec.reg, acs[i].keyrk, acs[i].valuerk, sline(ex))
			if !opIsK(acs[i].valuerk) {
				reg -= 1
			}
		}
	}
}

func compileRegAssignment(context *funcContext, names []string, exprs []qsa.Expr, reg int, nvars int, line int) {
	lennames := len(names)
	lenexprs := len(exprs)
	names_assigned := 0
	ec := &expcontext{}

	for names_assigned < lennames && names_assigned < lenexprs {
		if isVarArgReturnExpr(exprs[names_assigned]) && (lenexprs-names_assigned-1) <= 0 {

			varargopt := nvars - names_assigned
			ecupdate(ec, ecVararg, reg, varargopt-1)
			compileExpr(context, reg, exprs[names_assigned], ec)
			reg += varargopt
			names_assigned = lennames
		} else {
			ecupdate(ec, ecLocal, reg, 0)
			compileExpr(context, reg, exprs[names_assigned], ec)
			reg += 1
			names_assigned += 1
		}
	}

	// extra left names
	if lennames > names_assigned {
		restleft := lennames - names_assigned - 1
		context.Code.AddABC(OP_LOADNIL, reg, reg+restleft, 0, line)
		reg += restleft
	}

	// extra right exprs
	for i := names_assigned; i < lenexprs; i++ {
		varargopt := -1
		if i != lenexprs-1 {
			varargopt = 0
		}
		ecupdate(ec, ecNone, reg, varargopt)
		reg += compileExpr(context, reg, exprs[i], ec)
	}
}

func compileLocalAssignStmt(context *funcContext, stmt *qsa.LocalAssignStmt) {
	reg := context.RegTop()
	if len(stmt.Names) == 1 && len(stmt.Exprs) == 1 {
		if _, ok := stmt.Exprs[0].(*qsa.ProcExpr); ok {
			context.RegisterLocalVar(stmt.Names[0])
			compileRegAssignment(context, stmt.Names, stmt.Exprs, reg, len(stmt.Names), sline(stmt))
			return
		}
	}

	compileRegAssignment(context, stmt.Names, stmt.Exprs, reg, len(stmt.Names), sline(stmt))
	for _, name := range stmt.Names {
		context.RegisterLocalVar(name)
	}
}

func compileReturnStmt(context *funcContext, stmt *qsa.ReturnStmt) {
	lenexprs := len(stmt.Exprs)
	code := context.Code
	reg := context.RegTop()
	a := reg
	lastisvaarg := false

	if lenexprs == 1 {
		switch ex := stmt.Exprs[0].(type) {
		case *qsa.IdentExpr:
			if idx := context.FindLocalVar(ex.Value); idx > -1 {
				code.AddABC(OP_RETURN, idx, 2, 0, sline(stmt))
				return
			}
		case *qsa.FuncCallExpr:
			reg += compileExpr(context, reg, ex, ecnone(-2))
			code.SetOpCode(code.LastPC(), OP_TAILCALL)
			code.AddABC(OP_RETURN, a, 0, 0, sline(stmt))
			return
		}
	}

	for i, expr := range stmt.Exprs {
		if i == lenexprs-1 && isVarArgReturnExpr(expr) {
			compileExpr(context, reg, expr, ecnone(-2))
			lastisvaarg = true
		} else {
			reg += compileExpr(context, reg, expr, ecnone(0))
		}
	}
	count := reg - a + 1
	if lastisvaarg {
		count = 0
	}
	context.Code.AddABC(OP_RETURN, a, count, 0, sline(stmt))
}

func compileIfStmt(context *funcContext, stmt *qsa.IfStmt) {
	thenlabel := context.NewLabel()
	elselabel := context.NewLabel()
	endlabel := context.NewLabel()

	compileBranchCondition(context, context.RegTop(), stmt.Condition, thenlabel, elselabel, false)
	context.SetLabelPc(thenlabel, context.Code.LastPC())
	compileBlock(context, stmt.Then)
	if len(stmt.Else) > 0 {
		context.Code.AddASbx(OP_JMP, 0, endlabel, sline(stmt))
	}
	context.SetLabelPc(elselabel, context.Code.LastPC())
	if len(stmt.Else) > 0 {
		compileBlock(context, stmt.Else)
		context.SetLabelPc(endlabel, context.Code.LastPC())
	}

}

func compileBranchCondition(context *funcContext, reg int, expr qsa.Expr, thenlabel, elselabel int, hasnextcond bool) {
	code := context.Code
	flip := 0
	jumplabel := elselabel
	if hasnextcond {
		flip = 1
		jumplabel = thenlabel
	}

	switch ex := expr.(type) {
	case *qsa.FalseExpr, *qsa.NilExpr:
		if !hasnextcond {
			code.AddASbx(OP_JMP, 0, elselabel, sline(expr))
			return
		}
	case *qsa.TrueExpr, *qsa.NumberExpr, *qsa.StringExpr:
		if !hasnextcond {
			return
		}
	case *qsa.UnaryNotOpExpr:
		compileBranchCondition(context, reg, ex.Expr, elselabel, thenlabel, !hasnextcond)
		return
	case *qsa.LogicalOpExpr:
		switch ex.Operator {
		case "and":
			nextcondlabel := context.NewLabel()
			compileBranchCondition(context, reg, ex.Lhs, nextcondlabel, elselabel, false)
			context.SetLabelPc(nextcondlabel, context.Code.LastPC())
			compileBranchCondition(context, reg, ex.Rhs, thenlabel, elselabel, hasnextcond)
		case "or":
			nextcondlabel := context.NewLabel()
			compileBranchCondition(context, reg, ex.Lhs, thenlabel, nextcondlabel, true)
			context.SetLabelPc(nextcondlabel, context.Code.LastPC())
			compileBranchCondition(context, reg, ex.Rhs, thenlabel, elselabel, hasnextcond)
		}
		return
	case *qsa.RelationalOpExpr:
		compileRelationalOpExprAux(context, reg, ex, flip, jumplabel)
		return
	}

	a := reg
	compileExprWithMVPropagation(context, expr, &reg, &a)
	code.AddABC(OP_TEST, a, 0, 0^flip, sline(expr))
	code.AddASbx(OP_JMP, 0, jumplabel, sline(expr))
}

func compileWhileStmt(context *funcContext, stmt *qsa.WhileStmt) {
	thenlabel := context.NewLabel()
	elselabel := context.NewLabel()
	condlabel := context.NewLabel()

	context.SetLabelPc(condlabel, context.Code.LastPC())
	compileBranchCondition(context, context.RegTop(), stmt.Condition, thenlabel, elselabel, false)
	context.SetLabelPc(thenlabel, context.Code.LastPC())
	context.EnterBlock(elselabel, stmt)
	compileSeg(context, stmt.Stmts)
	context.CloseUpvalues()
	context.Code.AddASbx(OP_JMP, 0, condlabel, eline(stmt))
	context.LeaveBlock()
	context.SetLabelPc(elselabel, context.Code.LastPC())
}

func compileRepeatStmt(context *funcContext, stmt *qsa.RepeatStmt) {
	initlabel := context.NewLabel()
	thenlabel := context.NewLabel()
	elselabel := context.NewLabel()

	context.SetLabelPc(initlabel, context.Code.LastPC())
	context.SetLabelPc(elselabel, context.Code.LastPC())
	context.EnterBlock(thenlabel, stmt)
	compileSeg(context, stmt.Stmts)
	compileBranchCondition(context, context.RegTop(), stmt.Condition, thenlabel, elselabel, false)

	context.SetLabelPc(thenlabel, context.Code.LastPC())
	n := context.LeaveBlock()

	if n > -1 {
		label := context.NewLabel()
		context.Code.AddASbx(OP_JMP, 0, label, eline(stmt))
		context.SetLabelPc(elselabel, context.Code.LastPC())
		context.Code.AddABC(OP_CLOSE, n, 0, 0, eline(stmt))
		context.Code.AddASbx(OP_JMP, 0, initlabel, eline(stmt))
		context.SetLabelPc(label, context.Code.LastPC())
	}

}

func compileBreakStmt(context *funcContext, stmt *qsa.BreakStmt) {
	for block := context.Block; block != nil; block = block.Parent {
		if label := block.BreakLabel; label != labelNoJump {
			if block.RefUpvalue {
				context.Code.AddABC(OP_CLOSE, block.Parent.LocalVars.LastIndex(), 0, 0, sline(stmt))
			}
			context.Code.AddASbx(OP_JMP, 0, label, sline(stmt))
			return
		}
	}
	raiseCompileError(context, sline(stmt), "no loop to break")
}

func compileFuncDefStmt(context *funcContext, stmt *qsa.FuncDefStmt) {
	if stmt.Name.Func == nil {
		reg := context.RegTop()
		var treg, kreg int
		compileExprWithKMVPropagation(context, stmt.Name.Receiver, &reg, &treg)
		kreg = loadRk(context, &reg, stmt.Func, LString(stmt.Name.Method))
		compileExpr(context, reg, stmt.Func, ecfuncdef)
		context.Code.AddABC(OP_SETTABLE, treg, kreg, reg, sline(stmt.Name.Receiver))
	} else {
		astmt := &qsa.AssignStmt{Lhs: []qsa.Expr{stmt.Name.Func}, Rhs: []qsa.Expr{stmt.Func}}
		astmt.SetLine(sline(stmt.Func))
		astmt.SetLastLine(eline(stmt.Func))
		compileAssignStmt(context, astmt)
	}
}

func compileNumberForStmt(context *funcContext, stmt *qsa.NumberForStmt) {
	code := context.Code
	endlabel := context.NewLabel()
	ec := &expcontext{}

	context.EnterBlock(endlabel, stmt)
	reg := context.RegTop()
	rindex := context.RegisterLocalVar("(for index)")
	ecupdate(ec, ecLocal, rindex, 0)
	compileExpr(context, reg, stmt.Init, ec)

	reg = context.RegTop()
	rlimit := context.RegisterLocalVar("(for limit)")
	ecupdate(ec, ecLocal, rlimit, 0)
	compileExpr(context, reg, stmt.Limit, ec)

	reg = context.RegTop()
	rstep := context.RegisterLocalVar("(for step)")
	if stmt.Step == nil {
		stmt.Step = &qsa.NumberExpr{Value: "1"}
		stmt.Step.SetLine(sline(stmt.Init))
	}
	ecupdate(ec, ecLocal, rstep, 0)
	compileExpr(context, reg, stmt.Step, ec)

	code.AddASbx(OP_FORPREP, rindex, 0, sline(stmt))

	context.RegisterLocalVar(stmt.Name)

	bodypc := code.LastPC()
	compileSeg(context, stmt.Stmts)

	context.LeaveBlock()

	flpc := code.LastPC()
	code.AddASbx(OP_FORLOOP, rindex, bodypc-(flpc+1), sline(stmt))

	context.SetLabelPc(endlabel, code.LastPC())
	code.SetSbx(bodypc, flpc-bodypc)

}

func compileGenericForStmt(context *funcContext, stmt *qsa.GenericForStmt) {
	code := context.Code
	endlabel := context.NewLabel()
	bodylabel := context.NewLabel()
	fllabel := context.NewLabel()
	nnames := len(stmt.Names)

	context.EnterBlock(endlabel, stmt)
	rgen := context.RegisterLocalVar("(for generator)")
	context.RegisterLocalVar("(for state)")
	context.RegisterLocalVar("(for control)")

	compileRegAssignment(context, stmt.Names, stmt.Exprs, context.RegTop()-3, 3, sline(stmt))

	code.AddASbx(OP_JMP, 0, fllabel, sline(stmt))

	for _, name := range stmt.Names {
		context.RegisterLocalVar(name)
	}

	context.SetLabelPc(bodylabel, code.LastPC())
	compileSeg(context, stmt.Stmts)

	context.LeaveBlock()

	context.SetLabelPc(fllabel, code.LastPC())
	code.AddABC(OP_TFORLOOP, rgen, 0, nnames, sline(stmt))
	code.AddASbx(OP_JMP, 0, bodylabel, sline(stmt))

	context.SetLabelPc(endlabel, code.LastPC())
}

func compileExpr(context *funcContext, reg int, expr qsa.Expr, ec *expcontext) int {
	code := context.Code
	sreg := savereg(ec, reg)
	sused := 1
	if sreg < reg {
		sused = 0
	}

	switch ex := expr.(type) {
	case *qsa.StringExpr:
		code.AddABx(OP_LOADK, sreg, context.ConstIndex(LString(ex.Value)), sline(ex))
		return sused
	case *qsa.NumberExpr:
		num, err := parseNumber(ex.Value)
		if err != nil {
			num = LNumber(math.NaN())
		}
		code.AddABx(OP_LOADK, sreg, context.ConstIndex(num), sline(ex))
		return sused
	case *constLValueExpr:
		code.AddABx(OP_LOADK, sreg, context.ConstIndex(ex.Value), sline(ex))
		return sused
	case *qsa.NilExpr:
		code.AddABC(OP_LOADNIL, sreg, sreg, 0, sline(ex))
		return sused
	case *qsa.FalseExpr:
		code.AddABC(OP_LOADBOOL, sreg, 0, 0, sline(ex))
		return sused
	case *qsa.TrueExpr:
		code.AddABC(OP_LOADBOOL, sreg, 1, 0, sline(ex))
		return sused
	case *qsa.IdentExpr:
		switch getIdentRefType(context, context, ex) {
		case ecGlobal:
			code.AddABx(OP_GETGLOBAL, sreg, context.ConstIndex(LString(ex.Value)), sline(ex))
		case ecUpvalue:
			code.AddABC(OP_GETUPVAL, sreg, context.Upvalues.RegisterUnique(ex.Value), 0, sline(ex))
		case ecLocal:
			b := context.FindLocalVar(ex.Value)
			code.AddABC(OP_MOVE, sreg, b, 0, sline(ex))
		}
		return sused
	case *qsa.Comma3Expr:
		if context.Proto.IsVarArg == 0 {
			raiseCompileError(context, sline(ex), "cannot use '...' outside a vararg proc")
		}
		context.Proto.IsVarArg &= ^VarArgNeedsArg
		code.AddABC(OP_VARARG, sreg, 2+ec.varargopt, 0, sline(ex))
		if context.RegTop() > (sreg+2+ec.varargopt) || ec.varargopt < -1 {
			return 0
		}
		return (sreg + 1 + ec.varargopt) - reg
	case *qsa.AttrGetExpr:
		a := sreg
		b := reg
		compileExprWithMVPropagation(context, ex.Object, &reg, &b)
		c := reg
		compileExprWithKMVPropagation(context, ex.Key, &reg, &c)
		opcode := OP_GETTABLE
		if _, ok := ex.Key.(*qsa.StringExpr); ok {
			opcode = OP_GETTABLEKS
		}
		code.AddABC(opcode, a, b, c, sline(ex))
		return sused
	case *qsa.OAListExpr:
		compileOAListExpr(context, reg, ex, ec)
		return 1
	case *qsa.ArithmeticOpExpr:
		compileArithmeticOpExpr(context, reg, ex, ec)
		return sused
	case *qsa.StringConcatOpExpr:
		compileStringConcatOpExpr(context, reg, ex, ec)
		return sused
	case *qsa.UnaryMinusOpExpr, *qsa.UnaryNotOpExpr, *qsa.UnaryLenOpExpr:
		compileUnaryOpExpr(context, reg, ex, ec)
		return sused
	case *qsa.RelationalOpExpr:
		compileRelationalOpExpr(context, reg, ex, ec)
		return sused
	case *qsa.LogicalOpExpr:
		compileLogicalOpExpr(context, reg, ex, ec)
		return sused
	case *qsa.FuncCallExpr:
		return compileFuncCallExpr(context, reg, ex, ec)
	case *qsa.ProcExpr:
		childcontext := newFuncContext(context.Proto.SourceName, context)
		compileProcExpr(childcontext, ex, ec)
		protono := len(context.Proto.ProcPrototypes)
		context.Proto.ProcPrototypes = append(context.Proto.ProcPrototypes, childcontext.Proto)
		code.AddABx(OP_CLOSURE, sreg, protono, sline(ex))
		for _, upvalue := range childcontext.Upvalues.List() {
			localidx, block := context.FindLocalVarAndBlock(upvalue.Name)
			if localidx > -1 {
				code.AddABC(OP_MOVE, 0, localidx, 0, sline(ex))
				block.RefUpvalue = true
			} else {
				upvalueidx := context.Upvalues.Find(upvalue.Name)
				if upvalueidx < 0 {
					upvalueidx = context.Upvalues.RegisterUnique(upvalue.Name)
				}
				code.AddABC(OP_GETUPVAL, 0, upvalueidx, 0, sline(ex))
			}
		}
		return sused
	default:
		e := reflect.TypeOf(ex).Elem().Name()
		log.Error().Str("expr", e).
			Msgf("Expression %v not implemented", e)
		os.Exit(RCERROR)
	}

	log.Fatal().Msg("Q language failure, expression illogic")
	return sused
}

func compileExprWithPropagation(context *funcContext, expr qsa.Expr, reg *int, save *int, propergator func(int, *int, *int, int)) {
	reginc := compileExpr(context, *reg, expr, ecnone(0))
	if _, ok := expr.(*qsa.LogicalOpExpr); ok {
		*save = *reg
		*reg = *reg + reginc
	} else {
		propergator(context.RegTop(), save, reg, reginc)
	}
}

func compileExprWithKMVPropagation(context *funcContext, expr qsa.Expr, reg *int, save *int) {
	compileExprWithPropagation(context, expr, reg, save, context.Code.PropagateKMV)
}

func compileExprWithMVPropagation(context *funcContext, expr qsa.Expr, reg *int, save *int) {
	compileExprWithPropagation(context, expr, reg, save, context.Code.PropagateMV)
}

func constFold(exp qsa.Expr) qsa.Expr {
	switch expr := exp.(type) {
	case *qsa.ArithmeticOpExpr:
		lvalue, lisconst := lnumberValue(expr.Lhs)
		rvalue, risconst := lnumberValue(expr.Rhs)
		if lisconst && risconst {
			switch expr.Operator {
			case "+":
				return &constLValueExpr{Value: lvalue + rvalue}
			case "-":
				return &constLValueExpr{Value: lvalue - rvalue}
			case "*":
				return &constLValueExpr{Value: lvalue * rvalue}
			case "/":
				return &constLValueExpr{Value: lvalue / rvalue}
			case "%":
				return &constLValueExpr{Value: oaModulo(lvalue, rvalue)}
			case "^":
				return &constLValueExpr{Value: LNumber(math.Pow(float64(lvalue), float64(rvalue)))}
			default:
				log.Error().Msgf("Binary operator %s invalid", expr.Operator)
			}
		} else {
			retexpr := *expr
			retexpr.Lhs = constFold(expr.Lhs)
			retexpr.Rhs = constFold(expr.Rhs)
			return &retexpr
		}
	case *qsa.UnaryMinusOpExpr:
		expr.Expr = constFold(expr.Expr)
		if value, ok := lnumberValue(expr.Expr); ok {
			return &constLValueExpr{Value: LNumber(-value)}
		}
		return expr
	default:

		return exp
	}
	return exp
}

func compileProcExpr(context *funcContext, funcexpr *qsa.ProcExpr, ec *expcontext) {
	context.Proto.LineDefined = sline(funcexpr)
	context.Proto.LastLineDefined = eline(funcexpr)
	if len(funcexpr.ParList.Names) > maxRegisters {
		raiseCompileError(context, context.Proto.LineDefined, "register overflow")
	}
	context.Proto.NumParameters = uint8(len(funcexpr.ParList.Names))
	if ec.ctype == ecMethod {
		context.Proto.NumParameters += 1
		context.RegisterLocalVar("self")
	}
	for _, name := range funcexpr.ParList.Names {
		context.RegisterLocalVar(name)
	}
	if funcexpr.ParList.HasVargs {
		if CompatVarArg {
			context.Proto.IsVarArg = VarArgHasArg | VarArgNeedsArg
			if context.Parent != nil {
				context.RegisterLocalVar("arg")
			}
		}
		context.Proto.IsVarArg |= VarArgIsVarArg
	}

	compileSeg(context, funcexpr.Stmts)

	context.Code.AddABC(OP_RETURN, 0, 1, 0, eline(funcexpr))
	context.EndScope()
	context.Proto.Code = context.Code.List()
	context.Proto.DbgSourcePositions = context.Code.PosList()
	context.Proto.DbgUpvalues = context.Upvalues.Names()
	context.Proto.NumUpvalues = uint8(len(context.Proto.DbgUpvalues))
	for _, clv := range context.Proto.Constants {
		sv := ""
		if slv, ok := clv.(LString); ok {
			sv = string(slv)
		}
		context.Proto.stringConstants = append(context.Proto.stringConstants, sv)
	}
	patchCode(context)
}

func compileOAListExpr(context *funcContext, reg int, ex *qsa.OAListExpr, ec *expcontext) {
	code := context.Code
	listreg := reg
	reg++
	code.AddABC(OP_NEWTABLE, listreg, 0, 0, sline(ex))
	listpc := code.LastPC()
	regbase := reg

	arraycount := 0
	lastvararg := false
	for i, field := range ex.Fields {
		islast := i == len(ex.Fields)-1
		if field.Key == nil {
			if islast && isVarArgReturnExpr(field.Value) {
				reg += compileExpr(context, reg, field.Value, ecnone(-2))
				lastvararg = true
			} else {
				reg += compileExpr(context, reg, field.Value, ecnone(0))
				arraycount += 1
			}
		} else {
			regorg := reg
			b := reg
			compileExprWithKMVPropagation(context, field.Key, &reg, &b)
			c := reg
			compileExprWithKMVPropagation(context, field.Value, &reg, &c)
			opcode := OP_SETTABLE
			if _, ok := field.Key.(*qsa.StringExpr); ok {
				opcode = OP_SETTABLEKS
			}
			code.AddABC(opcode, listreg, b, c, sline(ex))
			reg = regorg
		}
		flush := arraycount % FieldsPerFlush
		if (arraycount != 0 && (flush == 0 || islast)) || lastvararg {
			reg = regbase
			num := flush
			if num == 0 {
				num = FieldsPerFlush
			}
			c := (arraycount-1)/FieldsPerFlush + 1
			b := num
			if islast && isVarArgReturnExpr(field.Value) {
				b = 0
			}
			line := field.Value
			if field.Key != nil {
				line = field.Key
			}
			if c > 511 {
				c = 0
			}
			code.AddABC(OP_SETLIST, listreg, b, c, sline(line))
			if c == 0 {
				code.Add(uint32(c), sline(line))
			}
		}
	}
	code.SetB(listpc, int2Fb(arraycount))
	code.SetC(listpc, int2Fb(len(ex.Fields)-arraycount))
	if ec.ctype == ecLocal && ec.reg != listreg {
		code.AddABC(OP_MOVE, ec.reg, listreg, 0, sline(ex))
	}
}

func compileArithmeticOpExpr(context *funcContext, reg int, expr *qsa.ArithmeticOpExpr, ec *expcontext) {
	exp := constFold(expr)
	if ex, ok := exp.(*constLValueExpr); ok {
		exp.SetLine(sline(expr))
		compileExpr(context, reg, ex, ec)
		return
	}
	expr, _ = exp.(*qsa.ArithmeticOpExpr)
	a := savereg(ec, reg)
	b := reg
	compileExprWithKMVPropagation(context, expr.Lhs, &reg, &b)
	c := reg
	compileExprWithKMVPropagation(context, expr.Rhs, &reg, &c)

	op := 0
	switch expr.Operator {
	case "+":
		op = OP_ADD
	case "-":
		op = OP_SUB
	case "*":
		op = OP_MUL
	case "/":
		op = OP_DIV
	case "%":
		op = OP_MOD
	case "^":
		op = OP_POW
	}
	context.Code.AddABC(op, a, b, c, sline(expr))
}

func compileStringConcatOpExpr(context *funcContext, reg int, expr *qsa.StringConcatOpExpr, ec *expcontext) {
	code := context.Code
	crange := 1
	for current := expr.Rhs; current != nil; {
		if ex, ok := current.(*qsa.StringConcatOpExpr); ok {
			crange += 1
			current = ex.Rhs
		} else {
			current = nil
		}
	}
	a := savereg(ec, reg)
	basereg := reg
	reg += compileExpr(context, reg, expr.Lhs, ecnone(0))
	reg += compileExpr(context, reg, expr.Rhs, ecnone(0))
	for pc := code.LastPC(); pc != 0 && opGetOpCode(code.At(pc)) == OP_CONCAT; pc-- {
		code.Pop()
	}
	code.AddABC(OP_CONCAT, a, basereg, basereg+crange, sline(expr))
}

func compileUnaryOpExpr(context *funcContext, reg int, expr qsa.Expr, ec *expcontext) {
	opcode := 0
	code := context.Code
	var operandexpr qsa.Expr
	switch ex := expr.(type) {
	case *qsa.UnaryMinusOpExpr:
		exp := constFold(ex)
		if lvexpr, ok := exp.(*constLValueExpr); ok {
			exp.SetLine(sline(expr))
			compileExpr(context, reg, lvexpr, ec)
			return
		}
		ex, _ = exp.(*qsa.UnaryMinusOpExpr)
		operandexpr = ex.Expr
		opcode = OP_UNM
	case *qsa.UnaryNotOpExpr:
		switch ex.Expr.(type) {
		case *qsa.TrueExpr:
			code.AddABC(OP_LOADBOOL, savereg(ec, reg), 0, 0, sline(expr))
			return
		case *qsa.FalseExpr, *qsa.NilExpr:
			code.AddABC(OP_LOADBOOL, savereg(ec, reg), 1, 0, sline(expr))
			return
		default:
			opcode = OP_NOT
			operandexpr = ex.Expr
		}
	case *qsa.UnaryLenOpExpr:
		opcode = OP_LEN
		operandexpr = ex.Expr
	}

	a := savereg(ec, reg)
	b := reg
	compileExprWithMVPropagation(context, operandexpr, &reg, &b)
	code.AddABC(opcode, a, b, 0, sline(expr))
}

func compileRelationalOpExprAux(context *funcContext, reg int, expr *qsa.RelationalOpExpr, flip int, label int) {
	code := context.Code
	b := reg
	compileExprWithKMVPropagation(context, expr.Lhs, &reg, &b)
	c := reg
	compileExprWithKMVPropagation(context, expr.Rhs, &reg, &c)
	switch expr.Operator {
	case "<":
		code.AddABC(OP_LT, 0^flip, b, c, sline(expr))
	case ">":
		code.AddABC(OP_LT, 0^flip, c, b, sline(expr))
	case "<=":
		code.AddABC(OP_LE, 0^flip, b, c, sline(expr))
	case ">=":
		code.AddABC(OP_LE, 0^flip, c, b, sline(expr))
	case "==":
		code.AddABC(OP_EQ, 0^flip, b, c, sline(expr))
	case "~=":
		code.AddABC(OP_EQ, 1^flip, b, c, sline(expr))
	}
	code.AddASbx(OP_JMP, 0, label, sline(expr))
}

func compileRelationalOpExpr(context *funcContext, reg int, expr *qsa.RelationalOpExpr, ec *expcontext) {
	a := savereg(ec, reg)
	code := context.Code
	jumplabel := context.NewLabel()
	compileRelationalOpExprAux(context, reg, expr, 1, jumplabel)
	code.AddABC(OP_LOADBOOL, a, 0, 1, sline(expr))
	context.SetLabelPc(jumplabel, code.LastPC())
	code.AddABC(OP_LOADBOOL, a, 1, 0, sline(expr))
}

func compileLogicalOpExpr(context *funcContext, reg int, expr *qsa.LogicalOpExpr, ec *expcontext) {
	a := savereg(ec, reg)
	code := context.Code
	endlabel := context.NewLabel()
	lb := &lblabels{context.NewLabel(), context.NewLabel(), endlabel, false}
	nextcondlabel := context.NewLabel()
	if expr.Operator == "and" {
		compileLogicalOpExprAux(context, reg, expr.Lhs, ec, nextcondlabel, endlabel, false, lb)
		context.SetLabelPc(nextcondlabel, code.LastPC())
		compileLogicalOpExprAux(context, reg, expr.Rhs, ec, endlabel, endlabel, false, lb)
	} else {
		compileLogicalOpExprAux(context, reg, expr.Lhs, ec, endlabel, nextcondlabel, true, lb)
		context.SetLabelPc(nextcondlabel, code.LastPC())
		compileLogicalOpExprAux(context, reg, expr.Rhs, ec, endlabel, endlabel, false, lb)
	}

	if lb.b {
		context.SetLabelPc(lb.f, code.LastPC())
		code.AddABC(OP_LOADBOOL, a, 0, 1, sline(expr))
		context.SetLabelPc(lb.t, code.LastPC())
		code.AddABC(OP_LOADBOOL, a, 1, 0, sline(expr))
	}

	lastinst := code.Last()
	if opGetOpCode(lastinst) == OP_JMP && opGetArgSbx(lastinst) == endlabel {
		code.Pop()
	}

	context.SetLabelPc(endlabel, code.LastPC())
}

func compileLogicalOpExprAux(context *funcContext, reg int, expr qsa.Expr, ec *expcontext, thenlabel, elselabel int, hasnextcond bool, lb *lblabels) {
	code := context.Code
	flip := 0
	jumplabel := elselabel
	if hasnextcond {
		flip = 1
		jumplabel = thenlabel
	}

	switch ex := expr.(type) {
	case *qsa.FalseExpr:
		if elselabel == lb.e {
			code.AddASbx(OP_JMP, 0, lb.f, sline(expr))
			lb.b = true
		} else {
			code.AddASbx(OP_JMP, 0, elselabel, sline(expr))
		}
		return
	case *qsa.NilExpr:
		if elselabel == lb.e {
			compileExpr(context, reg, expr, ec)
			code.AddASbx(OP_JMP, 0, lb.e, sline(expr))
		} else {
			code.AddASbx(OP_JMP, 0, elselabel, sline(expr))
		}
		return
	case *qsa.TrueExpr:
		if thenlabel == lb.e {
			code.AddASbx(OP_JMP, 0, lb.t, sline(expr))
			lb.b = true
		} else {
			code.AddASbx(OP_JMP, 0, thenlabel, sline(expr))
		}
		return
	case *qsa.NumberExpr, *qsa.StringExpr:
		if thenlabel == lb.e {
			compileExpr(context, reg, expr, ec)
			code.AddASbx(OP_JMP, 0, lb.e, sline(expr))
		} else {
			code.AddASbx(OP_JMP, 0, thenlabel, sline(expr))
		}
		return
	case *qsa.LogicalOpExpr:
		switch ex.Operator {
		case "and":
			nextcondlabel := context.NewLabel()
			compileLogicalOpExprAux(context, reg, ex.Lhs, ec, nextcondlabel, elselabel, false, lb)
			context.SetLabelPc(nextcondlabel, context.Code.LastPC())
			compileLogicalOpExprAux(context, reg, ex.Rhs, ec, thenlabel, elselabel, hasnextcond, lb)
		case "or":
			nextcondlabel := context.NewLabel()
			compileLogicalOpExprAux(context, reg, ex.Lhs, ec, thenlabel, nextcondlabel, true, lb)
			context.SetLabelPc(nextcondlabel, context.Code.LastPC())
			compileLogicalOpExprAux(context, reg, ex.Rhs, ec, thenlabel, elselabel, hasnextcond, lb)
		}
		return
	case *qsa.RelationalOpExpr:
		if thenlabel == elselabel {
			flip ^= 1
			jumplabel = lb.t
			lb.b = true
		} else if thenlabel == lb.e {
			jumplabel = lb.t
			lb.b = true
		} else if elselabel == lb.e {
			jumplabel = lb.f
			lb.b = true
		}
		compileRelationalOpExprAux(context, reg, ex, flip, jumplabel)
		return
	}

	if !hasnextcond && thenlabel == elselabel {
		reg += compileExpr(context, reg, expr, ec)
	} else {
		a := reg
		sreg := savereg(ec, a)
		reg += compileExpr(context, reg, expr, ecnone(0))
		if sreg == a {
			code.AddABC(OP_TEST, a, 0, 0^flip, sline(expr))
		} else {
			code.AddABC(OP_TESTSET, sreg, a, 0^flip, sline(expr))
		}
	}
	code.AddASbx(OP_JMP, 0, jumplabel, sline(expr))
}

func compileFuncCallExpr(context *funcContext, reg int, expr *qsa.FuncCallExpr, ec *expcontext) int {
	funcreg := reg
	if ec.ctype == ecLocal && ec.reg == (int(context.Proto.NumParameters)-1) {
		funcreg = ec.reg
		reg = ec.reg
	}
	argc := len(expr.Args)
	islastvararg := false
	name := "(anonymous)"

	if expr.Func != nil { // hoge.func()
		reg += compileExpr(context, reg, expr.Func, ecnone(0))
		name = getExprName(context, expr.Func)
	} else { // hoge:method()
		b := reg
		compileExprWithMVPropagation(context, expr.Receiver, &reg, &b)
		c := loadRk(context, &reg, expr, LString(expr.Method))
		context.Code.AddABC(OP_SELF, funcreg, b, c, sline(expr))
		// increments a register for an implicit "self"
		reg = b + 1
		reg2 := funcreg + 2
		if reg2 > reg {
			reg = reg2
		}
		argc += 1
		name = string(expr.Method)
	}

	for i, ar := range expr.Args {
		islastvararg = (i == len(expr.Args)-1) && isVarArgReturnExpr(ar)
		if islastvararg {
			compileExpr(context, reg, ar, ecnone(-2))
		} else {
			reg += compileExpr(context, reg, ar, ecnone(0))
		}
	}
	b := argc + 1
	if islastvararg {
		b = 0
	}
	context.Code.AddABC(OP_CALL, funcreg, b, ec.varargopt+2, sline(expr))
	context.Proto.DbgCalls = append(context.Proto.DbgCalls, DbgCall{Pc: context.Code.LastPC(), Name: name})

	if ec.varargopt == 0 && ec.ctype == ecLocal && funcreg != ec.reg {
		context.Code.AddABC(OP_MOVE, ec.reg, funcreg, 0, sline(expr))
		return 1
	}
	if context.RegTop() > (funcreg+2+ec.varargopt) || ec.varargopt < -1 {
		return 0
	}
	return ec.varargopt + 1
}

func loadRk(context *funcContext, reg *int, expr qsa.Expr, cnst LValue) int {
	cindex := context.ConstIndex(cnst)
	if cindex <= opMaxIndexRk {
		return opRkAsk(cindex)
	} else {
		ret := *reg
		*reg++
		context.Code.AddABx(OP_LOADK, ret, cindex, sline(expr))
		return ret
	}
}

func getIdentRefType(context *funcContext, current *funcContext, expr *qsa.IdentExpr) expContextType {
	if current == nil {
		return ecGlobal
	} else if current.FindLocalVar(expr.Value) > -1 {
		if current == context {
			return ecLocal
		}
		return ecUpvalue
	}
	return getIdentRefType(context, current.Parent, expr)
}

func getExprName(context *funcContext, expr qsa.Expr) string {
	switch ex := expr.(type) {
	case *qsa.IdentExpr:
		return ex.Value
	case *qsa.AttrGetExpr:
		switch kex := ex.Key.(type) {
		case *qsa.StringExpr:
			return kex.Value
		}
		return "?"
	}
	return "?"
}

func patchCode(context *funcContext) {
	maxreg := 1
	if np := int(context.Proto.NumParameters); np > 1 {
		maxreg = np
	}
	moven := 0
	code := context.Code.List()
	for pc := 0; pc < len(code); pc++ {
		inst := code[pc]
		curop := opGetOpCode(inst)
		switch curop {
		case OP_CLOSURE:
			pc += int(context.Proto.ProcPrototypes[opGetArgBx(inst)].NumUpvalues)
			moven = 0
			continue
		case OP_SETGLOBAL, OP_SETUPVAL, OP_EQ, OP_LT, OP_LE, OP_TEST,
			OP_TAILCALL, OP_RETURN, OP_FORPREP, OP_FORLOOP, OP_TFORLOOP,
			OP_SETLIST, OP_CLOSE:
			/* nothing to do */
		case OP_CALL:
			if reg := opGetArgA(inst) + opGetArgC(inst) - 2; reg > maxreg {
				maxreg = reg
			}
		case OP_VARARG:
			if reg := opGetArgA(inst) + opGetArgB(inst) - 1; reg > maxreg {
				maxreg = reg
			}
		case OP_SELF:
			if reg := opGetArgA(inst) + 1; reg > maxreg {
				maxreg = reg
			}
		case OP_LOADNIL:
			if reg := opGetArgB(inst); reg > maxreg {
				maxreg = reg
			}
		case OP_JMP: // jump to jump optimization
			distance := 0
			count := 0 // avoiding infinite loops
			for jmp := inst; opGetOpCode(jmp) == OP_JMP && count < 5; jmp = context.Code.At(pc + distance + 1) {
				d := context.GetLabelPc(opGetArgSbx(jmp)) - pc
				if d > opMaxArgSbx {
					if distance == 0 {
						raiseCompileError(context, context.Proto.LineDefined, "too long to jump.")
					}
					break
				}
				distance = d
				count++
			}
			if distance == 0 {
				context.Code.SetOpCode(pc, OP_NOP)
			} else {
				context.Code.SetSbx(pc, distance)
			}
		default:
			if reg := opGetArgA(inst); reg > maxreg {
				maxreg = reg
			}
		}

		// bulk move optimization(reducing op dipatch costs)
		if curop == OP_MOVE {
			moven++
		} else {
			if moven > 1 {
				context.Code.SetOpCode(pc-moven, OP_MOVEN)
				context.Code.SetC(pc-moven, intMin(moven-1, opMaxArgsC))
			}
			moven = 0
		}
	}
	maxreg++
	if maxreg > maxRegisters {
		raiseCompileError(context, context.Proto.LineDefined, "register overflow(too many local variables)")
	}
	context.Proto.NumUsedRegisters = uint8(maxreg)
}

func Compile(segment []qsa.Stmt, name string) (proto *ProcProto, err error) {
	defer func() {
		if rcv := recover(); rcv != nil {
			if _, ok := rcv.(*CompileError); ok {
				err = rcv.(error)
			} else {
				panic(rcv)
			}
		}
	}()
	err = nil
	parlist := &qsa.ParList{HasVargs: true, Names: []string{}}
	funcexpr := &qsa.ProcExpr{ParList: parlist, Stmts: segment}
	context := newFuncContext(name, nil)
	compileProcExpr(context, funcexpr, ecnone(0))
	proto = context.Proto
	return
}
