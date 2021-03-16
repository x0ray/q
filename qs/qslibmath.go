// Package qs - q scripting language
package qs

import (
	"math"
	"math/rand"
	"sort"
)

func mathAbs(L *LState) int {
	L.Push(LNumber(math.Abs(float64(L.CheckNumber(1)))))
	return 1
}

func mathAcos(L *LState) int {
	L.Push(LNumber(math.Acos(float64(L.CheckNumber(1)))))
	return 1
}

func mathAsin(L *LState) int {
	L.Push(LNumber(math.Asin(float64(L.CheckNumber(1)))))
	return 1
}

func mathAtan(L *LState) int {
	L.Push(LNumber(math.Atan(float64(L.CheckNumber(1)))))
	return 1
}

func mathAtan2(L *LState) int {
	L.Push(LNumber(math.Atan2(float64(L.CheckNumber(1)), float64(L.CheckNumber(2)))))
	return 1
}

func mathCeil(L *LState) int {
	L.Push(LNumber(math.Ceil(float64(L.CheckNumber(1)))))
	return 1
}

func mathCos(L *LState) int {
	L.Push(LNumber(math.Cos(float64(L.CheckNumber(1)))))
	return 1
}

func mathCosh(L *LState) int {
	L.Push(LNumber(math.Cosh(float64(L.CheckNumber(1)))))
	return 1
}

func mathDeg(L *LState) int {
	L.Push(LNumber(float64(L.CheckNumber(1)) * 180 / math.Pi))
	return 1
}

func mathExp(L *LState) int {
	L.Push(LNumber(math.Exp(float64(L.CheckNumber(1)))))
	return 1
}

func mathFact(L *LState) int {
	f := int(L.CheckNumber(1))
	var fact int = 1
	for i := f; i > 0; i-- {
		fact = fact * i
	}
	L.Push(LNumber(fact))
	return 1
}

func fib(n int) int {
	if n < 2 {
		return n
	}
	return fib(n-2) + fib(n-1)
}

func mathFib(L *LState) int {
	f := int(L.CheckNumber(1))
	var r int
	r = fib(f)
	L.Push(LNumber(r))
	return 1
}

func mathFloor(L *LState) int {
	L.Push(LNumber(math.Floor(float64(L.CheckNumber(1)))))
	return 1
}

func mathFmod(L *LState) int {
	L.Push(LNumber(math.Mod(float64(L.CheckNumber(1)), float64(L.CheckNumber(2)))))
	return 1
}

func mathFrexp(L *LState) int {
	v1, v2 := math.Frexp(float64(L.CheckNumber(1)))
	L.Push(LNumber(v1))
	L.Push(LNumber(v2))
	return 2
}

func mathLdexp(L *LState) int {
	L.Push(LNumber(math.Ldexp(float64(L.CheckNumber(1)), L.CheckInt(2))))
	return 1
}

func mathLog(L *LState) int {
	L.Push(LNumber(math.Log(float64(L.CheckNumber(1)))))
	return 1
}

func mathLog10(L *LState) int {
	L.Push(LNumber(math.Log10(float64(L.CheckNumber(1)))))
	return 1
}

func mathMax(L *LState) int {
	if L.GetTop() == 0 {
		L.RaiseError("max no arguments")
	}
	max := L.CheckNumber(1)
	top := L.GetTop()
	for i := 2; i <= top; i++ {
		v := L.CheckNumber(i)
		if v > max {
			max = v
		}
	}
	L.Push(LNumber(max))
	return 1
}

func mathMean(L *LState) int {
	if L.GetTop() == 0 {
		L.RaiseError("mean no arguments")
	}
	sum := float64(L.CheckNumber(1))
	count := 1.0
	top := L.GetTop()
	for i := 2; i <= top; i++ {
		sum = sum + float64(L.CheckNumber(i))
		count++
	}
	mean := sum / count
	L.Push(LNumber(mean))
	return 1
}

type ByNumber []float64

func (a ByNumber) Len() int           { return len(a) }
func (a ByNumber) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByNumber) Less(i, j int) bool { return a[i] < a[j] }

func mathMedian(L *LState) int {
	if L.GetTop() == 0 {
		L.RaiseError("median no arguments")
	}
	var median float64
	if L.GetTop() == 1 {
		median = float64(L.CheckNumber(1))
	} else {
		// collect argument values
		var vals []float64
		top := L.GetTop() // number of values
		for i := 1; i <= top; i++ {
			vals = append(vals, float64(L.CheckNumber(i)))
		}
		// sort
		sort.Sort(ByNumber(vals))
		// get middle
		if (top % 2) == 0 { // even number of values -- get avg of two middle
			ix := len(vals) / 2
			median = (vals[ix] + vals[ix+1]) / 2
		} else { // odd number of values -- get middle value
			ix := (len(vals) / 2) + 1
			median = vals[ix]
		}
	}
	L.Push(LNumber(median))
	return 1
}

func mathMin(L *LState) int {
	if L.GetTop() == 0 {
		L.RaiseError("min no arguments")
	}
	min := L.CheckNumber(1)
	top := L.GetTop()
	for i := 2; i <= top; i++ {
		v := L.CheckNumber(i)
		if v < min {
			min = v
		}
	}
	L.Push(LNumber(min))
	return 1
}

func mathMod(L *LState) int {
	lhs := L.CheckNumber(1)
	rhs := L.CheckNumber(2)
	L.Push(oaModulo(lhs, rhs))
	return 1
}

func mathMode(L *LState) int {
	if L.GetTop() == 0 {
		L.RaiseError("mode no arguments")
	}
	var mode float64
	if L.GetTop() == 1 {
		mode = float64(L.CheckNumber(1))
	} else {
		// collect argument values
		var vals []float64
		top := L.GetTop() // number of values
		for i := 1; i <= top; i++ {
			vals = append(vals, float64(L.CheckNumber(i)))
		}
		// sort
		sort.Sort(sort.Reverse(ByNumber(vals)))
		// get mode -- most common
		var lv, v float64
		var maxcnt, cnt int
		for i := 0; i < len(vals); i++ {
			lv = v
			v = vals[i]
			if v != lv { // transition
				if cnt > maxcnt { // new max cnt
					maxcnt = cnt
					mode = lv
				}
				cnt = 1
			} else {
				cnt++
			}
		}
		if cnt > maxcnt { // new max cnt
			maxcnt = cnt
			mode = lv
		}
	}
	L.Push(LNumber(mode))
	return 1
}

func mathModf(L *LState) int {
	v1, v2 := math.Modf(float64(L.CheckNumber(1)))
	L.Push(LNumber(v1))
	L.Push(LNumber(v2))
	return 2
}

func mathPow(L *LState) int {
	L.Push(LNumber(math.Pow(float64(L.CheckNumber(1)), float64(L.CheckNumber(2)))))
	return 1
}

func mathRad(L *LState) int {
	L.Push(LNumber(float64(L.CheckNumber(1)) * math.Pi / 180))
	return 1
}

func mathRandom(L *LState) int {
	switch L.GetTop() {
	case 0:
		L.Push(LNumber(rand.Float64()))
	case 1:
		n := L.CheckInt(1)
		L.Push(LNumber(rand.Intn(n-1) + 1))
	default:
		min := L.CheckInt(1)
		max := L.CheckInt(2) + 1
		L.Push(LNumber(rand.Intn(max-min) + min))
	}
	return 1
}

func mathRandomseed(L *LState) int {
	rand.Seed(L.CheckInt64(1))
	return 0
}

func mathRange(L *LState) int {
	if L.GetTop() == 0 {
		L.RaiseError("range no arguments")
	}
	min := float64(L.CheckNumber(1))
	max := min
	top := L.GetTop()
	for i := 2; i <= top; i++ {
		v := float64(L.CheckNumber(i))
		if v < min {
			min = v
		} else if v > max {
			max = v
		}
	}
	rnge := max - min
	L.Push(LNumber(rnge))
	return 1
}

func mathRms(L *LState) int {
	if L.GetTop() == 0 {
		L.RaiseError("rms no arguments")
	}
	n := float64(L.CheckNumber(1))
	ssq := n * n
	count := 1.0
	top := L.GetTop()
	for i := 2; i <= top; i++ {
		n = float64(L.CheckNumber(i))
		ssq = ssq + n*n
		count++
	}
	rms := math.Sqrt(ssq / count)
	L.Push(LNumber(rms))
	return 1
}

func mathSin(L *LState) int {
	L.Push(LNumber(math.Sin(float64(L.CheckNumber(1)))))
	return 1
}

func mathSinh(L *LState) int {
	L.Push(LNumber(math.Sinh(float64(L.CheckNumber(1)))))
	return 1
}

func mathSqrt(L *LState) int {
	L.Push(LNumber(math.Sqrt(float64(L.CheckNumber(1)))))
	return 1
}

func mathStdDev(L *LState) int {
	if L.GetTop() == 0 {
		L.RaiseError("stddev no arguments")
	}
	n := float64(L.CheckNumber(1))
	sum := n
	ssq := n * n
	count := 1.0
	top := L.GetTop()
	for i := 2; i <= top; i++ {
		n = float64(L.CheckNumber(i))
		sum = sum + n
		ssq = ssq + n*n
		count++
	}
	sumsq := (sum * sum) / count
	vari := (ssq - sumsq) / (count - 1)
	sd := math.Sqrt(vari)
	L.Push(LNumber(sd))
	return 1
}

func mathSum(L *LState) int {
	if L.GetTop() == 0 {
		L.RaiseError("sum no arguments")
	}
	sum := float64(L.CheckNumber(1))
	top := L.GetTop()
	for i := 2; i <= top; i++ {
		sum = sum + float64(L.CheckNumber(i))
	}
	L.Push(LNumber(sum))
	return 1
}

func mathTan(L *LState) int {
	L.Push(LNumber(math.Tan(float64(L.CheckNumber(1)))))
	return 1
}

func mathTanh(L *LState) int {
	L.Push(LNumber(math.Tanh(float64(L.CheckNumber(1)))))
	return 1
}

func mathVariance(L *LState) int {
	if L.GetTop() == 0 {
		L.RaiseError("variance no arguments")
	}
	n := float64(L.CheckNumber(1))
	sum := n
	ssq := n * n
	count := 1.0
	top := L.GetTop()
	for i := 2; i <= top; i++ {
		n = float64(L.CheckNumber(i))
		sum = sum + n
		ssq = ssq + n*n
		count++
	}
	sumsq := (sum * sum) / count
	vari := (ssq - sumsq) / (count - 1)
	L.Push(LNumber(vari))
	return 1
}
