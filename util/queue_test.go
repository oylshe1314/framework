package util

import (
	"fmt"
	"testing"
)

func TestNewLinkedQueue(t *testing.T) {
	var q = NewLinkedQueue[int]()

	q.Push(1)
	fmt.Println(q.Len())
	q.Push(2)
	fmt.Println(q.Len())
	q.Push(3)
	fmt.Println(q.Len())
	q.Push(4)
	fmt.Println(q.Len())
	q.Push(5)
	fmt.Println(q.Len())

	fmt.Println("len:", q.Len())
	fmt.Println("val:", q.Pop())
	fmt.Println()

	fmt.Println("len:", q.Len())
	fmt.Println("val:", q.Pop())
	fmt.Println()

	fmt.Println("len:", q.Len())
	fmt.Println("val:", q.Pop())
	fmt.Println()

	fmt.Println("len:", q.Len())
	fmt.Println("val:", q.Pop())
	fmt.Println()

	fmt.Println("len:", q.Len())
	fmt.Println("val:", q.Pop())
	fmt.Println()

	fmt.Println("len:", q.Len())
	fmt.Println("val:", q.Pop())
	fmt.Println()

}
