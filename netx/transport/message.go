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
	ModId() uint16
	MsgId() uint16
	Body() io.Reader

	Read(any) error
	Reply(any) error
}

type MessageHandler interface {
	Handle(Conn, Message)
}

type MessageHandlerFunc func(Conn, Message)

func (fun MessageHandlerFunc) Handle(conn Conn, msg Message) {
	fun(conn, msg)
}
