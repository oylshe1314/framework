package net

import (
	"github.com/oylshe1314/framework/log"
	"github.com/oylshe1314/framework/message"
	"github.com/oylshe1314/framework/util"
	"io"
	"net"
	"runtime/debug"
	"sync"
)

const HeaderLength uint32 = 8

type Message struct {
	ModId  uint16
	MsgId  uint16
	Length uint32
	Body   []byte

	Conn *Conn

	codec message.Codec
}

func newMessage(modId, msgId uint16, length uint32, body []byte, conn *Conn, codec message.Codec) *Message {
	return &Message{ModId: modId, MsgId: msgId, Length: length, Body: body, Conn: conn, codec: codec}
}

func (this *Message) Read(v interface{}) error {
	if v == nil || len(this.Body) == 0 {
		if this.Conn.logger.IsDebugEnabled() {
			this.Conn.logger.Debugf("[%s:%d] <- ModId: %d, MsgId: %d, Msg: %s", this.Conn.RemoteAddr(), this.Conn.ObjectUid(), this.ModId, this.MsgId, util.ToJsonString(nil))
		}
		return nil
	}

	var err = this.codec.Decode(this.Body, v)
	if err != nil {
		return err
	}

	if this.Conn.logger.IsDebugEnabled() {
		this.Conn.logger.Debugf("[%s:%d] <- ModId: %d, MsgId: %d, Msg: %s", this.Conn.RemoteAddr(), this.Conn.ObjectUid(), this.ModId, this.MsgId, util.ToJsonString(v))
	}
	return nil
}

func (this *Message) Reply(v interface{}) error {
	return this.Conn.Send(this.ModId, this.MsgId, v)
}

type MessageHandler func(*Message)

type Handler interface {
	handleMessage(*Message)
	handleConnect(*Conn)
	handleDisconnect(*Conn)
	messageCodec() message.Codec
}

type Conn struct {
	conn net.Conn

	closed bool

	locker sync.Mutex
	logger log.Logger

	handler Handler

	object interface{}
}

func NewConn(conn net.Conn, logger log.Logger, handler Handler) *Conn {
	return &Conn{conn: conn, logger: logger, handler: handler}
}

func (this *Conn) LocalAddr() string {
	return this.conn.LocalAddr().String()
}

func (this *Conn) RemoteAddr() string {
	return this.conn.RemoteAddr().String()
}

func (this *Conn) BindObject(object interface{}) {
	this.object = object
}

func (this *Conn) ClearObject() {
	this.object = nil
}

func (this *Conn) Object() interface{} {
	return this.object
}

func (this *Conn) ObjectUid() (uid uint64) {
	obj, ok := this.Object().(interface{ Uid() uint64 })
	if ok {
		uid = obj.Uid()
	}
	return
}

func (this *Conn) Read() (msg *Message, err error) {
	var head = make([]byte, HeaderLength)
	_, err = io.ReadFull(this.conn, head)
	if err != nil {
		return
	}

	var modId = util.BytesToUint16(head[0:2])
	var msgId = util.BytesToUint16(head[2:4])
	var length = util.BytesToUint32(head[4:8])

	var body []byte
	if length > 0 {
		body = make([]byte, length)
		_, err = io.ReadFull(this.conn, body)
		if err != nil {
			return
		}
	}

	msg = newMessage(modId, msgId, length, body, this, this.handler.messageCodec())
	return
}

func (this *Conn) send(modId, msgId uint16, body []byte) (err error) {
	var head = make([]byte, HeaderLength)

	util.PutUint16ToBytes(head[0:2], modId)
	util.PutUint16ToBytes(head[2:4], msgId)
	util.PutUint32ToBytes(head[4:8], uint32(len(body)))

	this.locker.Lock()
	defer this.locker.Unlock()

	_, err = this.conn.Write(head)
	if err != nil {
		return err
	}

	if len(body) == 0 {
		return nil
	}

	_, err = this.conn.Write(body)
	return err
}

func (this *Conn) Send(modId, msgId uint16, v interface{}) (err error) {
	if this.logger.IsDebugEnabled() {
		this.logger.Debugf("[%s:%d] -> ModId: %d, MsgId: %d, Msg: %s", this.RemoteAddr(), this.ObjectUid(), modId, msgId, util.ToJsonString(v))
	}

	body, err := this.handler.messageCodec().Encode(v)
	if err != nil {
		return err
	}
	return this.send(modId, msgId, body)
}

func (this *Conn) Serve() error {
	defer func() {
		if this.closed {
			return
		}

		this.handler.handleDisconnect(this)
		_ = this.conn.Close()

		var err = recover()
		if err != nil {
			this.logger.Error(err)
			this.logger.Error(string(debug.Stack()))
		}
	}()

	this.handler.handleConnect(this)
	for {
		msg, err := this.Read()
		if err != nil {
			if err == io.EOF {
				return nil
			}

			if this.closed {
				return nil
			}

			this.logger.Error("Read message failed, ", err)
			return err
		}

		this.handler.handleMessage(msg)
	}
}

func (this *Conn) Close() (err error) {
	this.closed = true
	return this.conn.Close()
}

type ConnMux struct {
	codec message.Codec

	connectHandler    func(*Conn)
	disconnectHandler func(*Conn)
	defaultHandler    MessageHandler
	messageHandlers   map[uint32]MessageHandler
}

func (this *ConnMux) ConnectHandler(handler func(*Conn)) {
	this.connectHandler = handler
}

func (this *ConnMux) DisconnectHandler(handler func(*Conn)) {
	this.disconnectHandler = handler
}

func (this *ConnMux) MessageHandler(modId, msgId uint16, handler MessageHandler) {
	if this.messageHandlers == nil {
		this.messageHandlers = make(map[uint32]MessageHandler)
	}
	this.messageHandlers[util.Compose2uint16(modId, msgId)] = handler
}

func (this *ConnMux) DefaultHandler(handler MessageHandler) {
	this.defaultHandler = handler
}

func (this *ConnMux) handleMessage(msg *Message) {
	if this.messageHandlers == nil {
		return
	}

	var handler = this.messageHandlers[util.Compose2uint16(msg.ModId, msg.MsgId)]
	if handler != nil {
		handler(msg)
	} else {
		if this.defaultHandler != nil {
			this.defaultHandler(msg)
		}
	}
}

func (this *ConnMux) handleConnect(conn *Conn) {
	if this.connectHandler != nil {
		this.connectHandler(conn)
	}
}

func (this *ConnMux) handleDisconnect(conn *Conn) {
	if this.disconnectHandler != nil {
		this.disconnectHandler(conn)
	}
}

func (this *ConnMux) SetCodec(codec message.Codec) {
	this.codec = codec
}

func (this *ConnMux) messageCodec() message.Codec {
	if this.codec == nil {
		return message.DefaultCodec
	}
	return this.codec
}
