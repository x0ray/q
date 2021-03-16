// Package qs - q scripting language
package qs

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"
)

func listSort(L *LState) int {
	list := L.CheckOAList(1)
	sorter := lValueArraySorter{L, nil, list.array}
	if L.GetTop() != 1 {
		sorter.Fn = L.CheckProc(2)
	}
	sort.Sort(sorter)
	return 0
}

func listGetN(L *LState) int {
	L.Push(LNumber(L.CheckOAList(1).Len()))
	return 1
}

func listMaxN(L *LState) int {
	L.Push(LNumber(L.CheckOAList(1).MaxN()))
	return 1
}

func listErase(L *LState) int {
	list := L.CheckOAList(1)
	if L.GetTop() == 1 {
		L.Push(list.Remove(-1))
	} else {
		L.Push(list.Remove(L.CheckInt(2)))
	}
	return 1
}

func listConcat(L *LState) int {
	list := L.CheckOAList(1)
	sep := LString(L.OptString(2, ""))
	i := L.OptInt(3, 1)
	j := L.OptInt(4, list.Len())
	if L.GetTop() == 3 {
		if i > list.Len() || i < 1 {
			L.Push(LString(""))
			return 1
		}
	}
	i = intMax(intMin(i, list.Len()), 1)
	j = intMin(intMin(j, list.Len()), list.Len())
	if i > j {
		L.Push(LString(""))
		return 1
	}
	// TODO should do flushing?
	retbottom := L.GetTop()
	for ; i <= j; i++ {
		L.Push(list.RawGetInt(i))
		if i != j {
			L.Push(sep)
		}
	}
	L.Push(stringConcat(L, L.GetTop()-retbottom, L.reg.Top()-1))
	return 1
}

func listInsert(L *LState) int {
	list := L.CheckOAList(1)
	nargs := L.GetTop()
	if nargs == 1 {
		L.RaiseError("wrong number of arguments")
	}

	if L.GetTop() == 2 {
		list.Append(L.Get(2))
		return 0
	}
	list.Insert(int(L.CheckInt(2)), L.CheckAny(3))
	return 0
}

var ( // JSON pretty formatting strings
	SPC string = ""
	NL  string = ""
)

// listMarshal - marshal list of elements into JSON string
func listMarshal(L *LState) int {
	list := L.CheckOAList(1)
	sep := LString(L.OptString(2, ""))
	nargs := L.GetTop()
	if nargs > 2 || nargs < 1 {
		L.RaiseError("wrong number of arguments")
	}
	lev := 0
	if sep != "" {
		SPC = string(sep)
		NL = "\n"
		lev = 1
	}
	jsonout := ""
	spcs := strings.Repeat(SPC, lev)
	if list.array != nil {
		jsonout = jsonout + spcs + "[" + NL
	} else {
		jsonout = jsonout + spcs + "{" + NL
	}
	jsonout = list.mjson(jsonout, lev+1)
	if list.array != nil {
		jsonout = jsonout + spcs + "]" + NL
	} else {
		jsonout = jsonout + spcs + "}" + NL
	}
	L.Push(LString(jsonout))
	return 1
}

// mjson - recursive function marshalling elements into JSON string
func (lst *LOAList) mjson(jsonout string, lev int) string {
	var spcs string = ""
	if lev > 0 {
		spcs = strings.Repeat(SPC, lev)
	}
	if lst.array != nil {
		i := 0
		la := len(lst.array)
		for j, v := range lst.array {
			if v != LNil {
				x := j + 1
				i++
				switch v.(type) {
				case LBool:
					jsonout = jsonout + spcs + v.String()
				case LNumber:
					jsonout = jsonout + spcs + v.String()
				case LString:
					out, _ := json.Marshal(v.String())
					jsonout = jsonout + spcs + string(out)
				default:
					if lv, ok := v.(*LOAList); ok {
						jsonout = jsonout + spcs + "[" + NL
						jsonout = lv.mjson(jsonout, lev+1)
						jsonout = jsonout + spcs + "]"
					} else {
						fmt.Printf("=== %s(%s)%d: %v\n", spcs, v.Type().String(), x, v)
					}
				}
				if i < la {
					jsonout = jsonout + ","
				}
				jsonout = jsonout + NL
			}
		}
	}
	if lst.strdict != nil {
		i := 0
		ls := len(lst.strdict)
		for k, v := range lst.strdict {
			if v != LNil {
				i++
				switch v.(type) {
				case LBool:
					jsonout = jsonout + spcs + "\"" + k + "\":" + v.String()
				case LNumber:
					jsonout = jsonout + spcs + "\"" + k + "\":" + v.String()
				case LString:
					out, _ := json.Marshal(v.String())
					jsonout = jsonout + spcs + "\"" + k + "\":" + string(out)
				default:
					if lv, ok := v.(*LOAList); ok {
						jsonout = jsonout + spcs + "\"" + k + "\":{" + NL
						jsonout = lv.mjson(jsonout, lev+1)
						jsonout = jsonout + spcs + "}"
					} else {
						fmt.Printf("=== %s(%s)%v: %v\n", spcs, v.Type().String(), k, v)
					}
				}
				if i < ls {
					jsonout = jsonout + ","
				}
				jsonout = jsonout + NL
			}
		}
	}
	if lst.dict != nil {
		i := 0
		ld := len(lst.dict)
		for k, v := range lst.dict {
			if v != LNil {
				i++
				switch v.(type) {
				case LBool:
					jsonout = jsonout + spcs + "\"" + k.String() + "\":" + v.String()
				case LNumber:
					jsonout = jsonout + spcs + "\"" + k.String() + "\":" + v.String()
				case LString:
					out, _ := json.Marshal(v.String())
					jsonout = jsonout + spcs + "\"" + k.String() + "\":" + string(out)
				default:
					if lv, ok := v.(*LOAList); ok {
						if k.Type() == LTString {
							jsonout = jsonout + spcs + "\"" + k.String() + "\":{" + NL
						} else {
							jsonout = jsonout + spcs + "{" + NL
						}
						jsonout = lv.mjson(jsonout, lev+1)
						jsonout = jsonout + spcs + "}"
					} else {
						fmt.Printf("=== %s(%s)%v: %v\n", spcs, v.Type().String(), k, v)
					}
				}
				if i < ld {
					jsonout = jsonout + ","
				}
				jsonout = jsonout + NL
			}
		}
	}
	return jsonout
}

// listMarshal - marshal list of elements into JSON string
func listUnMarshal(L *LState) int {
	jsonBlob := L.CheckString(1)
	nargs := L.GetTop()
	if nargs != 1 {
		L.RaiseError("wrong number of arguments")
	}
	var object interface{}
	var lst *LOAList
	var err error
	err = json.Unmarshal([]byte(jsonBlob), &object)
	if err != nil {
		L.RaiseError(err.Error())
	} else {
		lst = L.NewOAList()
		switch v := object.(type) {
		case int:
			lst.Append(LNumber(v))
		case float64:
			lst.Append(LNumber(v))
		case bool:
			lst.Append(LBool(v))
		case nil:
			lst.Append(LNil)
		case string:
			lst.Append(LString(v))
		case []interface{}:
			ls := newLOAList(0, 0)
			lst.Append(LValue(ls))
			ms := object.([]interface{})
			err = ls.prs(ms)
		case map[string]interface{}:
			ls := newLOAList(0, 0)
			lst.Append(LValue(ls))
			ms := object.(map[string]interface{})
			err = ls.prm(ms)
		default:
			err = errors.New("Invalid data type: " + fmt.Sprintf("%T for %v", v, v))
		}
		//mapStrIfce := object.(map[string]interface{})
		//lst.prm(mapStrIfce)
	}
	if err != nil {
		L.RaiseError(err.Error())
	}

	L.Push(lst)
	return 1
}

func (lst *LOAList) prm(obj map[string]interface{}) error {
	var err error
	for k, d := range obj {
		switch v := d.(type) {
		case int:
			lst.RawSetString(k, LNumber(v))
		case float64:
			lst.RawSetString(k, LNumber(v))
		case bool:
			lst.RawSetString(k, LBool(v))
		case nil:
			lst.RawSetString(k, LNil)
		case string:
			lst.RawSetString(k, LString(v))
		case []interface{}:
			ls := newLOAList(0, 0)
			lst.RawSetString(k, LValue(ls))
			ms := d.([]interface{})
			err = ls.prs(ms)
		case map[string]interface{}:
			ls := newLOAList(0, 0)
			lst.RawSetString(k, LValue(ls))
			ms := d.(map[string]interface{})
			err = ls.prm(ms)
		default:
			err = errors.New("Invalid list data type for key: " + k)
		}
	}
	return err
}

func (lst *LOAList) prs(obj []interface{}) error {
	var err error
	for i, d := range obj {
		switch v := d.(type) {
		case int:
			lst.Append(LNumber(v))
		case float64:
			lst.Append(LNumber(v))
		case bool:
			lst.Append(LBool(v))
		case nil:
			lst.Append(LNil)
		case string:
			lst.Append(LString(v))
		case []interface{}:
			ls := newLOAList(0, 0)
			lst.Append(LValue(ls))
			ms := d.([]interface{})
			err = ls.prs(ms)
		case map[string]interface{}:
			ls := newLOAList(0, 0)
			lst.Append(LValue(ls))
			ms := d.(map[string]interface{})
			err = ls.prm(ms)
		default:
			err = errors.New("Invalid list data type for index: " + strconv.Itoa(i))
		}
	}
	return err
}

// listDump - dump the contents of a list
func listDump(L *LState) int {
	list := L.CheckOAList(1)
	nargs := L.GetTop()
	if nargs != 1 {
		L.RaiseError("wrong number of arguments")
	}
	lev := 1
	spcs := strings.Repeat(SPC, lev)
	if list.array != nil {
		fmt.Printf("%s(array): [\n", spcs)
	} else if list.strdict != nil {
		fmt.Printf("%s(strdict): {\n", spcs)
	} else {
		fmt.Printf("%s(dict): {\n", spcs)
	}
	list.dumplist(lev + 1)
	if list.array != nil {
		fmt.Printf("%s]\n", spcs)
	} else {
		fmt.Printf("%s}\n", spcs)
	}
	return 0
}

// dumplist - iterates over list of elements, printing each
func (lst *LOAList) dumplist(lev int) {
	spcs := strings.Repeat(SPC, lev)
	if lst.array != nil {
		for i, v := range lst.array {
			if v != LNil {
				x := i + 1
				switch v.(type) {
				case LNumber:
					fmt.Printf("%s(num)%d: %v\n", spcs, x, v)
				case LString:
					fmt.Printf("%s(str)%d: \"%v\"\n", spcs, x, v)
				default:
					if lv, ok := v.(*LOAList); ok {
						fmt.Printf("%s(array): [\n", spcs)
						lv.dumplist(lev + 1)
						fmt.Printf("%s]\n", spcs)
					} else {
						fmt.Printf("%s(%s)%d: %v\n", spcs, v.Type().String(), x, v)
					}
				}
			}
		}
	}
	if lst.strdict != nil {
		for k, v := range lst.strdict {
			if v != LNil {
				switch v.(type) {
				case LNumber:
					fmt.Printf("%s(num)%v: %v\n", spcs, k, v)
				case LString:
					fmt.Printf("%s(str)%v: \"%v\"\n", spcs, k, v)
				default:
					if lv, ok := v.(*LOAList); ok {
						fmt.Printf("%s(strdict)%v: {\n", spcs, k)
						lv.dumplist(lev + 1)
						fmt.Printf("%s}\n", spcs)
					} else {
						fmt.Printf("%s(%s)%v: %v\n", spcs, v.Type().String(), k, v)
					}
				}
			}
		}
	}
	if lst.dict != nil {
		for k, v := range lst.dict {
			if v != LNil {
				switch v.(type) {
				case LNumber:
					fmt.Printf("%s(num)%v: %v\n", spcs, k, v)
				case LString:
					fmt.Printf("%s(str)%v: \"%v\"\n", spcs, k, v)
				default:
					if lv, ok := v.(*LOAList); ok {
						if k.Type() == LTString {
							fmt.Printf("%s(dict)%v: {\n", spcs, k)
						} else {
							fmt.Printf("%s(dict): {\n", spcs)
						}
						lv.dumplist(lev + 1)
						fmt.Printf("%s}\n", spcs)
					} else {
						fmt.Printf("%s(%s)%v: %v\n", spcs, v.Type().String(), k, v)
					}
				}
			}
		}
	}
}

// listMarshalXml - marshal list of elements into XML string
func listMarshalXml(L *LState) int {
	list := L.CheckOAList(1)
	name := L.OptString(2, "list")
	rep := L.OptString(3, "")
	// verify name valid name and not blank
	if !IsXmlTagName(name) {
		if rep == "" {
			name, _ = MakeXmlTagName(name, "c", "", "list")
		} else {
			name, _ = MakeXmlTagName(name, "r", rep, "list")
		}
	}
	sep := LString(L.OptString(3, ""))
	nargs := L.GetTop()
	if nargs > 3 || nargs < 1 {
		L.RaiseError("wrong number of arguments")
	}
	lev := 0
	if sep != "" {
		SPC = string(sep)
		NL = "\n"
		lev = 1
	}
	xmlout := ""
	spcs := strings.Repeat(SPC, lev)
	xmlout = xmlout + spcs + "<" + name + ">" + NL
	xmlout = list.mxml(xmlout, name, lev+1)
	xmlout = xmlout + spcs + "</" + name + ">" + NL
	L.Push(LString(xmlout))
	return 1
}

// mxml - recursive function marshalling elements into XML string
func (lst *LOAList) mxml(xmlout string, name string, lev int) string {
	var spcs string = ""
	if lev > 0 {
		spcs = strings.Repeat(SPC, lev)
	}
	if lst.array != nil {
		for j, v := range lst.array {
			if v != LNil {
				x := j + 1
				switch v.(type) {
				case LBool:
					xmlout = xmlout + spcs + "<" + name + ">" + v.String() + "</" + name + ">"
				case LNumber:
					xmlout = xmlout + spcs + "<" + name + ">" + v.String() + "</" + name + ">"
				case LString:
					out, _ := xml.Marshal(v.String())
					str := string(out)
					str = strings.TrimPrefix(str, "<string>")
					str = strings.TrimSuffix(str, "</string>")
					xmlout = xmlout + spcs + "<" + name + ">" + str + "</" + name + ">"
				default:
					if lv, ok := v.(*LOAList); ok {
						xmlout = xmlout + spcs + "<" + name + ">" + NL
						xmlout = lv.mxml(xmlout, name, lev+1)
						xmlout = xmlout + spcs + "</" + name + ">"
					} else {
						fmt.Printf("<!-- %s(%s)%d: %v -->\n", spcs, v.Type().String(), x, v)
					}
				}
				xmlout = xmlout + NL
			}
		}
	}
	if lst.strdict != nil {
		for k, v := range lst.strdict {
			if v != LNil {
				switch v.(type) {
				case LBool:
					xmlout = xmlout + spcs + "<" + k + ">" + v.String() + "</" + k + ">"
				case LNumber:
					xmlout = xmlout + spcs + "<" + k + ">" + v.String() + "</" + k + ">"
				case LString:
					out, _ := xml.Marshal(v.String())
					str := string(out)
					str = strings.TrimPrefix(str, "<string>")
					str = strings.TrimSuffix(str, "</string>")
					xmlout = xmlout + spcs + "<" + k + ">" + str + "</" + k + ">"
				default:
					if lv, ok := v.(*LOAList); ok {
						xmlout = xmlout + spcs + "<" + k + ">" + NL
						xmlout = lv.mxml(xmlout, k, lev+1)
						xmlout = xmlout + spcs + "</" + k + ">"
					} else {
						fmt.Printf("<!-- %s(%s)%v: %v -->\n", spcs, v.Type().String(), k, v)
					}
				}
				xmlout = xmlout + NL
			}
		}
	}
	if lst.dict != nil {
		for k, v := range lst.dict {
			if v != LNil {
				switch v.(type) {
				case LBool:
					xmlout = xmlout + spcs + "<" + k.String() + ">" + v.String() + "</" + k.String() + ">"
				case LNumber:
					xmlout = xmlout + spcs + "<" + k.String() + ">" + v.String() + "</" + k.String() + ">"
				case LString:
					out, _ := xml.Marshal(v.String())
					str := string(out)
					str = strings.TrimPrefix(str, "<string>")
					str = strings.TrimSuffix(str, "</string>")
					xmlout = xmlout + spcs + "<" + k.String() + ">" + str + "</" + k.String() + ">"
				default:
					if lv, ok := v.(*LOAList); ok {
						if k.Type() == LTString {
							xmlout = xmlout + spcs + "<" + k.String() + ">" + NL
							xmlout = lv.mxml(xmlout, k.String(), lev+1)
							xmlout = xmlout + spcs + "</" + k.String() + ">" + NL
						}
					} else {
						fmt.Printf("<!-- %s(%s)%v: %v -->\n", spcs, v.Type().String(), k, v)
					}
				}
				xmlout = xmlout + NL
			}
		}
	}
	return xmlout
}
