package util

func If[T any](ok bool, v1, v2 T) T {
	if ok {
		return v1
	} else {
		return v2
	}
}
