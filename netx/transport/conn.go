package transport

import (
	"io"
	"net"

	"github.com/oylshe1314/framework/store"
)

type Conn interface {
	io.Closer
	store.AttributeProvider

	LocalAddr() net.Addr
	RemoteAddr() net.Addr

	Read() (Message, error)

	Send(cmd uint32, v any) error

	Serve(handler MessageHandler) error
}
