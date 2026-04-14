package slicex

func FindIndex[T any](ss []T, test func(int) bool) int {
	if ss == nil {
		return -1
	}
	for i := range ss {
		if test(i) {
			return i
		}
	}
	return -1
}

func FindValue[T any](ss []T, test func(int) bool) (t T) {
	if ss == nil {
		return
	}
	for i := range ss {
		if test(i) {
			t = ss[i]
			return
		}
	}
	return
}
