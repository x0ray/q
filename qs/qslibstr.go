// Package qs - q scripting language
package qs

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"regexp"
	"strings"
	"unsafe"

	"github.com/x0ray/q/qs/qsm"
)

// strAfter - returns he portion of the string 'str' that follows the first
//   occurrence of 'sub' within 'str'. If 'sub' does not occur within 'str', the
//   value returned is the null string.
func strAfter(L *LState) int {
	str := L.CheckString(1)
	sub := L.CheckString(2)
	lstr := len(str)
	lsub := len(sub)
	if lsub > lstr { // substring bigger than string
		L.Push(LString(""))
		return 1
	}
	p := strings.Index(str, sub)
	if p == -1 { // substring not in string
		L.Push(LString(""))
		return 1
	}
	p = p + lsub
	if p > lstr { // no characers after substring
		L.Push(LString(""))
		return 1
	}
	s := str[p:]
	L.Push(LString(s))
	return 1
}

// strBefore - returns the portion of the string 'str' that precedes the first
//   occurrence of 'sub' within 'str'.  If 'sub' does not occur within 'str', the
//   value returned is the null string.
func strBefore(L *LState) int {
	str := L.CheckString(1)
	sub := L.CheckString(2)
	lstr := len(str)
	lsub := len(sub)
	if lsub > lstr { // substring bigger than string
		L.Push(LString(""))
		return 1
	}
	p := strings.Index(str, sub)
	if p == -1 { // substring not in string
		L.Push(LString(""))
		return 1
	}
	s := str[:p]
	L.Push(LString(s))
	return 1
}

// strByte - return list of bytes from substring of str from start to end
func strByte(L *LState) int {
	str := L.CheckString(1)
	start := L.OptInt(2, 1) - 1
	end := L.OptInt(3, -1)
	l := len(str)
	if start < 0 {
		start = l + start + 1
	}
	if end < 0 {
		end = l + end + 1
	}

	if L.GetTop() == 2 {
		if start < 0 || start >= l {
			return 0
		}
		L.Push(LNumber(str[start]))
		return 1
	}

	start = intMax(start, 0)
	end = intMin(end, l)
	if end < 0 || end <= start || start >= l {
		return 0
	}

	for i := start; i < end; i++ {
		L.Push(LNumber(str[i]))
	}
	return end - start
}

// strChar - return string from array of bytes
func strChar(L *LState) int {
	top := L.GetTop()
	bytes := make([]byte, L.GetTop())
	for i := 1; i <= top; i++ {
		bytes[i-1] = uint8(L.CheckInt(i))
	}
	L.Push(LString(string(bytes)))
	return 1
}

// strContains - contains reports whether str is within cont.
func strContains(L *LState) int {
	str := L.CheckString(1)
	cont := L.CheckString(2)
	L.Push(LBool(strings.Contains(str, cont)))
	return 1
}

// strContainsAny -  ContainsAny reports whether any Unicode code points in chars
//   are within str.
func strContainsAny(L *LState) int {
	str := L.CheckString(1)
	chars := L.CheckString(2)
	L.Push(LBool(strings.ContainsAny(str, chars)))
	return 1
}

// strCount - count counts the number of non-overlapping instances of sep in str. If
//   sep is an empty string, Count returns 1 + the number of Unicode code points
//   contained in str.
func strCount(L *LState) int {
	str := L.CheckString(1)
	sep := L.CheckString(2)
	L.Push(LNumber(strings.Count(str, sep)))
	return 1
}

// strDecodeBase64 - decode a base 64 string to a string
func strDecodeBase64(L *LState) int {
	str := L.CheckString(1)
	if strings.HasSuffix(str, "=") {
		dstr, err := base64.StdEncoding.DecodeString(str)
		if err != nil {
			L.RaiseError("base64 decode error, " + err.Error())
		} else {
			L.Push(LString(dstr))
		}
	} else {
		L.RaiseError("invalid base64 string")
	}
	return 1
}

// strDump - create a dump format string of the contents of the string parameter
func strDump(L *LState) int {
	str := L.CheckString(1)
	in := []byte(str)
	str = hex.Dump(in)
	L.Push(LString(str))
	return 1
}

// strEncodeBase64 - encode a string to a base 64 string
func strEncodeBase64(L *LState) int {
	str := L.CheckString(1)
	b := []byte(str) // make string byte slice
	str = base64.StdEncoding.EncodeToString(b)
	L.Push(LString(str))
	return 1
}

// strEscapeXmlData - returns str with characters needing to be escaped for XML
//   repleced with the the correct XML escape sequence. For example: < --> &lt;
func strEscapeXmlData(L *LState) int {
	str := L.CheckString(1)
	L.Push(LString(EscapeXmlData(str)))
	return 1
}

func strFind(L *LState) int {
	str := L.CheckString(1)
	pattern := L.CheckString(2)
	if len(pattern) == 0 {
		L.Push(LNumber(1))
		L.Push(LNumber(0))
		return 2
	}
	init := oaIndex2StringIndex(str, L.OptInt(3, 1), true)
	plain := false
	if L.GetTop() == 4 {
		plain = LVAsBool(L.Get(4))
	}

	if plain {
		pos := strings.Index(str[init:], pattern)
		if pos < 0 {
			L.Push(LNil)
			return 1
		}
		L.Push(LNumber(init+pos) + 1)
		L.Push(LNumber(init + pos + len(pattern)))
		return 2
	}

	mds, err := qsm.Find(pattern, *(*[]byte)(unsafe.Pointer(&str)), init, 1)
	if err != nil {
		L.RaiseError(err.Error())
	}
	if len(mds) == 0 {
		L.Push(LNil)
		return 1
	}
	md := mds[0]
	L.Push(LNumber(md.Capture(0) + 1))
	L.Push(LNumber(md.Capture(1)))
	for i := 2; i < md.CaptureLength(); i += 2 {
		if md.IsPosCapture(i) {
			L.Push(LNumber(md.Capture(i)))
		} else {
			L.Push(LString(str[md.Capture(i):md.Capture(i+1)]))
		}
	}
	return md.CaptureLength()/2 + 1
}

func strFormat(L *LState) int {
	str := L.CheckString(1)
	args := make([]interface{}, L.GetTop()-1)
	top := L.GetTop()
	for i := 2; i <= top; i++ {
		args[i-2] = L.Get(i)
	}
	npat := strings.Count(str, "%") - strings.Count(str, "%%")
	L.Push(LString(fmt.Sprintf(str, args[:intMin(npat, len(args))]...)))
	return 1
}

func strGsub(L *LState) int {
	str := L.CheckString(1)
	pat := L.CheckString(2)
	L.CheckTypes(3, LTString, LTOAList, LTProc)
	repl := L.CheckAny(3)
	limit := L.OptInt(4, -1)

	mds, err := qsm.Find(pat, *(*[]byte)(unsafe.Pointer(&str)), 0, limit)
	if err != nil {
		L.RaiseError(err.Error())
	}
	if len(mds) == 0 {
		L.SetTop(1)
		L.Push(LNumber(0))
		return 2
	}
	switch lv := repl.(type) {
	case LString:
		L.Push(LString(strGsubStr(L, str, string(lv), mds)))
	case *LOAList:
		L.Push(LString(strGsubOAList(L, str, lv, mds)))
	case *LProc:
		L.Push(LString(strGsubFunc(L, str, lv, mds)))
	}
	L.Push(LNumber(len(mds)))
	return 2
}

type replaceInfo struct {
	Indicies []int
	String   string
}

func checkCaptureIndex(L *LState, m *qsm.MatchData, idx int) {
	if idx <= 2 {
		return
	}
	if idx >= m.CaptureLength() {
		L.RaiseError("invalid capture index")
	}
}

func capturedString(L *LState, m *qsm.MatchData, str string, idx int) string {
	checkCaptureIndex(L, m, idx)
	if idx >= m.CaptureLength() && idx == 2 {
		idx = 0
	}
	if m.IsPosCapture(idx) {
		return fmt.Sprint(m.Capture(idx))
	} else {
		return str[m.Capture(idx):m.Capture(idx+1)]
	}

}

func strGsubDoReplace(str string, info []replaceInfo) string {
	offset := 0
	buf := []byte(str)
	for _, replace := range info {
		oldlen := len(buf)
		b1 := append([]byte(""), buf[0:offset+replace.Indicies[0]]...)
		b2 := []byte("")
		index2 := offset + replace.Indicies[1]
		if index2 <= len(buf) {
			b2 = append(b2, buf[index2:len(buf)]...)
		}
		buf = append(b1, replace.String...)
		buf = append(buf, b2...)
		offset += len(buf) - oldlen
	}
	return string(buf)
}

func strGsubStr(L *LState, str string, repl string, matches []*qsm.MatchData) string {
	infoList := make([]replaceInfo, 0, len(matches))
	for _, match := range matches {
		start, end := match.Capture(0), match.Capture(1)
		sc := newFlagScanner('%', "", "", repl)
		for c, eos := sc.Next(); !eos; c, eos = sc.Next() {
			if !sc.ChangeFlag {
				if sc.HasFlag {
					if c >= '0' && c <= '9' {
						sc.AppendString(capturedString(L, match, str, 2*(int(c)-48)))
					} else {
						sc.AppendChar('%')
						sc.AppendChar(c)
					}
					sc.HasFlag = false
				} else {
					sc.AppendChar(c)
				}
			}
		}
		infoList = append(infoList, replaceInfo{[]int{start, end}, sc.String()})
	}

	return strGsubDoReplace(str, infoList)
}

func strGsubOAList(L *LState, str string, repl *LOAList, matches []*qsm.MatchData) string {
	infoList := make([]replaceInfo, 0, len(matches))
	for _, match := range matches {
		idx := 0
		if match.CaptureLength() > 2 { // has captures
			idx = 2
		}
		var value LValue
		if match.IsPosCapture(idx) {
			value = L.GetOAList(repl, LNumber(match.Capture(idx)))
		} else {
			value = L.GetField(repl, str[match.Capture(idx):match.Capture(idx+1)])
		}
		if !LVIsFalse(value) {
			infoList = append(infoList, replaceInfo{[]int{match.Capture(0), match.Capture(1)}, LVAsString(value)})
		}
	}
	return strGsubDoReplace(str, infoList)
}

func strGsubFunc(L *LState, str string, repl *LProc, matches []*qsm.MatchData) string {
	infoList := make([]replaceInfo, 0, len(matches))
	for _, match := range matches {
		start, end := match.Capture(0), match.Capture(1)
		L.Push(repl)
		nargs := 0
		if match.CaptureLength() > 2 { // has captures
			for i := 2; i < match.CaptureLength(); i += 2 {
				if match.IsPosCapture(i) {
					L.Push(LNumber(match.Capture(i)))
				} else {
					L.Push(LString(capturedString(L, match, str, i)))
				}
				nargs++
			}
		} else {
			L.Push(LString(capturedString(L, match, str, 0)))
			nargs++
		}
		L.Call(nargs, 1)
		ret := L.reg.Pop()
		if !LVIsFalse(ret) {
			infoList = append(infoList, replaceInfo{[]int{start, end}, LVAsString(ret)})
		}
	}
	return strGsubDoReplace(str, infoList)
}

type strMatchData struct {
	str     string
	pos     int
	matches []*qsm.MatchData
}

func strGmatchIter(L *LState) int {
	md := L.CheckUserData(1).Value.(*strMatchData)
	str := md.str
	matches := md.matches
	idx := md.pos
	md.pos += 1
	if idx == len(matches) {
		return 0
	}
	L.Push(L.Get(1))
	match := matches[idx]
	if match.CaptureLength() == 2 {
		L.Push(LString(str[match.Capture(0):match.Capture(1)]))
		return 1
	}

	for i := 2; i < match.CaptureLength(); i += 2 {
		if match.IsPosCapture(i) {
			L.Push(LNumber(match.Capture(i)))
		} else {
			L.Push(LString(str[match.Capture(i):match.Capture(i+1)]))
		}
	}
	return match.CaptureLength()/2 - 1
}

func strGmatch(L *LState) int {
	str := L.CheckString(1)
	pattern := L.CheckString(2)
	mds, err := qsm.Find(pattern, []byte(str), 0, -1)
	if err != nil {
		L.RaiseError(err.Error())
	}
	L.Push(L.Get(UpvalueIndex(1)))
	ud := L.NewUserData()
	ud.Value = &strMatchData{str, 0, mds}
	L.Push(ud)
	return 2
}

// strHasPrefix - tests whether the string str begins with string prefix.
func strHasPrefix(L *LState) int {
	str := L.CheckString(1)
	prefix := L.CheckString(2)
	L.Push(LBool(strings.HasPrefix(str, prefix)))
	return 1
}

// strHasSuffix - tests whether the string str ends with string prefix.
func strHasSuffix(L *LState) int {
	str := L.CheckString(1)
	suffix := L.CheckString(2)
	L.Push(LBool(strings.HasSuffix(str, suffix)))
	return 1
}

// strIndex - index returns the index of the first instance of sep in str, or -1 if
//   sep is not present in str.
func strIndex(L *LState) int {
	str := L.CheckString(1)
	sep := L.CheckString(2)
	L.Push(LNumber(strings.Index(str, sep)))
	return 1
}

// strIndexAny - indexAny returns the index of the first instance of any Unicode
//   code point from chars in s, or -1 if no Unicode code point from chars is
//   present in str
func strIndexAny(L *LState) int {
	str := L.CheckString(1)
	chars := L.CheckString(2)
	L.Push(LNumber(strings.IndexAny(str, chars)))
	return 1
}

// strIsName - returns true if the string is a name
//   Names consist of one alphabetic character or underscore followed by
//   1 to 31 alphanumeric characters that may include underscore
func strIsName(L *LState) int {
	str := L.CheckString(1)
	L.Push(LBool(IsName(str)))
	return 1
}

// strIsXmlTagName - returns true if the string is a valid XML tag name
//   Tag names cannot contain any of the characters !"#$%&'()*+,/;<=>?@[\]^`{|}~,
//   nor a space character, and cannot begin with "-", ".", or a numeric digit.
func strIsXmlTagName(L *LState) int {
	str := L.CheckString(1)
	L.Push(LBool(IsXmlTagName(str)))
	return 1
}

// strLastIndex - lastIndex returns the index of the last instance of sep in str,
//   or -1 if sep is not present in str.
func strLastIndex(L *LState) int {
	str := L.CheckString(1)
	sep := L.CheckString(2)
	L.Push(LNumber(strings.LastIndex(str, sep)))
	return 1
}

// strLastIndexAny - LastIndexAny returns the index of the last instance of
//   any Unicode code point from chars in str, or -1 if no Unicode code point
//   from chars is present in str.
func strLastIndexAny(L *LState) int {
	str := L.CheckString(1)
	chars := L.CheckString(2)
	L.Push(LNumber(strings.LastIndexAny(str, chars)))
	return 1
}

//  strLen - returns the langth of the string
func strLen(L *LState) int {
	str := L.CheckString(1)
	L.Push(LNumber(len(str)))
	return 1
}

// strLower - converts string str to lower case
func strLower(L *LState) int {
	str := L.CheckString(1)
	L.Push(LString(strings.ToLower(str)))
	return 1
}

// strMakeXmlTagName - returns a valid XML tag based on an input string. The
//   input string may have invalid characters.
//   Valid XML tag names consist of a string of alphanumeric characters that can
//   also include _ - . characters. The first character must not be numeric or
//   contain the . or - character
//   Either an expansion or contraction method can be selected to process
//   invalid characters automatically.
//     c - Remove invalid characters
//     r - Replace invalid characters with supplied value
//     x - Expand invalid characters - NB: non reversable
//     e - Do nothing with invalid characters and return an error
func strMakeXmlTagName(L *LState) int {
	str := L.CheckString(1)
	opt := L.OptString(2, "c")
	rep := L.OptString(3, ".")
	prefix := L.OptString(4, "")
	tag, err := MakeXmlTagName(str, opt, rep, prefix)
	if err != nil {
		L.RaiseError(err.Error())
	}
	L.Push(LString(tag))
	return 1
}

func strMatch(L *LState) int {
	str := L.CheckString(1)
	pattern := L.CheckString(2)
	offset := L.OptInt(3, 1)
	l := len(str)
	if offset < 0 {
		offset = l + offset + 1
	}
	offset--
	if offset < 0 {
		offset = 0
	}

	mds, err := qsm.Find(pattern, *(*[]byte)(unsafe.Pointer(&str)), offset, 1)
	if err != nil {
		L.RaiseError(err.Error())
	}
	if len(mds) == 0 {
		L.Push(LNil)
		return 0
	}
	md := mds[0]
	nsubs := md.CaptureLength() / 2
	switch nsubs {
	case 1:
		L.Push(LString(str[md.Capture(0):md.Capture(1)]))
		return 1
	default:
		for i := 2; i < md.CaptureLength(); i += 2 {
			if md.IsPosCapture(i) {
				L.Push(LNumber(md.Capture(i)))
			} else {
				L.Push(LString(str[md.Capture(i):md.Capture(i+1)]))
			}
		}
		return nsubs - 1
	}
}

// strPrxMatch - return position in 'str' of regular expression 'rx' match
func strPrxMatch(L *LState) int {
	rx := L.CheckString(1)
	str := L.CheckString(2)

	re, err := regexp.Compile(rx)
	if err != nil {
		L.Push(LNil)
		return 0
	}

	ixs := re.FindAllStringIndex(str, -1)
	if len(ixs) > 0 {
		ret := L.NewOAList()
		for i, v := range ixs { // copy slice number values to OA list
			// NB: i+1 converst 0 based go to 1 based OA list
			ret.RawSetInt(i+1, LNumber(v[0]))
		}
		L.Push(ret)
		return 1
	}
	L.Push(LNil)
	return 0
}

// strPrxChange - return 'str' with regular expression 'rx' matches replaced by 'rep'
func strPrxChange(L *LState) int {
	rx := L.CheckString(1)
	str := L.CheckString(2)
	rep := L.CheckString(3)

	re, err := regexp.Compile(rx)
	if err != nil {
		L.Push(LNil)
		return 0
	}

	s := re.ReplaceAllString(str, rep)
	L.Push(LString(s))
	return 1
}

// strRep - repeat the input str string n times, concatenating the repeated string as the output string
func strRep(L *LState) int {
	str := L.CheckString(1)
	n := L.CheckInt(2)
	L.Push(LString(strings.Repeat(str, n)))
	return 1
}

// strReplace - replace n occurrences of olds string with news string in str
func strReplace(L *LState) int {
	str := L.CheckString(1)
	olds := L.CheckString(2)
	news := L.CheckString(3)
	n := L.CheckInt(4)
	str = strings.Replace(str, olds, news, n)
	L.Push(LString(str))
	return 1
}

// strReverse - reverse all characters in the string
func strReverse(L *LState) int {
	str := L.CheckString(1)
	bts := []byte(str)
	out := make([]byte, len(bts))
	for i, j := 0, len(bts)-1; j >= 0; i, j = i+1, j-1 {
		out[i] = bts[j]
	}
	L.Push(LString(string(out)))
	return 1
}

// strScan - returns substring 'n' of 'str' delimited by characters from 'delimstr'
func strScan(L *LState) int {
	str := L.CheckString(1)
	delim := L.CheckString(2)
	pos := L.CheckInt(3)
	s := Scan(pos, str, delim)
	L.Push(LString(s))
	return 1
}

// strScanAll - returns list of 'str' delimited by characters from 'delimstr'
func strScanAll(L *LState) int {
	str := L.CheckString(1)
	delim := L.CheckString(2)
	sl := ScanAll(str, delim) // scan text into slice
	ret := L.NewOAList()
	for i, v := range sl { // copy slice values to OA list
		// NB: i+1 converst 0 based go to 1 based OA list
		ret.RawSetInt(i+1, LString(v))
	}
	L.Push(ret)
	return 1
}

// strSub - returns a substring of str from start to end
func strSub(L *LState) int {
	str := L.CheckString(1)
	start := oaIndex2StringIndex(str, L.CheckInt(2), true)
	end := oaIndex2StringIndex(str, L.OptInt(3, -1), false)
	l := len(str)
	if start >= l || end < start {
		L.Push(LString(""))
	} else {
		L.Push(LString(str[start:end]))
	}
	return 1
}

// strSubstr - returns a substring of str from start to end
func strSubstr(L *LState) int {
	str := L.CheckString(1)
	start := oaIndex2StringIndex(str, L.CheckInt(2), true)
	length := oaIndex2StringIndex(str, L.OptInt(3, -1), false)
	l := len(str)
	if start >= l || start < 0 {
		L.Push(LString(""))
	} else if length == -1 || (start+length) >= l {
		L.Push(LString(str[start:]))
	} else {
		L.Push(LString(str[start : start+length]))
	}
	return 1
}

// strTitle - title returns a copy of the string str with all Unicode letters
//   that begin words mapped to their title case.
func strTitle(L *LState) int {
	str := L.CheckString(1)
	L.Push(LString(strings.Title(str)))
	return 1
}

// strTrim - trim returns a slice of the string str with all leading and
//   trailing Unicode code points contained in cutset removed.
func strTrim(L *LState) int {
	str := L.CheckString(1)
	cutset := L.CheckString(2)
	L.Push(LString(strings.Trim(str, cutset)))
	return 1
}

// strTrimLeft - trimLeft returns a slice of the string str with all leading
//   Unicode code points contained in cutset removed.
func strTrimLeft(L *LState) int {
	str := L.CheckString(1)
	cutset := L.CheckString(2)
	L.Push(LString(strings.TrimLeft(str, cutset)))
	return 1
}

// strTrimPrefix - trimPrefix returns str without the provided leading
//   prefix string. If str doesn't start with prefix, str is returned unchanged.
func strTrimPrefix(L *LState) int {
	str := L.CheckString(1)
	prefix := L.CheckString(2)
	L.Push(LString(strings.TrimPrefix(str, prefix)))
	return 1
}

// strTrimRight - trimRight returns a slice of the string str, with all
//   trailing Unicode code points contained in cutset removed.
func strTrimRight(L *LState) int {
	str := L.CheckString(1)
	cutset := L.CheckString(2)
	L.Push(LString(strings.TrimRight(str, cutset)))
	return 1
}

// strTrimSpace - trimSpace returns a slice of the string str, with
//   all leading and trailing white space removed, as defined by Unicode.
func strTrimSpace(L *LState) int {
	str := L.CheckString(1)
	L.Push(LString(strings.TrimSpace(str)))
	return 1
}

// strTrimSuffix - trimSuffix returns str without the provided trailing
//   suffix string. If str doesn't end with suffix, str is returned unchanged.
func strTrimSuffix(L *LState) int {
	str := L.CheckString(1)
	suffix := L.CheckString(2)
	L.Push(LString(strings.TrimSuffix(str, suffix)))
	return 1
}

// strUnEscapeXmlData - returns str with XML escape sequences repleced with the
//   characters that were escaped. For example: &lt; --> <
func strUnEscapeXmlData(L *LState) int {
	str := L.CheckString(1)
	L.Push(LString(UnEscapeXmlData(str)))
	return 1
}

// strUpper - converts string str to upper case
func strUpper(L *LState) int {
	str := L.CheckString(1)
	L.Push(LString(strings.ToUpper(str)))
	return 1
}

// oaIndex2StringIndex -
func oaIndex2StringIndex(str string, i int, start bool) int {
	if start && i != 0 {
		i -= 1
	}
	l := len(str)
	if i < 0 {
		i = l + i + 1
	}
	i = intMax(0, i)
	if !start && i > l {
		i = l
	}
	return i
}

//-----------------------------------------------------------------------------------
// Support functions
//-----------------------------------------------------------------------------------

func iabs(x int) int {
	switch {
	case x >= 0:
		return x
	case x > MinInt:
		return -x
	}
	panic("iabs: invalid argument")
}

func pos(c byte, s string) int {
	for i, v := range s {
		if c == byte(v) {
			return i
		}
	} // of for
	return -1
} /* of pos() */

func Scan(n int, InText string, Delim string) string {
	/*
		The scan function will return the nth word from the input string
		InText whereby the returned word is delimited by characters from the
		string Delim.

		If InText or Delim are empty, or n=0 the result will be an empty string.

		If n is negative the scan will start at tbe tail of the string and
		proceed to the front of the string until the nth word is located.

		If the nth word is not located an empty string will be returned.
	*/
	result := ""
	if (n == 0) || (InText == "") || (Delim == "") { // invalid parameter
		return result
	}
	cwdStart := 0
	cwdEnd := 0
	InWord := false
	WordCnt := 0

	Reverse := false // scan is forwards
	InTextLen := len(InText) - 1
	i := 0
	if n < 0 { // scan is backwards
		n = iabs(n)
		Reverse = true
		i = InTextLen
		InTextLen = 0
	}

	//fmt.Printf("InTextLen=%d i=%d\n", InTextLen, i)
	for i != InTextLen { // scan characters of the input text
		//fmt.Printf("cwdStart=%d cwdEnd=%d InText=%s i=%d \n", cwdStart, cwdEnd, InText, i)
		ch := InText[i]
		if pos(ch, Delim) < 0 { // Character from word
			if !InWord {
				cwdStart = i
				cwdEnd = i
			}
			InWord = true
		} else { // Character is a delimiter
			if InWord { // exiting a word
				InWord = false
				cwdEnd = i
				WordCnt++
				if n == WordCnt {
					break
				} else {
					cwdStart = i
					cwdEnd = i
				}
			}
		}
		if Reverse {
			i--
		} else {
			i++
		}
	} // of for

	if InWord { // exiting a word
		cwdEnd = i
		WordCnt++
	}

	if n == WordCnt { // got nth word ?
		if Reverse {
			cwdEnd++
			cwdStart++
			// NB positions have reversed meanings in next stmt.
			result = InText[cwdEnd:cwdStart]
		} else {
			result = InText[cwdStart:cwdEnd]
		}
	}
	return result
} /* of Scan() */

func ScanAll(InText string, Delim string) []string {
	/*
		The ScanAll function will return all the words from the input string
		whereby the returned words are delimited by one or more characters from the
		string Delim.

		If InText or Delim are empty, the result will be an empty string.
	*/

	result := []string{}
	if (InText == "") || (Delim == "") { // invalid parameter
		return result
	}
	word := ""
	cwdStart := 0
	cwdEnd := 0
	InWord := false
	WordCnt := 0
	InTextLen := len(InText)
	i := 0

	//fmt.Printf("InTextLen=%d \n", InTextLen)
	for i != InTextLen { // scan characters of the input text
		//fmt.Printf("cwdStart=%d cwdEnd=%d InText=%s i=%d \n", cwdStart, cwdEnd, InText, i)
		ch := InText[i]
		if pos(ch, Delim) < 0 { // Character from word
			if !InWord {
				cwdStart = i
				cwdEnd = i
			}
			InWord = true
		} else { // Character is a delimiter
			if InWord { // exiting a word
				InWord = false
				cwdEnd = i
				WordCnt++
				word = InText[cwdStart:cwdEnd]
				//fmt.Printf("-- word=%s WordCnt=%d cwdStart=%d cwdEnd=%d \n", word, WordCnt, cwdStart, cwdEnd)
				result = append(result, word)
			}
		}
		i++
	} // of for
	if InWord { // exiting a word
		InWord = false
		cwdEnd = i
		WordCnt++
		word = InText[cwdStart:cwdEnd]
		result = append(result, word)
	}

	return result
} /* of ScanAll() */

func MultiScan(InText string, Delim ...string) []string {
	/*
		The MultiScan function will return all the words from the input string
		whereby the ith word is delimited by one or more characters from the
		ith string of Delim.

		If InText or Delim are empty, the result will be an empty string.

		If there are more words in InText than strings in Delim then the last
		string of delimiters in Delim will be used for all the remaining words
		in InText
	*/

	result := []string{}
	DelimLen := len(Delim)
	if (InText == "") || (DelimLen == 0) { // invalid parameter
		return result
	}
	word := ""
	cwdStart := 0
	cwdEnd := 0
	InWord := false
	WordCnt := 0
	DelimCnt := 0
	InTextLen := len(InText) - 1
	i := 0

	//fmt.Printf("InTextLen=%d DelimLen=%d \n", InTextLen, DelimLen)
	for i != InTextLen { // scan characters of the input text
		//fmt.Printf("cwdStart=%d cwdEnd=%d InText=%s i=%d DelimCnt=%d\n", cwdStart, cwdEnd, InText, i, DelimCnt)
		ch := InText[i]
		if pos(ch, Delim[DelimCnt]) < 0 { // Character from word
			if !InWord {
				cwdStart = i
				cwdEnd = i
			}
			InWord = true
		} else { // Character is a delimiter
			if InWord { // exiting a word
				InWord = false
				cwdEnd = i
				WordCnt++
				word = InText[cwdStart:cwdEnd]
				//fmt.Printf("-- word=%s WordCnt=%d cwdStart=%d cwdEnd=%d \n", word, WordCnt, cwdStart, cwdEnd)
				result = append(result, word)
				if DelimCnt < (DelimLen - 1) {
					DelimCnt++
				}
			}
		}
		i++
	} // of for

	return result
} /* of MultiScan() */
