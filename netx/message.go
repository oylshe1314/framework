package netx

import (
	"bytes"
	"io"

	"github.com/oylshe1314/framework/netx/codec"
)

type message struct {
	cmd  uint32
	body []byte

	codec codec.Codec
}

func newMessage(cmd uint32, body []byte, codec codec.Codec) *message {
	return &message{cmd: cmd, body: body, codec: codec}
}

func (this *message) Command() uint32 {
	return this.cmd
}

func (this *message) Length() uint32 {
	return uint32(len(this.body))
}

func (this *message) Body() io.Reader {
	return bytes.NewReader(this.body)
}

func (this *message) Read(v any) error {
	return this.codec.Decode(this.body, v)
}
