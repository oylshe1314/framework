package store

type Pair[Key comparable, Value any] struct {
	Key   Key   `json:"key"`
	Value Value `json:"value"`
}

func NewPair[Key comparable, Value any](key Key, value Value) *Pair[Key, Value] {
	return &Pair[Key, Value]{
		Key:   key,
		Value: value,
	}
}
