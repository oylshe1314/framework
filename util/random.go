package util

import (
	"crypto/md5"
	"encoding/hex"
	"math/rand"
)

func NewRandom() *rand.Rand {
	return rand.New(rand.NewSource(UnixMilli()))
}

var defaultRandom = NewRandom()

// RandomRange 在给定范围内随机一个数
func RandomRange[NT Number](min, max NT) NT {
	return min + NT(defaultRandom.Float64()*float64(max-min))
}

// RandomSelect 从给定列表随机选择N个
func RandomSelect[T any](list []T, num int) []T {
	if len(list) == 0 {
		return nil
	}

	if len(list) <= num {
		return list
	}

	var is = make([]int, len(list))
	for i := range list {
		is[i] = i
	}

	for i := 0; i < num; i++ {
		var ri = defaultRandom.Intn(len(is))
		is[i], is[ri] = is[ri], is[i]
	}

	var nl = make([]T, num)
	for i := 0; i < num; i++ {
		nl[i] = list[is[i]]
	}

	return nl
}

// RandomWeights 按权重随机
func RandomWeights[T any, W Number](list []T, weightFunc func(i int) W) (t T) {
	if len(list) == 0 {
		return
	}

	var rates float64
	for i := range list {
		rates += float64(weightFunc(i))
	}

	var factor = defaultRandom.Float64() * rates
	for i, e := range list {
		rates -= float64(weightFunc(i))
		if factor >= rates {
			return e
		}
	}
	return list[defaultRandom.Intn(len(list))]
}

const (
	// CharsUpperLetter 大写字母
	CharsUpperLetter = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"

	// CharsLowerletter 大写字母
	CharsLowerletter = "abcdefghijklmnopqrstuvwxyz"

	//CharsAllLetter 所有所有字母
	CharsAllLetter = CharsUpperLetter + CharsLowerletter

	// CharsAllNumber 所有数字
	CharsAllNumber = "0123456789"

	// CharsNumbersAndLetter 数字加字母
	CharsNumbersAndLetter = CharsAllNumber + CharsAllLetter
)

func RandomStrings(chars string, num int, repeated bool) string {
	var cl = len(chars)
	if cl == 0 {
		chars = CharsNumbersAndLetter
		cl = len(chars)
	}

	if repeated {
		var ret = make([]byte, num)
		for i := 0; i < num; i++ {
			ret[i] = chars[defaultRandom.Intn(cl)]
		}
		return string(ret)
	} else {
		if len(chars) < num {
			return chars
		}

		return string(RandomSelect([]byte(chars), num))
	}
}

func RandomToken() string {
	var src []byte
	src = Uint64ToBytes(src, uint64(UnixMilli())<<48|(uint64(defaultRandom.Int63n(65536))&0xFFFF))
	src = Uint64ToBytes(src, defaultRandom.Uint64())
	var h = md5.New()
	h.Write(src)
	return hex.EncodeToString(h.Sum(nil))
}
