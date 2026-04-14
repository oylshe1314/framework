package slicex

func Set[T any](s []T, t T) {
	for i := range s {
		s[i] = t
	}
}

func Convert[S, R any](ss []S, converter func(s S) R) []R {
	if ss == nil {
		return nil
	}

	var rs = make([]R, len(ss))
	for i := range ss {
		rs[i] = converter(ss[i])
	}
	return rs
}

func ToMap[T any, K comparable](ss []T, f func(T) K) map[K]T {
	if ss == nil {
		return nil
	}

	var m = make(map[K]T)
	for i := range ss {
		m[f(ss[i])] = ss[i]
	}
	return m
}
