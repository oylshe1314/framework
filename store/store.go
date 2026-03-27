package store

type Store[Key comparable, Value any] interface {
	Put(key Key, value Value)
	Get(key Key) (Value, bool)
	Remove(key Key)
	Clear()
	Len() int
	Foreach(func(key Key, value Value))
}

type store[Key comparable, Value any] struct {
	m map[Key]Value
}

func New[Key comparable, Value any]() Store[Key, Value] {
	return &store[Key, Value]{m: make(map[Key]Value)}
}

func (this *store[Key, Value]) Put(key Key, value Value) {
	this.m[key] = value
}

func (this *store[Key, Value]) Get(key Key) (Value, bool) {
	v, ok := this.m[key]
	return v, ok
}

func (this *store[Key, Value]) Remove(key Key) {
	delete(this.m, key)
}

func (this *store[Key, Value]) Clear() {
	this.m = make(map[Key]Value)
}

func (this *store[Key, Value]) Len() int {
	return len(this.m)
}

func (this *store[Key, Value]) Foreach(f func(key Key, value Value)) {
	for k, v := range this.m {
		f(k, v)
	}
}
