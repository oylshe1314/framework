package util

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
	l int
	h *node[T]
	t *node[T]
}

func NewLinkedQueue[T any]() Queue[T] {
	return &linkedQueue[T]{}
}

func (this *linkedQueue[T]) Push(t T) {
	if this.h == nil {
		this.t = &node[T]{t: t, p: nil, n: nil}
		this.h = this.t
	} else {
		this.t.n = &node[T]{t: t, p: this.t, n: nil}
		this.t = this.t.n
	}
	this.l += 1
}

func (this *linkedQueue[T]) Pop() (t T) {
	if this.h == nil {
		return
	}

	var head = this.h
	this.h = head.n
	if this.h != nil {
		this.h.p = nil
	}

	this.l -= 1

	return head.t
}

func (this *linkedQueue[T]) Head() (t T) {
	if this.h == nil {
		return
	}
	return this.h.t
}

func (this *linkedQueue[T]) Len() int {
	return this.l
}
