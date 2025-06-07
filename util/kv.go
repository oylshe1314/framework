package util

type Pair[K comparable, V any] struct {
	Key   K `json:"key" bson:"key"`
	Value V `json:"value" bson:"value"`
}

func NewPair[K comparable, V any](k K, v V) *Pair[K, V] {
	return &Pair[K, V]{Key: k, Value: v}
}

type Pairs[K comparable, V any] []*Pair[K, V]

func (pairs Pairs[K, V]) Append(k K, v V) Pairs[K, V] {
	return append(pairs, NewPair(k, v))
}

func (pairs Pairs[K, V]) ToMap() map[K]V {
	if pairs == nil {
		return nil
	}

	var m = map[K]V{}
	for _, pair := range pairs {
		m[pair.Key] = pair.Value
	}
	return m
}
