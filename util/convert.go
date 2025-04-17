package util

import (
	"github.com/oylshe1314/framework/errors"
	"strconv"
	"strings"
)

type Integer interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr
}

type Float interface {
	~float32 | ~float64
}

type Number interface {
	Integer | Float
}

func StringToInteger1[T Integer](s string) (t T, err error) {
	err = StringToInteger2(s, &t)
	return
}

func StringToInteger2[T Integer](s string, t *T) error {
	if len(s) == 0 {
		*t = 0
		return nil
	}

	n, err := strconv.ParseInt(s, 10, 64)
	*t = T(n)
	return err
}

func StringToFloat1[T Float](s string) (t T, err error) {
	err = StringToFloat2(s, &t)
	return
}

func StringToFloat2[T Float](s string, t *T) error {
	if len(s) == 0 {
		*t = 0
		return nil
	}

	n, err := strconv.ParseFloat(s, 64)
	*t = T(n)
	return err
}

func StringsToIntegers1[T Integer](ss []string) (ts []T, err error) {
	err = StringsToIntegers2(ss, &ts)
	return
}

func StringsToIntegers2[T Integer](ss []string, ts *[]T) (err error) {
	*ts = make([]T, len(ss))
	for i, s := range ss {
		err = StringToInteger2(s, &(*ts)[i])
		if err != nil {
			return err
		}
	}
	return nil
}

func StringsToFloats1[T Float](ss []string) (ts []T, err error) {
	err = StringsToFloats2(ss, &ts)
	return
}

func StringsToFloats2[T Float](ss []string, ts *[]T) (err error) {
	*ts = make([]T, len(ss))
	for i, s := range ss {
		err = StringToFloat2(s, &(*ts)[i])
		if err != nil {
			return err
		}
	}
	return nil
}

func IntegerToString[T Integer](t T) string {
	return strconv.FormatInt(int64(t), 10)
}

func FloatToString[T Float](t T, prec int) string {
	return strconv.FormatFloat(float64(t), 'f', prec, 64)
}

func IntegersToStrings[T Integer](ts []T) []string {
	var ss = make([]string, len(ts))
	for i, t := range ts {
		ss[i] = IntegerToString(t)
	}
	return ss
}

func FloatsToStrings[T Float](ts []T, prec int) []string {
	var ss = make([]string, len(ts))
	for i, t := range ts {
		ss[i] = FloatToString(t, prec)
	}
	return ss
}

func SplitToIntegers1[T Integer](str, sep string) ([]T, error) {
	return StringsToIntegers1[T](strings.Split(str, sep))
}

func SplitToIntegers2[T Integer](str, sep string, ts *[]T) error {
	return StringsToIntegers2(strings.Split(str, sep), ts)
}

func SplitToFloats1[T Float](str, sep string) ([]T, error) {
	if sep == "." {
		return nil, errors.Error("can not use '.' to split the floats in the string")
	}
	return StringsToFloats1[T](strings.Split(str, sep))
}

func SplitToFloats2[T Float](str, sep string, ts *[]T) error {
	if sep == "." {
		return errors.Error("can not use '.' to split the floats in the string")
	}
	return StringsToFloats2(strings.Split(str, sep), ts)
}

func IntegersJoinToString[T Integer](ts []T, sep string) string {
	if len(ts) == 0 {
		return ""
	}
	var ss = make([]string, len(ts))
	for i, t := range ts {
		ss[i] = IntegerToString(t)
	}
	return strings.Join(ss, sep)
}

func NumbersConvert1[T1 Number, T2 Number](s1 []T1, s2 *[]T2) {
	*s2 = make([]T2, len(s1))
	for i := range s1 {
		(*s2)[i] = T2(s1[i])
	}
}

func NumbersConvert2[T1 Number, T2 Number](s1 []T1, _ T2) (s2 []T2) {
	NumbersConvert1(s1, &s2)
	return
}
