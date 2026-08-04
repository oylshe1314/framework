package util

import (
	"cmp"
)

type LinkedMap[K cmp.Ordered, V any] interface {
	Put(K, V)
	Get(K) (V, bool)

	Remove(K) (V, bool)

	Len() int
	Clear()

	ForEach(func(K, V))
}

type mapNode[K cmp.Ordered, V any] struct {
	key   K
	value V
}

type linkedMap[K cmp.Ordered, V any] struct {
	deque *deque[*mapNode[K, V]]
	nodes map[K]*dequeNode[*mapNode[K, V]]
}

func NewLinkedMap[K cmp.Ordered, V any]() LinkedMap[K, V] {
	return &linkedMap[K, V]{deque: &deque[*mapNode[K, V]]{}, nodes: map[K]*dequeNode[*mapNode[K, V]]{}}
}

func (this *linkedMap[K, V]) Put(k K, v V) {
	this.deque.PushBack(&mapNode[K, V]{key: k, value: v})
	this.nodes[k] = this.deque.tail

}

func (this *linkedMap[K, V]) Get(k K) (v V, ok bool) {
	var n *dequeNode[*mapNode[K, V]]
	n, ok = this.nodes[k]
	if !ok {
		return v, false
	} else {
		return n.value.value, true
	}
}

func (this *linkedMap[K, V]) Remove(k K) (v V, ok bool) {
	var n *dequeNode[*mapNode[K, V]]
	n, ok = this.nodes[k]
	if !ok {
		return v, false
	} else {
		switch {
		case n.prev == nil && n.next == nil:
			this.deque.head = nil
			this.deque.tail = nil
		case n.prev == nil:
			this.deque.head = n.next
		case n.next == nil:
			this.deque.tail = n.prev
		default:
			n.prev.next = n.next
			n.next.prev = n.prev
		}
		delete(this.nodes, k)
		return n.value.value, true
	}
}

func (this *linkedMap[K, V]) Len() int {
	return this.deque.Len()
}

func (this *linkedMap[K, V]) Clear() {
	this.deque.Clear()
	clear(this.nodes)
}

func (this *linkedMap[K, V]) ForEach(f func(K, V)) {
	this.deque.Foreach(func(n *mapNode[K, V]) {
		f(n.key, n.value)
	})
}
