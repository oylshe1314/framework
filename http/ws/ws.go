package ws

import (
	"github.com/gorilla/websocket"
	"github.com/oylshe1314/framework/log"
	"github.com/oylshe1314/framework/message"
	"github.com/oylshe1314/framework/util"
	"io"
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
}

func newMessage(modId, msgId uint16, length uint32, body []byte, conn *Conn) *Message {
	return &Message{ModId: modId, MsgId: msgId, Length: length, Body: body, Conn: conn}
}

func (this *Message) Read(v interface{}) error {
	if v == nil || len(this.Body) == 0 {
		if this.Conn.logger.IsDebugEnabled() {
			this.Conn.logger.Debugf("[%s:%d] <- ModId: %d, MsgId: %d, Msg: %s", this.Conn.RemoteAddr(), this.Conn.ObjectUid(), this.ModId, this.MsgId, util.ToJsonString(nil))
		}
		return nil
	}

	var err = this.Conn.handler.messageCodec().Decode(this.Body, v)
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
	handleWsMessage(*Message)
	handleWsConnect(*Conn)
	handleWsDisconnect(*Conn)
	messageCodec() message.Codec
}

type Conn struct {
	conn *websocket.Conn

	closed bool

	locker sync.Mutex
	logger log.Logger

	handler Handler

	object interface{}
}

func NewConn(conn *websocket.Conn, logger log.Logger, handler Handler) *Conn {
	return &Conn{conn: conn, logger: logger, handler: handler}
}

func (this *Conn) Read() (*Message, error) {
	_, msg, err := this.conn.ReadMessage()
	if err != nil {
		return nil, err
	}

	return newMessage(util.BytesToUint16(msg[0:2]), util.BytesToUint16(msg[2:4]), util.BytesToUint32(msg[4:8]), msg[8:], this), nil
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

func (this *Conn) send(modId, msgId uint16, body []byte) (err error) {
	var msg = make([]byte, HeaderLength+uint32(len(body)))

	util.PutUint16ToBytes(msg[0:2], modId)
	util.PutUint16ToBytes(msg[2:4], msgId)
	util.PutUint32ToBytes(msg[4:8], uint32(len(body)))

	if len(body) > 0 {
		copy(msg[8:], body)
	}

	this.locker.Lock()
	defer this.locker.Unlock()
	return this.conn.WriteMessage(websocket.BinaryMessage, msg)
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
		this.handler.handleWsDisconnect(this)

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

	this.handler.handleWsConnect(this)
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

		this.handler.handleWsMessage(msg)
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

func (this *ConnMux) WsConnectHandler(handler func(*Conn)) {
	this.connectHandler = handler
}

func (this *ConnMux) WsDisconnectHandler(handler func(*Conn)) {
	this.disconnectHandler = handler
}

func (this *ConnMux) WsMessageHandler(modId, msgId uint16, handler MessageHandler) {
	if this.messageHandlers == nil {
		this.messageHandlers = make(map[uint32]MessageHandler)
	}
	this.messageHandlers[util.Compose2uint16(modId, msgId)] = handler
}

func (this *ConnMux) WsDefaultHandler(handler MessageHandler) {
	this.defaultHandler = handler
}

func (this *ConnMux) handleWsMessage(msg *Message) {
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

func (this *ConnMux) handleWsConnect(conn *Conn) {
	if this.connectHandler != nil {
		this.connectHandler(conn)
	}
}

func (this *ConnMux) handleWsDisconnect(conn *Conn) {
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
