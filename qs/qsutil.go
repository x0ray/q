// Package qs - q scripting language
package qs

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// intMin - return mimimum of two integers
func intMin(a, b int) int {
	if a < b {
		return a
	} else {
		return b
	}
}

// intMax - return maximum of two integers
func intMax(a, b int) int {
	if a > b {
		return a
	} else {
		return b
	}
}

// defaultFormat -
func defaultFormat(v interface{}, f fmt.State, c rune) {
	buf := make([]string, 0, 10)
	buf = append(buf, "%")
	for i := 0; i < 128; i++ {
		if f.Flag(i) {
			buf = append(buf, string(i))
		}
	}

	if w, ok := f.Width(); ok {
		buf = append(buf, strconv.Itoa(w))
	}
	if p, ok := f.Precision(); ok {
		buf = append(buf, "."+strconv.Itoa(p))
	}
	buf = append(buf, string(c))
	format := strings.Join(buf, "")
	fmt.Fprintf(f, format, v)
}

type flagScanner struct {
	flag       byte
	start      string
	end        string
	buf        []byte
	str        string
	Length     int
	Pos        int
	HasFlag    bool
	ChangeFlag bool
}

func newFlagScanner(flag byte, start, end, str string) *flagScanner {
	return &flagScanner{flag, start, end, make([]byte, 0, len(str)), str, len(str), 0, false, false}
}

func (fs *flagScanner) AppendString(str string) { fs.buf = append(fs.buf, str...) }

func (fs *flagScanner) AppendChar(ch byte) { fs.buf = append(fs.buf, ch) }

func (fs *flagScanner) String() string { return string(fs.buf) }

func (fs *flagScanner) Next() (byte, bool) {
	c := byte('\000')
	fs.ChangeFlag = false
	if fs.Pos == fs.Length {
		if fs.HasFlag {
			fs.AppendString(fs.end)
		}
		return c, true
	} else {
		c = fs.str[fs.Pos]
		if c == fs.flag {
			if fs.Pos < (fs.Length-1) && fs.str[fs.Pos+1] == fs.flag {
				fs.HasFlag = false
				fs.AppendChar(fs.flag)
				fs.Pos += 2
				return fs.Next()
			} else if fs.Pos != fs.Length-1 {
				if fs.HasFlag {
					fs.AppendString(fs.end)
				}
				fs.AppendString(fs.start)
				fs.ChangeFlag = true
				fs.HasFlag = true
			}
		}
	}
	fs.Pos++
	return c, false
}

// cDateFlagToGo - maps a C date flag onto a Golang date format
var cDateFlagToGo = map[byte]string{
	'a': "mon", 'A': "Monday", 'b': "Jan", 'B': "January", 'c': "02 Jan 06 15:04 MST", 'd': "02",
	'F': "2006-01-02", 'H': "15", 'I': "03", 'm': "01", 'M': "04", 'p': "PM", 'P': "pm", 'S': "05",
	'x': "15/04/05", 'X': "15:04:05", 'y': "06", 'Y': "2006", 'z': "-0700", 'Z': "MST"}

// strftime - returns time t as a string using a C style format
func strftime(t time.Time, cfmt string) string {
	sc := newFlagScanner('%', "", "", cfmt)
	for c, eos := sc.Next(); !eos; c, eos = sc.Next() {
		if !sc.ChangeFlag {
			if sc.HasFlag {
				if v, ok := cDateFlagToGo[c]; ok {
					sc.AppendString(t.Format(v))
				} else {
					switch c {
					case 'w':
						sc.AppendString(fmt.Sprint(int(t.Weekday())))
					default:
						sc.AppendChar('%')
						sc.AppendChar(c)
					}
				}
				sc.HasFlag = false
			} else {
				sc.AppendChar(c)
			}
		}
	}

	return sc.String()
}

// isInteger - returns true if v is an integer else false
func isInteger(v LNumber) bool {
	return float64(v) == float64(int64(v))
}

// isArrayKey - returns true if v is an array key
func isArrayKey(v LNumber) bool {
	return isInteger(v) && v < LNumber(int((^uint(0))>>1)) && v > LNumber(0) && v < LNumber(MaxArrayIndex)
}

// parseNumber - converts string to a number
func parseNumber(number string) (LNumber, error) {
	var value LNumber
	number = strings.Trim(number, " \t\n")
	if v, err := strconv.ParseInt(number, 0, LNumberBit); err != nil {
		if v2, err2 := strconv.ParseFloat(number, LNumberBit); err2 != nil {
			return LNumber(0), err2
		} else {
			value = LNumber(v2)
		}
	} else {
		value = LNumber(v)
	}
	return value, nil
}

// popenArgs - formats args for use by systems shell
func popenArgs(arg string) (string, []string) {
	cmd := "/bin/sh"
	args := []string{"-c"}
	if QsOS == "windows" {
		cmd = "C:\\Windows\\system32\\cmd.exe"
		args = []string{"/c"}
	}
	args = append(args, arg)
	return cmd, args
}

// isGoroutineSafe - determines if value is safe for use in Go routine
func isGoroutineSafe(lv LValue) bool {
	switch v := lv.(type) {
	case *LProc, *LUserData, *LState:
		return false
	case *LOAList:
		return v.Metalist == LNil
	default:
		return true
	}
}

// readBufioSize - reads reader by size returning slice of bytes
func readBufioSize(reader *bufio.Reader, size int64) ([]byte, error, bool) {
	result := []byte{}
	read := int64(0)
	var err error
	var n int
	for read != size {
		buf := make([]byte, size-read)
		n, err = reader.Read(buf)
		if err != nil {
			break
		}
		read += int64(n)
		result = append(result, buf[:n]...)
	}
	e := err
	if e != nil && e == io.EOF {
		e = nil
	}

	return result, e, len(result) == 0 && err == io.EOF
}

// readBufioLine - reads reader into slice of bytes
func readBufioLine(reader *bufio.Reader) ([]byte, error, bool) {
	result := []byte{}
	var buf []byte
	var err error
	var isprefix bool = true
	for isprefix {
		buf, isprefix, err = reader.ReadLine()
		if err != nil {
			break
		}
		result = append(result, buf...)
	}
	e := err
	if e != nil && e == io.EOF {
		e = nil
	}

	return result, e, len(result) == 0 && err == io.EOF
}

// int2Fb - converts integer
func int2Fb(val int) int {
	e := 0
	x := val
	for x >= 16 {
		x = (x + 1) >> 1
		e++
	}
	if x < 8 {
		return x
	}
	return ((e + 1) << 3) | (x - 8)
}

// strCmp - compares two strings
func strCmp(s1, s2 string) int {
	len1 := len(s1)
	len2 := len(s2)
	for i := 0; ; i++ {
		c1 := -1
		if i < len1 {
			c1 = int(s1[i])
		}
		c2 := -1
		if i != len2 {
			c2 = int(s2[i])
		}
		switch {
		case c1 < c2:
			return -1
		case c1 > c2:
			return +1
		case c1 < 0:
			return 0
		}
	}
}

const tagsel = `([^a-zA-Z0-9_.-]*)` // bad tag char selector
var rxtag *regexp.Regexp = regexp.MustCompile(tagsel)

const tagself = `(^[^a-zA-Z_]*)` // bad first tag char selector
var rxtagf *regexp.Regexp = regexp.MustCompile(tagself)

// makeXmlTagName - Takes any string of characters and makes them into a valid
//   XML tag name (not - . digit or space, followed by word-characters _ . -)
//   Either an expansion or contraction method can be selected to process
//   invalid characters automatically.
//     c - Remove invalid characters
//     r - Replace invalid characters with supplied value
//     x - Expand invalid characters - NB: non reversable
//     e - Do nothing with invalid characters and return an error
func MakeXmlTagName(str string, opt string, rep string, prefix string) (string, error) {
	var err error
	var errf bool
	// remove characters not allowed any where in name
	str = rxtag.ReplaceAllStringFunc(str, func(s string) string {
		if opt == "c" {
			return ""
		} else if opt == "r" {
			return rep
		} else if opt == "x" {
			return fmt.Sprintf("%x", s)
		} else {
			errf = true
			return s
		}
	})
	// remove characters not allowed at the front of the name
	str = rxtagf.ReplaceAllStringFunc(str, func(s string) string {
		if opt == "c" {
			return ""
		} else if opt == "r" {
			if rep == "." || rep == "-" {
				return ""
			} else {
				return rep
			}
		} else if opt == "x" {
			if s != "" {
				return fmt.Sprintf("X%x", s)
			} else {
				return ""
			}
		} else {
			errf = true
			return s
		}
	})
	if errf {
		err = errors.New("Invalid characters found in tag name")
	}
	str = prefix + str
	return str, err
}

// IsName - tests if a string contains a name
func IsName(s string) bool {
	// Handy regex tool: http://www.regexr.com/
	// Names consist of one alphabetic character or underscore followed by
	//   1 to 31 alphanumeric characters that may include underscore
	rs := `^[A-Za-z_]\w{1,32}$`
	re := regexp.MustCompile(rs)
	matched := re.MatchString(s)
	return matched
}

const (
	nonXmlTagChars     = ` !"#$%&'()*+,/;<=>?@[\]^{|}~` + "`"
	nonXmlTagFirstChar = `-.0123456789`
)

// IsXmlTagName - tests if a string contains a valid XML tag
func IsXmlTagName(s string) bool {
	// Tag names cannot contain any of the characters !"#$%&'()*+,/;<=>?@[\]^`{|}~,
	// nor a space character, and cannot begin with "-", ".", or a numeric digit.
	// ref: https://en.wikipedia.org/wiki/XML
	//      https://www.w3.org/TR/2008/REC-xml-20081126/
	ok := !strings.ContainsAny(s, nonXmlTagChars)
	if ok {
		ok = !strings.ContainsAny(s[:1], nonXmlTagFirstChar)
	}
	return ok
}

// EscapeXmlDataName - replaces characters &<>'" with their XML escape sequence
func EscapeXmlData(v string) string {
	// TODO optimize - too many passes of string
	v = strings.Replace(v, "&", "&amp;", -1)
	v = strings.Replace(v, "<", "&lt;", -1)
	v = strings.Replace(v, ">", "&gt;", -1)
	v = strings.Replace(v, "'", "&apos;", -1)
	v = strings.Replace(v, `"`, "&quot;", -1)
	// TODO ? following ?
	v = strings.Replace(v, "\n", "&#xA;", -1)
	v = strings.Replace(v, "\t", "&#x9;", -1)
	v = strings.Replace(v, "\r", "&#xD;", -1)
	return v
}

// UnEscapeXmlDataName - replaces XML escape sequences with their &<>'" character
func UnEscapeXmlData(v string) string {
	// TODO optimize - too many passes of string
	v = strings.Replace(v, "&quot;", `"`, -1)
	v = strings.Replace(v, "&apos;", "'", -1)
	v = strings.Replace(v, "&gt;", ">", -1)
	v = strings.Replace(v, "&lt;", "<", -1)
	v = strings.Replace(v, "&amp;", "&", -1)
	// TODO ? following ?
	v = strings.Replace(v, "&#xA;", "\n", -1)
	v = strings.Replace(v, "&#x9;", "\t", -1)
	v = strings.Replace(v, "&#xD;", "\r", -1)
	return v
}
