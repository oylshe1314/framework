package netx

import (
	"encoding/binary"
	"framework/errors"
	"framework/log"
	"framework/netx/codec"
	"framework/netx/transport"
	"framework/store"
	"io"
	"net"
	"runtime/debug"
)

type conn struct {
	conn net.Conn

	logger log.Logger

	codec codec.Codec

	attr store.Attribute

	closed bool
}

func newConn(c net.Conn, logger log.Logger, codec codec.Codec) *conn {
	return &conn{conn: c, logger: logger, codec: codec}
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
	var head = make([]byte, transport.MessageHeadLength)
	_, err := this.conn.Read(head)
	if err != nil {
		return nil, err
	}

	var modId = binary.LittleEndian.Uint16(head[:2])
	var msgId = binary.LittleEndian.Uint16(head[2:4])
	var length = binary.LittleEndian.Uint32(head[4:8])

	var body []byte
	if length > 0 {
		if length >= uint32(transport.MessageMaxLength) {
			return nil, errors.New("message too long")
		}

		body = make([]byte, length)
		_, err = this.conn.Read(body)
		if err != nil {
			return nil, err
		}
	}
	return &message{conn: this, modId: modId, msgId: msgId, body: body}, nil
}

func (this *conn) Write(modId, msgId uint16, v any) error {
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

	_, err = this.conn.Write(buf)
	return err
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
