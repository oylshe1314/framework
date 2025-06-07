package util

import "sync"

type Queue[T any] interface {
	Push(t T)
	Pop() T
	Head() T
	Len() int
}

type node[T any] struct {
	t T
	p *node[T]
	n *node[T]
}

type linkedQueue[T any] struct {
	l int      //length
	h *node[T] //head
	t *node[T] //tail
}

func NewLinkedQueue[T any]() Queue[T] {
	return &linkedQueue[T]{}
}

func (this *linkedQueue[T]) Push(t T) {
	if this.t == nil {
		this.h = &node[T]{t: t}
		this.t = this.h
	} else {
		this.t.n = &node[T]{t: t, p: this.t}
		this.t = this.t.n
	}
	this.l += 1
}

func (this *linkedQueue[T]) Pop() (t T) {
	if this.h != nil {
		t = this.h.t
		this.h = this.h.n
		if this.h != nil {
			this.h.p = nil
		}
		this.l -= 1
	}
	return
}

func (this *linkedQueue[T]) Head() (t T) {
	if this.h != nil {
		t = this.h.t
	}
	return
}

func (this *linkedQueue[T]) Len() int {
	return this.l
}

type safeQueue[T any] struct {
	m sync.Mutex
	linkedQueue[T]
}

func NewSafeQueue[T any]() Queue[T] {
	return &safeQueue[T]{}
}
func (this *safeQueue[T]) Push(t T) {
	this.m.Lock()
	defer this.m.Unlock()
	this.linkedQueue.Push(t)
}

func (this *safeQueue[T]) Pop() (t T) {
	this.m.Lock()
	defer this.m.Unlock()
	return this.linkedQueue.Pop()
}

func (this *safeQueue[T]) Head() (t T) {
	this.m.Lock()
	defer this.m.Unlock()
	return this.linkedQueue.Head()
}

func (this *safeQueue[T]) Len() int {
	this.m.Lock()
	defer this.m.Unlock()
	return this.linkedQueue.Len()
}
