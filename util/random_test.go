package util

import (
	"testing"
)

func TestRandomToken(t *testing.T) {
	for range 100 {
		t.Log(RandomToken())
	}
}

func TestRandomUid(t *testing.T) {
	for range 100 {
		t.Log(RandomUid())
	}
}
