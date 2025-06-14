package util

func SliceSet[T any](s []T, t T) {
	for i := range s {
		s[i] = t
	}
}

func SliceConvert[S, R any](s []S, converter func(s S) R) []R {
	if s == nil {
		return nil
	}

	var r = make([]R, len(s))
	for i := range s {
		r[i] = converter(s[i])
	}
	return r
}

func SliceFindValue[T any](s []T, cmp func(int) bool) (t T) {
	for i := range s {
		if cmp(i) {
			t = s[i]
			return
		}
	}
	return
}
