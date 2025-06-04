package util

import (
	"crypto/sha256"
	"encoding/hex"
	"math/rand/v2"
	"strings"
	"sync"
)

type SafeRand struct {
	r *rand.Rand
	l sync.Locker
}

func NewSafeRand() *SafeRand {
	return &SafeRand{
		r: NewRand(),
		l: &sync.Mutex{},
	}
}

func (this *SafeRand) Int64() int64 {
	this.l.Lock()
	defer this.l.Unlock()
	return this.r.Int64()
}

func (this *SafeRand) Uint32() uint32 {
	this.l.Lock()
	defer this.l.Unlock()
	return this.r.Uint32()
}

func (this *SafeRand) Uint64() uint64 {
	this.l.Lock()
	defer this.l.Unlock()
	return this.r.Uint64()
}

func (this *SafeRand) Int32() int32 {
	this.l.Lock()
	defer this.l.Unlock()
	return this.r.Int32()
}

func (this *SafeRand) Int() int {
	this.l.Lock()
	defer this.l.Unlock()
	return this.r.Int()
}

func (this *SafeRand) Uint() uint {
	this.l.Lock()
	defer this.l.Unlock()
	return this.r.Uint()
}

func (this *SafeRand) Int64N(n int64) int64 {
	this.l.Lock()
	defer this.l.Unlock()
	return this.r.Int64N(n)
}

func (this *SafeRand) Uint64N(n uint64) uint64 {
	this.l.Lock()
	defer this.l.Unlock()
	return this.r.Uint64N(n)
}

func (this *SafeRand) Int32N(n int32) int32 {
	this.l.Lock()
	defer this.l.Unlock()
	return this.r.Int32N(n)
}

func (this *SafeRand) Uint32N(n uint32) uint32 {
	this.l.Lock()
	defer this.l.Unlock()
	return this.r.Uint32N(n)
}

func (this *SafeRand) IntN(n int) int {
	this.l.Lock()
	defer this.l.Unlock()
	return this.r.IntN(n)
}

func (this *SafeRand) UintN(n uint) uint {
	this.l.Lock()
	defer this.l.Unlock()
	return this.r.UintN(n)
}

func (this *SafeRand) Float64() float64 {
	this.l.Lock()
	defer this.l.Unlock()
	return this.r.Float64()
}

func (this *SafeRand) Float32() float32 {
	this.l.Lock()
	defer this.l.Unlock()
	return this.r.Float32()
}

func (this *SafeRand) Perm(n int) []int {
	this.l.Lock()
	defer this.l.Unlock()
	return this.r.Perm(n)
}

func (this *SafeRand) Shuffle(n int, swap func(i, j int)) {
	this.l.Lock()
	defer this.l.Unlock()
	this.r.Shuffle(n, swap)
}

var DefaultRandom = NewSafeRand()

func RandomToken() string {
	var src []byte
	var now = UnixMicro()
	src = Uint64ToBytes(src, uint64(now))
	src = Uint64ToBytes(src, DefaultRandom.Uint64())
	src = Uint64ToBytes(src, DefaultRandom.Uint64())
	src = Uint64ToBytes(src, DefaultRandom.Uint64())
	var h = sha256.New()
	h.Write(src)
	return strings.ToUpper(hex.EncodeToString(h.Sum(nil)))
}

func RandomUid() uint64 {
	return uint64(Unix()&0xFFFFFFFFFFFF)<<16 | DefaultRandom.Uint64N(0x10000)
}
