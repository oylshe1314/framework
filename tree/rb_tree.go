package tree

import (
	"cmp"
)

type Comparator[Key any] func(k1 Key, k2 Key) int

func compare[T cmp.Ordered](v1, v2 T) int {
	if v1 < v2 {
		return -1
	}
	if v1 > v2 {
		return 1
	}
	return 0
}

func defaultComparator[Key any](k1, k2 Key) int {
	var i1 interface{} = k1
	var i2 interface{} = k2
	switch i1.(type) {
	case int:
		return compare(i1.(int), i2.(int))
	case int8:
		return compare(i1.(int8), i2.(int8))
	case int16:
		return compare(i1.(int16), i2.(int16))
	case int32:
		return compare(i1.(int32), i2.(int32))
	case int64:
		return compare(i1.(int64), i2.(int64))
	case uint:
		return compare(i1.(uint), i2.(uint))
	case uint8:
		return compare(i1.(uint8), i2.(uint8))
	case uint16:
		return compare(i1.(uint16), i2.(uint16))
	case uint32:
		return compare(i1.(uint32), i2.(uint32))
	case uint64:
		return compare(i1.(uint64), i2.(uint64))
	case uintptr:
		return compare(i1.(uintptr), i2.(uintptr))
	case float32:
		return compare(i1.(float32), i2.(float32))
	case float64:
		return compare(i1.(float64), i2.(float64))
	case string:
		return compare(i1.(string), i2.(string))
	default:
		return 0
	}
}

type RBTree[Key any, Value any] interface {
	Put(key Key, value Value) Value
	Get(key Key) Value
	Delete(key Key) Value
	Clear()
	Size() int
	Foreach(f func(key Key, Value Value))
}

type rbNode[Key any, Value any] struct {
	red bool
	key Key
	val Value

	p *rbNode[Key, Value]
	l *rbNode[Key, Value]
	r *rbNode[Key, Value]
}

type rbTree[Key any, Value any] struct {
	size int

	comp Comparator[Key]

	root *rbNode[Key, Value]
	null *rbNode[Key, Value]
}

// NewRBTree 还没写完不能用
func NewRBTree[Key any, Value any](comparator ...Comparator[Key]) RBTree[Key, Value] {
	var tree = &rbTree[Key, Value]{null: &rbNode[Key, Value]{}}
	if len(comparator) > 0 {
		tree.comp = comparator[0]
	} else {
		tree.comp = defaultComparator[Key]
	}
	tree.root = tree.null
	return tree
}

func (this *rbTree[Key, Value]) leftRotate(x *rbNode[Key, Value]) {
	var y = x.r
	x.r = y.l
	if y.l != this.null {
		y.l.p = x
	}
	y.p = x.p
	if x.p == this.null {
		this.root = y
	} else {
		if x == x.p.l {
			x.p.l = y
		} else {
			x.p.r = y
		}
	}
	y.l = x
	x.p = y
}

func (this *rbTree[Key, Value]) rightRotate(x *rbNode[Key, Value]) {
	var y = x.l
	x.l = y.r
	if y.r != this.null {
		y.r.p = x
	}
	y.p = x.p
	if x.p == this.null {
		this.root = y
	} else {
		if x == x.p.l {
			x.p.l = y
		} else {
			x.p.r = y
		}
	}
	y.r = x
	x.p = y
}

func (this *rbTree[Key, Value]) insertFixup(x *rbNode[Key, Value]) {
	for x.p.red {
		if x.p == x.p.p.l {
			var y = x.p.p.r
			if y.red {
				x.p.red = false
				y.red = false
				x.p.p.red = true
				x = x.p.p
			} else {
				if x == x.p.r {
					x = x.p
					this.leftRotate(x)
				}
				x.p.red = false
				x.p.p.red = true
				this.rightRotate(x.p.p)
			}
		} else {
			var y = x.p.p.l
			if y.red {
				x.p.red = false
				y.red = false
				x.p.p.red = true
				x = x.p.p
			} else {
				if x == x.p.l {
					x = x.p
					this.rightRotate(x)
				}
				x.p.red = false
				x.p.p.red = true
				this.leftRotate(x.p.p)
			}
		}
	}
	this.root.red = false
}

func (this *rbTree[Key, Value]) Put(key Key, value Value) (old Value) {
	var x = &this.root
	var p *rbNode[Key, Value]

	for *x != this.null {
		var r = this.comp(key, (*x).key)
		switch {
		case r == 0:
			old = (*x).val
			(*x).val = value
			return
		case r < 0:
			p = *x
			x = &(*x).l
		case r > 0:
			p = *x
			x = &(*x).r
		}
	}

	*x = &rbNode[Key, Value]{red: true, key: key, val: value, p: this.null, l: this.null, r: this.null}
	if p != nil {
		(*x).p = p
	}

	this.size += 1

	this.insertFixup(*x)
	return
}

func (this *rbTree[Key, Value]) Get(key Key) (value Value) {
	var x = this.root
	for x != this.null {
		var r = this.comp(key, x.key)
		switch {
		case r == 0:
			value = x.val
			return
		case r < 0:
			x = x.l
		case r > 0:
			x = x.r
		}
	}
	return
}

func (this *rbTree[Key, Value]) transplant(x, y *rbNode[Key, Value]) {
	if x.p == this.null {
		this.root = y
	} else if x == x.p.l {
		x.p.l = y
	} else {
		x.p.r = y
	}
	if y != this.null {
		y.p = x.p
	}
}

func (this *rbTree[Key, Value]) deleteFixup(x *rbNode[Key, Value]) {

}

func (this *rbTree[Key, Value]) Delete(key Key) (value Value) {
	return
}

func (this *rbTree[Key, Value]) Clear() {
	this.root = this.null
}

func (this *rbTree[Key, Value]) Size() int {
	return this.size
}

func (this *rbTree[Key, Value]) foreach(x *rbNode[Key, Value], f func(key Key, Value Value)) {
	if x != this.null {
		f(x.key, x.val)
		this.foreach(x.l, f)
		this.foreach(x.r, f)
	}
}

func (this *rbTree[Key, Value]) Foreach(f func(key Key, Value Value)) {
	this.foreach(this.root, f)
}
