package util

type Deque[T any] interface {
	PushFront(v T)
	PushBack(v T)

	PopFront() (T, bool)
	PopBack() (T, bool)

	PeekFront() (T, bool)
	PeekBack() (T, bool)

	Len() int
	Clear()

	Foreach(func(v T))
}

type dequeNode[T any] struct {
	value T
	prev  *dequeNode[T]
	next  *dequeNode[T]
}

type deque[T any] struct {
	size int

	head *dequeNode[T]
	tail *dequeNode[T]
}

func NewDeque[T any]() Deque[T] {
	return &deque[T]{}
}

func (this *deque[T]) PushFront(v T) {
	if this.head == nil {
		this.head = &dequeNode[T]{value: v}
		this.tail = this.head
	} else {
		this.head.prev = &dequeNode[T]{value: v, next: this.head}
		this.head = this.head.prev
	}
	this.size += 1
}

func (this *deque[T]) PushBack(v T) {
	if this.tail == nil {
		this.tail = &dequeNode[T]{value: v}
		this.head = this.tail
	} else {
		this.tail.next = &dequeNode[T]{value: v, prev: this.tail}
		this.tail = this.tail.next
	}
	this.size += 1
}

func (this *deque[T]) PopFront() (v T, ok bool) {
	if this.tail == nil {
		return v, false
	} else {
		v = this.head.value
		this.head = this.head.next
		if this.head == nil {
			this.tail = nil
		} else {
			this.head.prev = nil
		}
		this.size -= 1
		return v, true
	}
}

func (this *deque[T]) PopBack() (v T, ok bool) {
	if this.head == nil {
		return v, false
	} else {
		v = this.tail.value
		this.tail = this.tail.prev
		if this.tail == nil {
			this.head = nil
		} else {
			this.tail.next = nil
		}
		this.size -= 1
		return v, true
	}
}

func (this *deque[T]) PeekFront() (v T, ok bool) {
	if this.tail == nil {
		return v, false
	} else {
		return this.head.value, true
	}
}

func (this *deque[T]) PeekBack() (v T, ok bool) {
	if this.head == nil {
		return v, false
	} else {
		return this.tail.value, true
	}
}

func (this *deque[T]) Len() int {
	return this.size
}

func (this *deque[T]) Clear() {
	this.size = 0
	this.head = nil
	this.tail = nil
}

func (this *deque[T]) Foreach(f func(v T)) {
	for n := this.head; n != nil; n = n.next {
		f(n.value)
	}
}
