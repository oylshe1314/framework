package message

import (
	"fmt"
	"testing"
)

func TestStringCodec(t *testing.T) {

	var c = DefaultCodec

	buf, err := c.Encode("example codec")
	if err != nil {
		t.Error(err)
		return
	}

	fmt.Println("Encoded:", buf)

	var s string
	err = c.Decode(buf, &s)
	if err != nil {
		t.Error(err)
		return
	}

	fmt.Println("Decoded：", s)
}

func TestJsonCodec(t *testing.T) {

	var c = NewJsonCodec()

	buf, err := c.Encode("{\"content\":\"example codec\"}")
	if err != nil {
		t.Error(err)
		return
	}

	fmt.Println("Encoded:", buf)

	var s string
	err = c.Decode(buf, &s)
	if err != nil {
		t.Error(err)
		return
	}

	fmt.Println("Decoded：", s)
}
