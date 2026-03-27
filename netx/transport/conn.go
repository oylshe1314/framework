package transport

import (
	"framework/store"
	"io"
	"net"
)

type Conn interface {
	io.Closer
	store.AttributeProvider

	LocalAddr() net.Addr
	RemoteAddr() net.Addr

	Read() (Message, error)
	Write(uint16, uint16, any) error

	Serve(handler MessageHandler) error
}
