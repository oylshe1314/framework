package util

import (
	"math/rand/v2"
)

func NewRand() *rand.Rand {
	return rand.New(rand.NewPCG(uint64(UnixMicro()), 0xC5))
}

func NewRandBySeed(seed int64) *rand.Rand {
	return rand.New(rand.NewPCG(uint64(seed), 0xC5))
}

// RandomRange 在给定范围内随机一个数
func RandomRange[NT Number](min, max NT) NT {
	return min + NT(NewRand().Float64()*float64(max-min))
}

// RandomSelect 从给定列表随机选择num个
func RandomSelect[T any](list []T, num int) []T {
	if len(list) == 0 {
		return nil
	}

	if len(list) <= num {
		return list
	}

	var r = NewRand()
	var is = make([]int, len(list))
	for i := range list {
		is[i] = i
	}

	for i := 0; i < num; i++ {
		var ri = r.IntN(len(is))
		is[i], is[ri] = is[ri], is[i]
	}

	var nl = make([]T, num)
	for i := 0; i < num; i++ {
		nl[i] = list[is[i]]
	}

	return nl
}

// RandomWeights 按权重随机列表的一个
func RandomWeights[T any, W Number](list []T, weightFunc func(i int) W) (t T) {
	if len(list) == 0 {
		return
	}

	var rates float64
	var r = NewRand()
	for i := range list {
		rates += float64(weightFunc(i))
	}

	var factor = r.Float64() * rates
	for i, e := range list {
		rates -= float64(weightFunc(i))
		if factor >= rates {
			return e
		}
	}
	return list[r.IntN(len(list))]
}

const (
	// CharsUpperLetter 大写字母
	CharsUpperLetter = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"

	// CharsLowerLetter 小写字母
	CharsLowerLetter = "abcdefghijklmnopqrstuvwxyz"

	// CharsAllNumber 所有数字
	CharsNumber = "0123456789"

	//CharsAllLetter 所有所有字母
	CharsAllLetter = CharsUpperLetter + CharsLowerLetter

	// CharsNumbersAndLetter 数字加大写字母
	CharsNumbersAndUpper = CharsNumber + CharsUpperLetter

	// CharsNumbersAndLetter 数字加小写字母
	CharsNumbersAndLower = CharsNumber + CharsLowerLetter

	// CharsNumbersAndLetter 数字加小写字母
	CharsNumbersAndAllLetter = CharsNumber + CharsUpperLetter + CharsLowerLetter
)

func RandomStrings(chars string, num int, repeated bool) string {
	var cl = len(chars)
	if cl == 0 {
		chars = CharsNumbersAndAllLetter
		cl = len(chars)
	}

	var r = NewRand()
	if repeated {
		var ret = make([]byte, num)
		for i := 0; i < num; i++ {
			ret[i] = chars[r.IntN(cl)]
		}
		return string(ret)
	} else {
		if len(chars) < num {
			return chars
		}

		return string(RandomSelect([]byte(chars), num))
	}
}
