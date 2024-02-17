package util

func MapKeys[K comparable, V any](m map[K]V) []K {
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

func MapValues[K comparable, V any](m map[K]V) []V {
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
