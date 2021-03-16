// Package qs - q scripting language
package qs

const defaultArrayCap = 32
const defaultHashCap = 32

type lValueArraySorter struct {
	L      *LState
	Fn     *LProc
	Values []LValue
}

// Len - list value sort method
func (lv lValueArraySorter) Len() int {
	return len(lv.Values)
}

// Swap - list value sort method
func (lv lValueArraySorter) Swap(i, j int) {
	lv.Values[i], lv.Values[j] = lv.Values[j], lv.Values[i]
}

// Less - list value sort method
func (lv lValueArraySorter) Less(i, j int) bool {
	if lv.Fn != nil {
		lv.L.Push(lv.Fn)
		lv.L.Push(lv.Values[i])
		lv.L.Push(lv.Values[j])
		lv.L.Call(2, 1)
		return LVAsBool(lv.L.reg.Pop())
	}
	return lessThan(lv.L, lv.Values[i], lv.Values[j])
}

// newLOAList - creates a new OAList
func newLOAList(acap int, hcap int) *LOAList {
	if acap < 0 {
		acap = 0
	}
	if hcap < 0 {
		hcap = 0
	}
	lst := &LOAList{}
	lst.keys = nil
	lst.k2i = nil
	lst.Metalist = LNil
	if acap != 0 {
		lst.array = make([]LValue, 0, acap)
	}
	if hcap != 0 {
		lst.strdict = make(map[string]LValue, hcap)
	}
	return lst
}

// Len - returns length of this LOAList.
func (lst *LOAList) Len() int {
	if lst.array == nil {
		return 0
	}
	var prev LValue = LNil
	for i := len(lst.array) - 1; i >= 0; i-- {
		v := lst.array[i]
		if prev == LNil && v != LNil {
			return i + 1
		}
		prev = v
	}
	return 0
}

// Append - appends a given LValue to this LOAList.
func (lst *LOAList) Append(value LValue) {
	if lst.array == nil {
		lst.array = make([]LValue, 0, defaultArrayCap)
	}
	lst.array = append(lst.array, value)
}

// Insert - inserts a given LValue at position `i` in this list.
func (lst *LOAList) Insert(i int, value LValue) {
	if lst.array == nil {
		lst.array = make([]LValue, 0, defaultArrayCap)
	}
	if i > len(lst.array) {
		lst.RawSetInt(i, value)
		return
	}
	if i <= 0 {
		lst.RawSet(LNumber(i), value)
		return
	}
	i -= 1
	lst.array = append(lst.array, LNil)
	copy(lst.array[i+1:], lst.array[i:])
	lst.array[i] = value
}

// MaxN - returns a maximum number key that nil value does not exist before it.
func (lst *LOAList) MaxN() int {
	if lst.array == nil {
		return 0
	}
	for i := len(lst.array) - 1; i >= 0; i-- {
		if lst.array[i] != LNil {
			return i + 1
		}
	}
	return 0
}

// Remove - removes from this list the element at a given position.
func (lst *LOAList) Remove(pos int) LValue {
	if lst.array == nil {
		return LNil
	}
	i := pos - 1
	larray := len(lst.array)
	oldval := LNil
	switch {
	case i >= larray:
		// nothing to do
	case i == larray-1 || i < 0:
		oldval = lst.array[larray-1]
		lst.array = lst.array[:larray-1]
	default:
		oldval = lst.array[i]
		copy(lst.array[i:], lst.array[i+1:])
		lst.array[larray-1] = nil
		lst.array = lst.array[:larray-1]
	}
	return oldval
}

// RawSet - sets a given LValue to a given index without the __newindex metamethod.
// It is recommended to use `RawSetString` or `RawSetInt` for performance
// if you already know the given LValue is a string or number.
func (lst *LOAList) RawSet(key LValue, value LValue) {
	switch v := key.(type) {
	case LNumber:
		if isArrayKey(v) {
			if lst.array == nil {
				lst.array = make([]LValue, 0, defaultArrayCap)
			}
			index := int(v) - 1
			alen := len(lst.array)
			switch {
			case index == alen:
				lst.array = append(lst.array, value)
			case index > alen:
				for i := 0; i < (index - alen); i++ {
					lst.array = append(lst.array, LNil)
				}
				lst.array = append(lst.array, value)
			case index < alen:
				lst.array[index] = value
			}
			return
		}
	case LString:
		lst.RawSetString(string(v), value)
		return
	}

	lst.RawSetH(key, value)
}

// RawSetInt - sets a given LValue at a position `key` without the __newindex metamethod.
func (lst *LOAList) RawSetInt(key int, value LValue) {
	if key < 1 || key >= MaxArrayIndex {
		lst.RawSetH(LNumber(key), value)
		return
	}
	if lst.array == nil {
		lst.array = make([]LValue, 0, 32)
	}
	index := key - 1
	alen := len(lst.array)
	switch {
	case index == alen:
		lst.array = append(lst.array, value)
	case index > alen:
		for i := 0; i < (index - alen); i++ {
			lst.array = append(lst.array, LNil)
		}
		lst.array = append(lst.array, value)
	case index < alen:
		lst.array[index] = value
	}
}

// RawSetString - sets a given LValue to a given string index without the __newindex metamethod.
func (lst *LOAList) RawSetString(key string, value LValue) {
	if lst.strdict == nil {
		lst.strdict = make(map[string]LValue, defaultHashCap)
	}
	if value == LNil {
		delete(lst.strdict, key)
	} else {
		lst.strdict[key] = value
	}
}

// RawSetH - sets a given LValue to a given index without the __newindex metamethod.
func (lst *LOAList) RawSetH(key LValue, value LValue) {
	if s, ok := key.(LString); ok {
		lst.RawSetString(string(s), value)
		return
	}
	if lst.dict == nil {
		lst.dict = make(map[LValue]LValue, len(lst.strdict))
	}

	if value == LNil {
		delete(lst.dict, key)
	} else {
		lst.dict[key] = value
	}
}

// RawGet - returns an LValue associated with a given key without __index metamethod.
func (lst *LOAList) RawGet(key LValue) LValue {
	switch v := key.(type) {
	case LNumber:
		if isArrayKey(v) {
			if lst.array == nil {
				return LNil
			}
			index := int(v) - 1
			if index >= len(lst.array) {
				return LNil
			}
			return lst.array[index]
		}
	case LString:
		if lst.strdict == nil {
			return LNil
		}
		if ret, ok := lst.strdict[string(v)]; ok {
			return ret
		}
		return LNil
	}
	if lst.dict == nil {
		return LNil
	}
	if v, ok := lst.dict[key]; ok {
		return v
	}
	return LNil
}

// RawGetInt - returns an LValue at position `key` without __index metamethod.
func (lst *LOAList) RawGetInt(key int) LValue {
	if lst.array == nil {
		return LNil
	}
	index := int(key) - 1
	if index >= len(lst.array) || index < 0 {
		return LNil
	}
	return lst.array[index]
}

// RawGet - returns an LValue associated with a given key without __index metamethod.
func (lst *LOAList) RawGetH(key LValue) LValue {
	if s, sok := key.(LString); sok {
		if lst.strdict == nil {
			return LNil
		}
		if v, vok := lst.strdict[string(s)]; vok {
			return v
		}
		return LNil
	}
	if lst.dict == nil {
		return LNil
	}
	if v, ok := lst.dict[key]; ok {
		return v
	}
	return LNil
}

// RawGetString - returns an LValue associated with a given key without __index metamethod.
func (lst *LOAList) RawGetString(key string) LValue {
	if lst.strdict == nil {
		return LNil
	}
	if v, vok := lst.strdict[string(key)]; vok {
		return v
	}
	return LNil
}

// ForEach - iterates over this list of elements, yielding each in turn to a given proc.
func (lst *LOAList) ForEach(cb func(LValue, LValue)) {
	if lst.array != nil {
		for i, v := range lst.array {
			if v != LNil {
				cb(LNumber(i+1), v)
			}
		}
	}
	if lst.strdict != nil {
		for k, v := range lst.strdict {
			if v != LNil {
				cb(LString(k), v)
			}
		}
	}
	if lst.dict != nil {
		for k, v := range lst.dict {
			if v != LNil {
				cb(k, v)
			}
		}
	}
}

// Next -
func (lst *LOAList) Next(key LValue) (LValue, LValue) {
	// TODO: inefficient way
	init := false
	if key == LNil {
		lst.keys = nil
		lst.k2i = nil
		key = LNumber(0)
		init = true
	}

	length := 0
	if lst.dict != nil {
		length += len(lst.dict)
	}
	if lst.strdict != nil {
		length += len(lst.strdict)
	}

	if lst.keys == nil {
		lst.keys = make([]LValue, length)
		lst.k2i = make(map[LValue]int)
		i := 0
		if lst.dict != nil {
			for k, _ := range lst.dict {
				lst.keys[i] = k
				lst.k2i[k] = i
				i++
			}
		}
		if lst.strdict != nil {
			for k, _ := range lst.strdict {
				lst.keys[i] = LString(k)
				lst.k2i[LString(k)] = i
				i++
			}
		}
	}

	if init || key != LNumber(0) {
		if kv, ok := key.(LNumber); ok && isInteger(kv) && int(kv) >= 0 {
			index := int(kv)
			if lst.array != nil {
				for ; index < len(lst.array); index++ {
					if v := lst.array[index]; v != LNil {
						return LNumber(index + 1), v
					}
				}
			}
			if lst.array == nil || index == len(lst.array) {
				if (lst.dict == nil || len(lst.dict) == 0) && (lst.strdict == nil || len(lst.strdict) == 0) {
					lst.keys = nil
					lst.k2i = nil
					return LNil, LNil
				}
				key = lst.keys[0]
				if v := lst.RawGetH(key); v != LNil {
					return key, v
				}
			}
		}
	}

	for i := lst.k2i[key] + 1; i < length; i++ {
		key = lst.keys[lst.k2i[key]+1]
		if v := lst.RawGetH(key); v != LNil {
			return key, v
		}
	}
	lst.keys = nil
	lst.k2i = nil
	return LNil, LNil
}

// HasArray - returns true if this list has array elements.
func (lst *LOAList) HasArray() bool {
	if lst.array == nil {
		return false
	}
	return true
}

// HasStrDict - returns true if this list has strdict elements.
func (lst *LOAList) HasStrDict() bool {
	if lst.strdict == nil {
		return false
	}
	return true
}

// HasDict - returns true if this list has dict elements.
func (lst *LOAList) HasDict() bool {
	if lst.dict == nil {
		return false
	}
	return true
}

// ForEachArray - iterates over this lists array elements, yielding each in turn to a given proc.
func (lst *LOAList) ForEachArray(cb func(LValue, LValue)) {
	if lst.array != nil {
		for i, v := range lst.array {
			if v != LNil {
				cb(LNumber(i+1), v)
			}
		}
	}
}

// ForEachStrDict - iterates over this lists strdict elements, yielding each in turn to a given proc.
func (lst *LOAList) ForEachStrDict(cb func(LValue, LValue)) {
	if lst.strdict != nil {
		for k, v := range lst.strdict {
			if v != LNil {
				cb(LString(k), v)
			}
		}
	}
}

// ForEachDict - iterates over this lists dict elements, yielding each in turn to a given proc.
func (lst *LOAList) ForEachDict(cb func(LValue, LValue)) {
	if lst.dict != nil {
		for k, v := range lst.dict {
			if v != LNil {
				cb(k, v)
			}
		}
	}
}
