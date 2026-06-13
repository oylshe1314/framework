package util

import "testing"

func TestNewLinkedMap(t *testing.T) {
	var m = NewLinkedMap[string, any]()

	m.Put("Name", "Lisa")
	m.Put("Age", 18)
	m.Put("Fuck", true)
	m.Put("Times", 100)
	m.Put("BoyFriend", "SK")

	t.Log("Len: ", m.Len())
	m.ForEach(func(k string, v interface{}) {
		t.Log("key: ", k, "value: ", v)
	})
}
