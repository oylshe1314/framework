package store

import "sync"

type concurrentStore[Key comparable, Value any] struct {
	sync.RWMutex
	store[Key, Value]
}

func NewConcurrent[Key comparable, Value any]() Store[Key, Value] {
	return &concurrentStore[Key, Value]{
		RWMutex: sync.RWMutex{},
		store: store[Key, Value]{
			m: make(map[Key]Value),
		},
	}
}

func (this *concurrentStore[Key, Value]) Put(key Key, value Value) {
	this.Lock()
	defer this.Unlock()
	this.store.Put(key, value)
}

func (this *concurrentStore[Key, Value]) Get(key Key) (Value, bool) {
	this.RLock()
	defer this.RUnlock()
	return this.store.Get(key)
}

func (this *concurrentStore[Key, Value]) Remove(key Key) {
	this.Lock()
	defer this.Unlock()
	this.store.Remove(key)
}

func (this *concurrentStore[Key, Value]) Clear() {
	this.Lock()
	defer this.Unlock()
	this.store.Clear()
}

func (this *concurrentStore[Key, Value]) Len() int {
	this.RLock()
	defer this.RUnlock()
	return this.store.Len()
}

func (this *concurrentStore[Key, Value]) Foreach(f func(key Key, value Value)) {
	this.Lock()
	defer this.Unlock()
	this.store.Foreach(f)
}
