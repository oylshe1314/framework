package transport

import (
	"io"
)

const (
	MessageHeadLength = 8
)

var (
	MessageMaxLength = 1024 * 1024 * 8
)

type Message interface {
	Command() uint32
	Length() uint32
	Body() io.Reader

	Read(any) error
}

type MessageHandler interface {
	Handle(Conn, Message)
}

type MessageHandleFunc func(Conn, Message)

func (fun MessageHandleFunc) Handle(conn Conn, msg Message) {
	fun(conn, msg)
}
