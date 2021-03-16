%{
package qsp

import (
	"github.com/x0ray/q/qs/qsa"
)
%}

%type<stmts> segment
%type<stmts> segment1
%type<stmts> block
%type<stmt>  stat
%type<stmts> elseifs
%type<stmt>  laststat
%type<funcname> funcname
%type<funcname> funcname1
%type<exprlist> varlist
%type<expr> var
%type<namelist> namelist
%type<exprlist> exprlist
%type<expr> expr
%type<expr> string
%type<expr> prefixexp
%type<expr> proccall
%type<expr> aproccall
%type<exprlist> args
%type<expr> proc
%type<funcexpr> funcbody
%type<parlist> parlist
%type<expr> listconstructor
%type<fieldlist> fieldlist
%type<field> field
%type<fieldsep> fieldsep

%union {
  token  qsa.Token

  stmts    []qsa.Stmt
  stmt     qsa.Stmt

  funcname *qsa.FuncName
  funcexpr *qsa.FunctionExpr

  exprlist []qsa.Expr
  expr   qsa.Expr

  fieldlist []*qsa.Field
  field     *qsa.Field
  fieldsep  string

  namelist []string
  parlist  *qsa.ParList
}

/* Reserved words */
%token<token> TAnd TBreak TDo TElse TElseIf TEnd TFalse TFor TProc TIf TIn TLocal TNil TNot TOr TReturn TRepeat TThen TTrue TUntil TWhile 

/* Literals */
%token<token> TEqeq TNeq TLte TGte T2Comma T3Comma TIdent TNumber TString '{' '('

/* Operators */
%left TOr
%left TAnd
%left '>' '<' TGte TLte TEqeq TNeq
%right T2Comma
%left '+' '-'
%left '*' '/' '%'
%right UNARY /* not # -(unary) */
%right '^'

%%

segment: 
        segment1 {
            $$ = $1
            if l, ok := yylex.(*Lexer); ok {
                l.Stmts = $$
            }
        } |
        segment1 laststat {
            $$ = append($1, $2)
            if l, ok := yylex.(*Lexer); ok {
                l.Stmts = $$
            }
        } | 
        segment1 laststat ';' {
            $$ = append($1, $2)
            if l, ok := yylex.(*Lexer); ok {
                l.Stmts = $$
            }
        }

segment1: 
        {
            $$ = []qsa.Stmt{}
        } |
        segment1 stat {
            $$ = append($1, $2)
        } | 
        segment1 ';' {
            $$ = $1
        }

block: 
        segment {
            $$ = $1
        }

stat:
        varlist '=' exprlist {
            $$ = &qsa.AssignStmt{Lhs: $1, Rhs: $3}
            $$.SetLine($1[0].Line())
        } |
        /* 'stat = proccal' causes a reduce/reduce conflict */
        prefixexp {
            if _, ok := $1.(*qsa.FuncCallExpr); !ok {
               yylex.(*Lexer).Error("parse error")
            } else {
              $$ = &qsa.FuncCallStmt{Expr: $1}
              $$.SetLine($1.Line())
            }
        } |
        TDo block TEnd {
            $$ = &qsa.DoBlockStmt{Stmts: $2}
            $$.SetLine($1.Pos.Line)
            $$.SetLastLine($3.Pos.Line)
        } |
        TWhile expr TDo block TEnd {
            $$ = &qsa.WhileStmt{Condition: $2, Stmts: $4}
            $$.SetLine($1.Pos.Line)
            $$.SetLastLine($5.Pos.Line)
        } |
        TRepeat block TUntil expr {
            $$ = &qsa.RepeatStmt{Condition: $4, Stmts: $2}
            $$.SetLine($1.Pos.Line)
            $$.SetLastLine($4.Line())
        } |
        TIf expr TThen block elseifs TEnd {
            $$ = &qsa.IfStmt{Condition: $2, Then: $4}
            cur := $$
            for _, elseif := range $5 {
                cur.(*qsa.IfStmt).Else = []qsa.Stmt{elseif}
                cur = elseif
            }
            $$.SetLine($1.Pos.Line)
            $$.SetLastLine($6.Pos.Line)
        } |
        TIf expr TThen block elseifs TElse block TEnd {
            $$ = &qsa.IfStmt{Condition: $2, Then: $4}
            cur := $$
            for _, elseif := range $5 {
                cur.(*qsa.IfStmt).Else = []qsa.Stmt{elseif}
                cur = elseif
            }
            cur.(*qsa.IfStmt).Else = $7
            $$.SetLine($1.Pos.Line)
            $$.SetLastLine($8.Pos.Line)
        } |
        TFor TIdent '=' expr ',' expr TDo block TEnd {
            $$ = &qsa.NumberForStmt{Name: $2.Str, Init: $4, Limit: $6, Stmts: $8}
            $$.SetLine($1.Pos.Line)
            $$.SetLastLine($9.Pos.Line)
        } |
        TFor TIdent '=' expr ',' expr ',' expr TDo block TEnd {
            $$ = &qsa.NumberForStmt{Name: $2.Str, Init: $4, Limit: $6, Step:$8, Stmts: $10}
            $$.SetLine($1.Pos.Line)
            $$.SetLastLine($11.Pos.Line)
        } |
        TFor namelist TIn exprlist TDo block TEnd {
            $$ = &qsa.GenericForStmt{Names:$2, Exprs:$4, Stmts: $6}
            $$.SetLine($1.Pos.Line)
            $$.SetLastLine($7.Pos.Line)
        } |
        TProc funcname funcbody {
            $$ = &qsa.FuncDefStmt{Name: $2, Func: $3}
            $$.SetLine($1.Pos.Line)
            $$.SetLastLine($3.LastLine())
        } |
        TLocal TProc TIdent funcbody {
            $$ = &qsa.LocalAssignStmt{Names:[]string{$3.Str}, Exprs: []qsa.Expr{$4}}
            $$.SetLine($1.Pos.Line)
            $$.SetLastLine($4.LastLine())
        } | 
        TLocal namelist '=' exprlist {
            $$ = &qsa.LocalAssignStmt{Names: $2, Exprs:$4}
            $$.SetLine($1.Pos.Line)
        } |
        TLocal namelist {
            $$ = &qsa.LocalAssignStmt{Names: $2, Exprs:[]qsa.Expr{}}
            $$.SetLine($1.Pos.Line)
        }

elseifs: 
        {
            $$ = []qsa.Stmt{}
        } | 
        elseifs TElseIf expr TThen block {
            $$ = append($1, &qsa.IfStmt{Condition: $3, Then: $5})
            $$[len($$)-1].SetLine($2.Pos.Line)
        }

laststat:
        TReturn {
            $$ = &qsa.ReturnStmt{Exprs:nil}
            $$.SetLine($1.Pos.Line)
        } |
        TReturn exprlist {
            $$ = &qsa.ReturnStmt{Exprs:$2}
            $$.SetLine($1.Pos.Line)
        } |
        TBreak  {
            $$ = &qsa.BreakStmt{}
            $$.SetLine($1.Pos.Line)
        }

funcname: 
        funcname1 {
            $$ = $1
        } |
        funcname1 ':' TIdent {
            $$ = &qsa.FuncName{Func:nil, Receiver:$1.Func, Method: $3.Str}
        }

funcname1:
        TIdent {
            $$ = &qsa.FuncName{Func: &qsa.IdentExpr{Value:$1.Str}}
            $$.Func.SetLine($1.Pos.Line)
        } | 
        funcname1 '.' TIdent {
            key:= &qsa.StringExpr{Value:$3.Str}
            key.SetLine($3.Pos.Line)
            fn := &qsa.AttrGetExpr{Object: $1.Func, Key: key}
            fn.SetLine($3.Pos.Line)
            $$ = &qsa.FuncName{Func: fn}
        }

varlist:
        var {
            $$ = []qsa.Expr{$1}
        } | 
        varlist ',' var {
            $$ = append($1, $3)
        }

var:
        TIdent {
            $$ = &qsa.IdentExpr{Value:$1.Str}
            $$.SetLine($1.Pos.Line)
        } |
        prefixexp '[' expr ']' {
            $$ = &qsa.AttrGetExpr{Object: $1, Key: $3}
            $$.SetLine($1.Line())
        } | 
        prefixexp '.' TIdent {
            key := &qsa.StringExpr{Value:$3.Str}
            key.SetLine($3.Pos.Line)
            $$ = &qsa.AttrGetExpr{Object: $1, Key: key}
            $$.SetLine($1.Line())
        }

namelist:
        TIdent {
            $$ = []string{$1.Str}
        } | 
        namelist ','  TIdent {
            $$ = append($1, $3.Str)
        }

exprlist:
        expr {
            $$ = []qsa.Expr{$1}
        } |
        exprlist ',' expr {
            $$ = append($1, $3)
        }

expr:
        TNil {
            $$ = &qsa.NilExpr{}
            $$.SetLine($1.Pos.Line)
        } | 
        TFalse {
            $$ = &qsa.FalseExpr{}
            $$.SetLine($1.Pos.Line)
        } | 
        TTrue {
            $$ = &qsa.TrueExpr{}
            $$.SetLine($1.Pos.Line)
        } | 
        TNumber {
            $$ = &qsa.NumberExpr{Value: $1.Str}
            $$.SetLine($1.Pos.Line)
        } | 
        T3Comma {
            $$ = &qsa.Comma3Expr{}
            $$.SetLine($1.Pos.Line)
        } |
        proc {
            $$ = $1
        } | 
        prefixexp {
            $$ = $1
        } |
        string {
            $$ = $1
        } |
        listconstructor {
            $$ = $1
        } |
        expr TOr expr {
            $$ = &qsa.LogicalOpExpr{Lhs: $1, Operator: "or", Rhs: $3}
            $$.SetLine($1.Line())
        } |
        expr TAnd expr {
            $$ = &qsa.LogicalOpExpr{Lhs: $1, Operator: "and", Rhs: $3}
            $$.SetLine($1.Line())
        } |
        expr '>' expr {
            $$ = &qsa.RelationalOpExpr{Lhs: $1, Operator: ">", Rhs: $3}
            $$.SetLine($1.Line())
        } |
        expr '<' expr {
            $$ = &qsa.RelationalOpExpr{Lhs: $1, Operator: "<", Rhs: $3}
            $$.SetLine($1.Line())
        } |
        expr TGte expr {
            $$ = &qsa.RelationalOpExpr{Lhs: $1, Operator: ">=", Rhs: $3}
            $$.SetLine($1.Line())
        } |
        expr TLte expr {
            $$ = &qsa.RelationalOpExpr{Lhs: $1, Operator: "<=", Rhs: $3}
            $$.SetLine($1.Line())
        } |
        expr TEqeq expr {
            $$ = &qsa.RelationalOpExpr{Lhs: $1, Operator: "==", Rhs: $3}
            $$.SetLine($1.Line())
        } |
        expr TNeq expr {
            $$ = &qsa.RelationalOpExpr{Lhs: $1, Operator: "~=", Rhs: $3}
            $$.SetLine($1.Line())
        } |
        expr T2Comma expr {
            $$ = &qsa.StringConcatOpExpr{Lhs: $1, Rhs: $3}
            $$.SetLine($1.Line())
        } |
        expr '+' expr {
            $$ = &qsa.ArithmeticOpExpr{Lhs: $1, Operator: "+", Rhs: $3}
            $$.SetLine($1.Line())
        } |
        expr '-' expr {
            $$ = &qsa.ArithmeticOpExpr{Lhs: $1, Operator: "-", Rhs: $3}
            $$.SetLine($1.Line())
        } |
        expr '*' expr {
            $$ = &qsa.ArithmeticOpExpr{Lhs: $1, Operator: "*", Rhs: $3}
            $$.SetLine($1.Line())
        } |
        expr '/' expr {
            $$ = &qsa.ArithmeticOpExpr{Lhs: $1, Operator: "/", Rhs: $3}
            $$.SetLine($1.Line())
        } |
        expr '%' expr {
            $$ = &qsa.ArithmeticOpExpr{Lhs: $1, Operator: "%", Rhs: $3}
            $$.SetLine($1.Line())
        } |
        expr '^' expr {
            $$ = &qsa.ArithmeticOpExpr{Lhs: $1, Operator: "^", Rhs: $3}
            $$.SetLine($1.Line())
        } |
        '-' expr %prec UNARY {
            $$ = &qsa.UnaryMinusOpExpr{Expr: $2}
            $$.SetLine($2.Line())
        } |
        TNot expr %prec UNARY {
            $$ = &qsa.UnaryNotOpExpr{Expr: $2}
            $$.SetLine($2.Line())
        } |
        '#' expr %prec UNARY {
            $$ = &qsa.UnaryLenOpExpr{Expr: $2}
            $$.SetLine($2.Line())
        }

string: 
        TString {
            $$ = &qsa.StringExpr{Value: $1.Str}
            $$.SetLine($1.Pos.Line)
        } 

prefixexp:
        var {
            $$ = $1
        } |
        aproccall {
            $$ = $1
        } |
        proccall {
            $$ = $1
        } |
        '(' expr ')' {
            $$ = $2
            $$.SetLine($1.Pos.Line)
        }

aproccall:
        '(' proccall ')' {
            $2.(*qsa.FuncCallExpr).AdjustRet = true
            $$ = $2
        }

proccall:
        prefixexp args {
            $$ = &qsa.FuncCallExpr{Func: $1, Args: $2}
            $$.SetLine($1.Line())
        } |
        prefixexp ':' TIdent args {
            $$ = &qsa.FuncCallExpr{Method: $3.Str, Receiver: $1, Args: $4}
            $$.SetLine($1.Line())
        }

args:
        '(' ')' {
            if yylex.(*Lexer).PNewLine {
               yylex.(*Lexer).TokenError($1, "ambiguous syntax (proc call x new statement)")
            }
            $$ = []qsa.Expr{}
        } |
        '(' exprlist ')' {
            if yylex.(*Lexer).PNewLine {
               yylex.(*Lexer).TokenError($1, "ambiguous syntax (proc call x new statement)")
            }
            $$ = $2
        } |
        listconstructor {
            $$ = []qsa.Expr{$1}
        } | 
        string {
            $$ = []qsa.Expr{$1}
        }

proc:
        TProc funcbody {
            $$ = &qsa.FunctionExpr{ParList:$2.ParList, Stmts: $2.Stmts}
            $$.SetLine($1.Pos.Line)
            $$.SetLastLine($2.LastLine())
        }

funcbody:
        '(' parlist ')' block TEnd {
            $$ = &qsa.FunctionExpr{ParList: $2, Stmts: $4}
            $$.SetLine($1.Pos.Line)
            $$.SetLastLine($5.Pos.Line)
        } | 
        '(' ')' block TEnd {
            $$ = &qsa.FunctionExpr{ParList: &qsa.ParList{HasVargs: false, Names: []string{}}, Stmts: $3}
            $$.SetLine($1.Pos.Line)
            $$.SetLastLine($4.Pos.Line)
        }

parlist:
        T3Comma {
            $$ = &qsa.ParList{HasVargs: true, Names: []string{}}
        } | 
        namelist {
          $$ = &qsa.ParList{HasVargs: false, Names: []string{}}
          $$.Names = append($$.Names, $1...)
        } | 
        namelist ',' T3Comma {
          $$ = &qsa.ParList{HasVargs: true, Names: []string{}}
          $$.Names = append($$.Names, $1...)
        }


listconstructor:
        '{' '}' {
            $$ = &qsa.OAListExpr{Fields: []*qsa.Field{}}
            $$.SetLine($1.Pos.Line)
        } |
        '{' fieldlist '}' {
            $$ = &qsa.OAListExpr{Fields: $2}
            $$.SetLine($1.Pos.Line)
        }


fieldlist:
        field {
            $$ = []*qsa.Field{$1}
        } | 
        fieldlist fieldsep field {
            $$ = append($1, $3)
        } | 
        fieldlist fieldsep {
            $$ = $1
        }

field:
        TIdent '=' expr {
            $$ = &qsa.Field{Key: &qsa.StringExpr{Value:$1.Str}, Value: $3}
            $$.Key.SetLine($1.Pos.Line)
        } | 
        '[' expr ']' '=' expr {
            $$ = &qsa.Field{Key: $2, Value: $5}
        } |
        expr {
            $$ = &qsa.Field{Value: $1}
        }

fieldsep:
        ',' {
            $$ = ","
        } | 
        ';' {
            $$ = ";"
        }

%%

func TokenName(c int) string {
	if c >= TAnd && c-TAnd < len(yyToknames) {
		if yyToknames[c-TAnd] != "" {
			return yyToknames[c-TAnd]
		}
	}
    return string([]byte{byte(c)})
}

