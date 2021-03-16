// package qsp q language scanner
package qsp

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"reflect"
	"strconv"
	"strings"

	"github.com/x0ray/q/qs/qsa"
)

const EOF = -1
const whitespace1 = 1<<'\t' | 1<<'\r' | 1<<' '
const whitespace2 = 1<<'\t' | 1<<'\n' | 1<<'\r' | 1<<' '

type Error struct {
	Pos     qsa.Position
	Message string
	Token   string
}

func (e *Error) Error() string {
	pos := e.Pos
	if pos.Line == EOF {
		return fmt.Sprintf("%v End input - %s\n", pos.Source, e.Message)
	} else {
		return fmt.Sprintf("%v (%d,%d) '%v'? - %s\n", pos.Source, pos.Line, pos.Column, e.Token, e.Message)
	}
}

func writeChar(buf *bytes.Buffer, c int) { buf.WriteByte(byte(c)) }

func isDecimal(ch int) bool { return '0' <= ch && ch <= '9' }

func isIdent(ch int, pos int) bool {
	return ch == '_' || 'A' <= ch && ch <= 'Z' || 'a' <= ch && ch <= 'z' || isDecimal(ch) && pos > 0
}

func isDigit(ch int) bool {
	return '0' <= ch && ch <= '9' || 'a' <= ch && ch <= 'f' || 'A' <= ch && ch <= 'F'
}

type Scanner struct {
	Pos    qsa.Position
	reader *bufio.Reader
}

func NewScanner(reader io.Reader, source string) *Scanner {
	return &Scanner{
		Pos:    qsa.Position{source, 1, 0},
		reader: bufio.NewReaderSize(reader, 4096),
	}
}

func (sc *Scanner) Error(tok string, msg string) *Error { return &Error{sc.Pos, msg, tok} }

func (sc *Scanner) TokenError(tok qsa.Token, msg string) *Error { return &Error{tok.Pos, msg, tok.Str} }

func (sc *Scanner) readNext() int {
	ch, err := sc.reader.ReadByte()
	if err == io.EOF {
		return EOF
	}
	return int(ch)
}

func (sc *Scanner) Newline(ch int) {
	if ch < 0 {
		return
	}
	sc.Pos.Line += 1
	sc.Pos.Column = 0
	next := sc.Peek()
	if ch == '\n' && next == '\r' || ch == '\r' && next == '\n' {
		sc.reader.ReadByte()
	}
}

func (sc *Scanner) Next() int {
	ch := sc.readNext()
	switch ch {
	case '\n', '\r':
		sc.Newline(ch)
		ch = int('\n')
	case EOF:
		sc.Pos.Line = EOF
		sc.Pos.Column = 0
	default:
		sc.Pos.Column++
	}
	return ch
}

func (sc *Scanner) Peek() int {
	ch := sc.readNext()
	if ch != EOF {
		sc.reader.UnreadByte()
	}
	return ch
}

func (sc *Scanner) skipWhiteSpace(whitespace int64) int {
	ch := sc.Next()
	for ; whitespace&(1<<uint(ch)) != 0; ch = sc.Next() {
	}
	return ch
}

func (sc *Scanner) skipCGoComments(ch int) error {
	// skip /* */ style comments (ch contains first '*')
	// eat contents of multi-line comment
	ch = sc.Next()
	for {
		if ch == '*' && sc.Peek() == '/' { // end of comments
			ch = sc.Next() // eat last '/'
			break
		} else {
			ch = sc.Next() // get next byte
		}
	}
	return nil
}

func (sc *Scanner) skipCppGoComments(ch int) error {
	// skip //..\n style comments (ch contains second '/')
	// eat contents of single line comment
	ch = sc.Next()
	for {
		if ch == '\n' || ch == '\r' { // end of line == end of comments
			break
		}
		ch = sc.Next() // get next byte
	}
	return nil
}

func (sc *Scanner) scanIdent(ch int, buf *bytes.Buffer) error {
	writeChar(buf, ch)
	for isIdent(sc.Peek(), 1) {
		writeChar(buf, sc.Next())
	}
	return nil
}

func (sc *Scanner) scanDecimal(ch int, buf *bytes.Buffer) error {
	writeChar(buf, ch)
	for isDecimal(sc.Peek()) {
		writeChar(buf, sc.Next())
	}
	return nil
}

func (sc *Scanner) scanNumber(ch int, buf *bytes.Buffer) error {
	if ch == '0' { // octal
		if sc.Peek() == 'x' || sc.Peek() == 'X' {
			writeChar(buf, ch)
			writeChar(buf, sc.Next())
			hasvalue := false
			for isDigit(sc.Peek()) {
				writeChar(buf, sc.Next())
				hasvalue = true
			}
			if !hasvalue {
				return sc.Error(buf.String(), "illegal hexadecimal number")
			}
			return nil
		} else if sc.Peek() != '.' && isDecimal(sc.Peek()) {
			ch = sc.Next()
		}
	}
	sc.scanDecimal(ch, buf)
	if sc.Peek() == '.' {
		sc.scanDecimal(sc.Next(), buf)
	}
	if ch = sc.Peek(); ch == 'e' || ch == 'E' {
		writeChar(buf, sc.Next())
		if ch = sc.Peek(); ch == '-' || ch == '+' {
			writeChar(buf, sc.Next())
		}
		sc.scanDecimal(sc.Next(), buf)
	}

	return nil
}

func (sc *Scanner) scanString(quote int, buf *bytes.Buffer) error {
	ch := sc.Next()
	for ch != quote {
		if ch == '\n' || ch == '\r' || ch < 0 {
			return sc.Error(buf.String(), "truncated string")
		}
		if ch == '\\' {
			if err := sc.scanEscape(ch, buf); err != nil {
				return err
			}
		} else {
			writeChar(buf, ch)
		}
		ch = sc.Next()
	}
	return nil
}

func (sc *Scanner) scanEscape(ch int, buf *bytes.Buffer) error {
	ch = sc.Next()
	switch ch {
	case 'a':
		buf.WriteByte('\a') // bell
	case 'b':
		buf.WriteByte('\b') // back space
	case 'f':
		buf.WriteByte('\f') // form feed - top of page
	case 'n':
		buf.WriteByte('\n') // new line
	case 'r':
		buf.WriteByte('\r') // carrage return
	case 't':
		buf.WriteByte('\t') // tab
	case 'v':
		buf.WriteByte('\v') //
	case '\\':
		buf.WriteByte('\\') // back slash
	case '"':
		buf.WriteByte('"') // double quote
	case '\'':
		buf.WriteByte('\'') // apostrophe
	case '\n':
		buf.WriteByte('\n') // new line
	case '\r':
		buf.WriteByte('\n') // new line
		sc.Newline('\r')
	default: // escaped numeric code point
		if '0' <= ch && ch <= '9' {
			bytes := []byte{byte(ch)}
			for i := 0; i < 2 && isDecimal(sc.Peek()); i++ {
				bytes = append(bytes, byte(sc.Next()))
			}
			val, _ := strconv.ParseInt(string(bytes), 10, 32)
			writeChar(buf, int(val))
		} else {
			buf.WriteByte('\\')
			writeChar(buf, ch)
			return sc.Error(buf.String(), "Escape sequence not correct")
		}
	}
	return nil
}

func (sc *Scanner) countSep(ch int) (int, int) {
	count := 0
	for ; ch == '='; count = count + 1 {
		ch = sc.Next()
	}
	return count, ch
}

func (sc *Scanner) scanMultilineGoString(ch int, buf *bytes.Buffer) error {
	for {
		if ch < 0 {
			return sc.Error(buf.String(), "multi-line string truncted")
		} else if ch == '`' {
			break
		}
		writeChar(buf, ch)
		ch = sc.Next()
	}
	return nil
}

var reservedWords = map[string]int{
	"and": TAnd, "break": TBreak, "do": TDo, "else": TElse, "elseif": TElseIf,
	"end": TEnd, "false": TFalse, "for": TFor, "func": TProc, "proc": TProc,
	"if": TIf, "in": TIn, "dcl": TLocal, "nil": TNil, "not": TNot, "or": TOr,
	"return": TReturn, "repeat": TRepeat, "then": TThen, "true": TTrue,
	"until": TUntil, "while": TWhile}

func (sc *Scanner) Scan(lexer *Lexer) (qsa.Token, error) {
redo:
	var err error
	tok := qsa.Token{}
	newline := false

	ch := sc.skipWhiteSpace(whitespace1)
	if ch == '\n' || ch == '\r' {
		newline = true
		ch = sc.skipWhiteSpace(whitespace2)
	}

	if ch == '(' {
		lexer.PNewLine = newline
	}

	var _buf bytes.Buffer
	buf := &_buf
	tok.Pos = sc.Pos

	switch {
	case isIdent(ch, 0):
		tok.Type = TIdent
		err = sc.scanIdent(ch, buf)
		tok.Str = buf.String()
		if err != nil {
			goto finally
		}
		if typ, ok := reservedWords[tok.Str]; ok {
			tok.Type = typ
		}
	case isDecimal(ch):
		tok.Type = TNumber
		err = sc.scanNumber(ch, buf)
		tok.Str = buf.String()
	default:
		switch ch {
		case EOF:
			tok.Type = EOF
		case '-':
			tok.Type = ch
			tok.Str = string(ch)
		case '/': // skip C, Cpp, or Go style comments
			pk := sc.Peek()
			if pk == '*' { // multi line /*..*/ comment
				err = sc.skipCGoComments(sc.Next())
				if err != nil {
					goto finally
				}
				goto redo
			} else if pk == '/' { // single line //..\n comment
				err = sc.skipCppGoComments(sc.Next())
				if err != nil {
					goto finally
				}
				goto redo
			} else {
				tok.Type = ch
				tok.Str = string(ch)
			}
		case '#': // skip Bash,sh style (first line) comments
			pk := sc.Peek()
			if pk == '!' { // single line #!..\n comment
				err = sc.skipCppGoComments(sc.Next())
				if err != nil {
					goto finally
				}
				goto redo
			} else { // allow #length unary operator
				tok.Type = ch
				tok.Str = string(ch)
			}
		case '"', '\'':
			tok.Type = TString
			err = sc.scanString(ch, buf)
			tok.Str = buf.String()
		case '`': /* long strings litterals are ` ... ` */
			tok.Type = TString
			err = sc.scanMultilineGoString(sc.Next(), buf)
			tok.Str = buf.String()
		case '[':
			tok.Type = ch
			tok.Str = string(ch)
		case '=':
			if sc.Peek() == '=' {
				tok.Type = TEqeq
				tok.Str = "=="
				sc.Next()
			} else {
				tok.Type = ch
				tok.Str = string(ch)
			}
		case '!':
			if sc.Peek() == '=' {
				tok.Type = TNeq
				tok.Str = "!="
				sc.Next()
			} else {
				err = sc.Error("!", "'!' is not valid here")
			}
		case '<':
			if sc.Peek() == '=' {
				tok.Type = TLte
				tok.Str = "<="
				sc.Next()
			} else {
				tok.Type = ch
				tok.Str = string(ch)
			}
		case '>':
			if sc.Peek() == '=' {
				tok.Type = TGte
				tok.Str = ">="
				sc.Next()
			} else {
				tok.Type = ch
				tok.Str = string(ch)
			}
		case '.':
			ch2 := sc.Peek()
			switch {
			case isDecimal(ch2):
				tok.Type = TNumber
				err = sc.scanNumber(ch, buf)
				tok.Str = buf.String()
			case ch2 == '.':
				writeChar(buf, ch)
				writeChar(buf, sc.Next())
				if sc.Peek() == '.' {
					writeChar(buf, sc.Next())
					tok.Type = T3Comma
				} else {
					tok.Type = T2Comma
				}
			default:
				tok.Type = '.'
			}
			tok.Str = buf.String()
		case '|': /* force concatenation token to be || as well as .. */
			ch2 := sc.Peek()
			switch {
			case ch2 == '|':
				writeChar(buf, '.')
				ch2 = sc.Next()
				writeChar(buf, '.')
				tok.Type = T2Comma
			default:
				tok.Type = '.'
			}
			tok.Str = buf.String()
		case '+', '*', '%', '^', '(', ')', '{', '}', ']', ';', ':', ',':
			tok.Type = ch
			tok.Str = string(ch)
		default:
			writeChar(buf, ch)
			err = sc.Error(buf.String(), "Symbol is not valid")
			goto finally
		}
	}

finally:
	tok.Name = TokenName(int(tok.Type))
	return tok, err
}

type Lexer struct {
	scanner  *Scanner
	Stmts    []qsa.Stmt
	PNewLine bool
	Token    qsa.Token
}

func (lx *Lexer) Lex(lval *yySymType) int {
	tok, err := lx.scanner.Scan(lx)
	if err != nil {
		panic(err)
	}
	if tok.Type < 0 {
		return 0
	}
	lval.token = tok
	lx.Token = tok
	return int(tok.Type)
}

func (lx *Lexer) Error(message string) {
	panic(lx.scanner.Error(lx.Token.Str, message))
}

func (lx *Lexer) TokenError(tok qsa.Token, message string) {
	panic(lx.scanner.TokenError(tok, message))
}

func Parse(reader io.Reader, name string) (segment []qsa.Stmt, err error) {
	lexer := &Lexer{NewScanner(reader, name), nil, false, qsa.Token{Str: ""}}
	segment = nil
	defer func() {
		if e := recover(); e != nil {
			err, _ = e.(error)
		}
	}()
	yyParse(lexer)
	segment = lexer.Stmts
	return
}

func isInlineDumpNode(rv reflect.Value) bool {
	switch rv.Kind() {
	case reflect.Struct, reflect.Slice, reflect.Interface, reflect.Ptr:
		return false
	default:
		return true
	}
}

func dump(node interface{}, level int, s string) string {
	rt := reflect.TypeOf(node)
	if fmt.Sprint(rt) == "<nil>" {
		return strings.Repeat(s, level) + "<nil>"
	}

	rv := reflect.ValueOf(node)
	buf := []string{}
	switch rt.Kind() {
	case reflect.Slice:
		if rv.Len() == 0 {
			return strings.Repeat(s, level) + "<empty>"
		}
		for i := 0; i < rv.Len(); i++ {
			buf = append(buf, dump(rv.Index(i).Interface(), level, s))
		}
	case reflect.Ptr:
		vt := rv.Elem()
		tt := rt.Elem()
		indicies := []int{}
		for i := 0; i < tt.NumField(); i++ {
			if strings.Index(tt.Field(i).Name, "Base") > -1 {
				continue
			}
			indicies = append(indicies, i)
		}
		switch {
		case len(indicies) == 0:
			return strings.Repeat(s, level) + "<empty>"
		case len(indicies) == 1 && isInlineDumpNode(vt.Field(indicies[0])):
			for _, i := range indicies {
				buf = append(buf, strings.Repeat(s, level)+"- Node$"+tt.Name()+": "+dump(vt.Field(i).Interface(), 0, s))
			}
		default:
			buf = append(buf, strings.Repeat(s, level)+"- Node$"+tt.Name())
			for _, i := range indicies {
				if isInlineDumpNode(vt.Field(i)) {
					inf := dump(vt.Field(i).Interface(), 0, s)
					buf = append(buf, strings.Repeat(s, level+1)+tt.Field(i).Name+": "+inf)
				} else {
					buf = append(buf, strings.Repeat(s, level+1)+tt.Field(i).Name+": ")
					buf = append(buf, dump(vt.Field(i).Interface(), level+2, s))
				}
			}
		}
	default:
		buf = append(buf, strings.Repeat(s, level)+fmt.Sprint(node))
	}
	return strings.Join(buf, "\n")
}

func Dump(segment []qsa.Stmt) string {
	return dump(segment, 0, "   ")
}
