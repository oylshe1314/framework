package route

import (
	"github.com/oylshe1314/framework/netx/transport"
)

type ConnRegistry interface {
	ConnectHandler(func(transport.Conn))
	DisconnectHandler(func(transport.Conn))
	DefaultHandler(transport.MessageHandleFunc)
	MessageHandler(uint32, transport.MessageHandleFunc)
}

type ConnHandler interface {
	HandleConnect(conn transport.Conn)
	HandleDisconnect(conn transport.Conn)
	HandleMessage(conn transport.Conn, msg transport.Message)
}

type ConnMux struct {
	connectHandlers    []func(transport.Conn)
	disconnectHandlers []func(transport.Conn)
	defaultHandler     transport.MessageHandler
	messageHandler     map[uint32]transport.MessageHandler
}

func NewConnMux() *ConnMux {
	return &ConnMux{messageHandler: make(map[uint32]transport.MessageHandler)}
}

func (this *ConnMux) ConnectHandler(handler func(transport.Conn)) {
	this.connectHandlers = append(this.connectHandlers, handler)
}

func (this *ConnMux) DisconnectHandler(handler func(transport.Conn)) {
	this.disconnectHandlers = append(this.disconnectHandlers, handler)
}

func (this *ConnMux) DefaultHandler(handleFunc transport.MessageHandleFunc) {
	this.defaultHandler = handleFunc
}

func (this *ConnMux) MessageHandler(command uint32, handleFunc transport.MessageHandleFunc) {
	this.messageHandler[command] = handleFunc
}

func (this *ConnMux) HandleConnect(conn transport.Conn) {
	for _, handler := range this.connectHandlers {
		handler(conn)
	}
}

func (this *ConnMux) HandleDisconnect(conn transport.Conn) {
	for _, handler := range this.disconnectHandlers {
		handler(conn)
	}
}

func (this *ConnMux) HandleMessage(conn transport.Conn, msg transport.Message) {
	var messageHandler = this.messageHandler[msg.Command()]
	if messageHandler != nil {
		messageHandler.Handle(conn, msg)
	} else {
		if this.defaultHandler != nil {
			this.defaultHandler.Handle(conn, msg)
		}
	}
}
