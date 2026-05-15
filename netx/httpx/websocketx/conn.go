package websocketx

import (
	"encoding/binary"
	"io"
	"net"
	"runtime/debug"

	"github.com/oylshe1314/framework/errors"
	"github.com/oylshe1314/framework/log"
	"github.com/oylshe1314/framework/netx/codec"
	"github.com/oylshe1314/framework/netx/transport"
	"github.com/oylshe1314/framework/store"

	"github.com/gorilla/websocket"
)

type conn struct {
	conn *websocket.Conn

	logger log.Logger

	codec codec.Codec

	attr store.Attribute

	closed bool
}

func newConn(cc *websocket.Conn, logger log.Logger, codec codec.Codec) *conn {
	return &conn{conn: cc, logger: logger, codec: codec}
}

func (this *conn) Close() error {
	this.closed = true
	return this.conn.Close()
}

func (this *conn) LocalAddr() net.Addr {
	return this.conn.LocalAddr()
}

func (this *conn) RemoteAddr() net.Addr {
	return this.conn.RemoteAddr()
}

func (this *conn) Attribute() store.Attribute {
	return this.attr
}

func (this *conn) read() (uint32, []byte, error) {
	_, buf, err := this.conn.ReadMessage()
	if err != nil {
		return 0, nil, err
	}

	var cmd = binary.LittleEndian.Uint32(buf[:4])
	var length = binary.LittleEndian.Uint32(buf[4:8])
	if length == 0 {
		return cmd, nil, nil
	}

	return cmd, buf[8:], nil
}

func (this *conn) write(cmd uint32, body []byte) error {
	var buf = make([]byte, transport.MessageHeadLength+len(body))

	binary.LittleEndian.PutUint32(buf[:4], cmd)
	binary.LittleEndian.PutUint32(buf[4:8], uint32(len(body)))

	if len(body) > 0 {
		copy(buf[8:], body)
	}

	return this.conn.WriteMessage(websocket.BinaryMessage, buf)
}

func (this *conn) Read() (transport.Message, error) {
	cmd, body, err := this.read()
	if err != nil {
		return nil, err
	}
	return newMessage(cmd, body, this.codec), nil
}

func (this *conn) Send(cmd uint32, v any) error {
	if v == nil {
		return this.write(cmd, nil)
	}

	body, err := this.codec.Encode(v)
	if err != nil {
		return err
	}
	return this.write(cmd, body)
}

func (this *conn) Serve(handleFunc transport.MessageHandleFunc) error {
	defer func() {
		if this.closed {
			return
		}
		_ = this.Close()

		var err = recover()
		if err != nil {
			this.logger.Error(err)
			this.logger.Error(string(debug.Stack()))
		}
	}()

	for {
		msg, err := this.Read()
		if err != nil {
			if this.closed || errors.Is(err, io.EOF) {
				return nil
			}

			this.logger.Error(err)
			return err
		}

		handleFunc(this, msg)
	}
}
