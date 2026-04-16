package mapx

func Keys[K comparable, V any](m map[K]V) []K {
	if len(m) == 0 {
		return nil
	}

	var i = 0
	var keys = make([]K, len(m))
	for k := range m {
		keys[i] = k
		i++
	}
	return keys
}

func Values[K comparable, V any](m map[K]V) []V {
	if len(m) == 0 {
		return nil
	}

	var i = 0
	var values = make([]V, len(m))
	for _, v := range m {
		values[i] = v
		i++
	}
	return values
}

func FindValue[K comparable, V any](m map[K]V, vCmp func(V) bool) (rv V) {
	for _, v := range m {
		if vCmp(v) {
			rv = v
			break
		}
	}
	return
}
