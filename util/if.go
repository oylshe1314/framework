package util

func If[T any](ok bool, v1, v2 T) T {
	if ok {
		return v1
	} else {
		return v2
	}
}

func Iff[T any](ok bool, f1, f2 func() T) T {
	if ok {
		return f1()
	} else {
		return f2()
	}
}
