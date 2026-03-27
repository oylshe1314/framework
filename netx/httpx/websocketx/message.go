package websocketx

import (
	"bytes"
	"io"
)

type message struct {
	conn *conn

	modId uint16
	msgId uint16
	body  []byte
}

func (this *message) ModId() uint16 {
	return this.modId
}

func (this *message) MsgId() uint16 {
	return this.msgId
}

func (this *message) Body() io.Reader {
	return bytes.NewBuffer(this.body)
}

func (this *message) Read(v any) error {
	return this.conn.codec.Decode(this.body, v)
}

func (this *message) Reply(v any) error {
	return this.conn.Write(this.modId, this.msgId, v)
}
