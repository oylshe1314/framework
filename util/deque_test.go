package util

import "testing"

func TestNewDeque(t *testing.T) {
	var dq = NewDeque[int]()

	var v int
	var ok bool

	v, ok = dq.PopFront()
	t.Log("PopFront: ", v, ok)

	v, ok = dq.PopBack()
	t.Log("PopBack: ", v, ok)

	dq.PushBack(3)
	dq.PushBack(4)
	dq.PushBack(5)
	dq.PushFront(2)
	dq.PushFront(1)

	t.Log("Len: ", dq.Len())
	dq.Foreach(func(n int) {
		t.Log("Value: ", n)
	})

	v, ok = dq.PopFront()
	t.Log("PopFront: ", v, ok)

	t.Log("Len: ", dq.Len())
	dq.Foreach(func(n int) {
		t.Log("Value: ", n)
	})

	v, ok = dq.PopBack()
	t.Log("PopBack: ", v, ok)

	t.Log("Len: ", dq.Len())
	dq.Foreach(func(n int) {
		t.Log("Value: ", n)
	})
}
