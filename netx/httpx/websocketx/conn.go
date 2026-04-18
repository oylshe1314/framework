package websocketx

import (
	"encoding/binary"
	"github.com/oylshe1314/framework/errors"
	"github.com/oylshe1314/framework/log"
	"github.com/oylshe1314/framework/netx/codec"
	"github.com/oylshe1314/framework/netx/transport"
	"github.com/oylshe1314/framework/store"
	"io"
	"net"
	"runtime/debug"

	"github.com/gorilla/websocket"
)

type conn struct {
	conn *websocket.Conn

	logger log.Logger

	codec codec.Codec

	attr store.Attribute

	closed bool

	remain []byte
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

func (this *conn) Read() (transport.Message, error) {
	for {
		_, buf, err := this.conn.ReadMessage()
		if err != nil {
			return nil, err
		}

		if len(this.remain) > 0 {
			var newBuf = make([]byte, len(this.remain)+len(buf))
			copy(newBuf, this.remain)
			copy(newBuf[len(this.remain):], buf)
			this.remain = nil
			buf = newBuf
		}

		if len(buf) < transport.MessageHeadLength {
			continue
		}

		var modId = binary.LittleEndian.Uint16(buf[:2])
		var msgId = binary.LittleEndian.Uint16(buf[2:4])
		var length = binary.LittleEndian.Uint32(buf[4:8])
		var body = buf[8:]

		if len(body) < int(length) {
			continue
		}

		if len(body) == int(length) {
			return &message{conn: this, modId: modId, msgId: msgId, body: body}, nil
		} else {
			this.remain = body[length:]
			body = body[:length]
			return &message{conn: this, modId: modId, msgId: msgId, body: body}, nil
		}
	}
}

func (this *conn) Write(modId uint16, msgId uint16, v any) error {
	var err error
	var body []byte

	if v != nil {
		body, err = this.codec.Encode(v)
		if err != nil {
			return err
		}
	}

	var buf = make([]byte, transport.MessageHeadLength+len(body))

	binary.LittleEndian.PutUint16(buf[:2], modId)
	binary.LittleEndian.PutUint16(buf[2:4], msgId)
	binary.LittleEndian.PutUint32(buf[4:8], uint32(len(body)))

	if len(body) > 0 {
		copy(buf[8:], body)
	}

	return this.conn.WriteMessage(websocket.BinaryMessage, buf)
}

func (this *conn) Serve(handler transport.MessageHandler) error {
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

		handler.Handle(this, msg)
	}
}
