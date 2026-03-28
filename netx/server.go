package netx

import (
	"context"
	"framework"
	"framework/errors"
	"framework/log"
	"framework/netx/codec"
	"framework/netx/heartbeat"
	"framework/netx/route"
	"framework/netx/transport"
	"framework/option"
	"framework/store"
	"net"
	"runtime/debug"
)

type NetServer struct {
	option.Optional[*Option]

	closed bool

	codec  codec.Codec
	logger log.Logger

	address  net.Addr
	listener net.Listener

	mux *route.ConnMux

	connStore store.Store[*conn, struct{}]
}

//func (this *NetServer) SetOption(option *Option) {
//	this.option = option
//}

func (this *NetServer) Init(ctx context.Context) error {
	if this.GetOption() == nil {
		return errors.New("'NetServer' option is nil")
	}

	var err error
	switch this.GetOption().Network {
	case "tcp":
		this.address, err = net.ResolveTCPAddr(this.GetOption().Network, this.GetOption().Address)
	case "udp":
		this.address, err = net.ResolveUDPAddr(this.GetOption().Network, this.GetOption().Address)
	case "unix":
		this.address, err = net.ResolveUnixAddr(this.GetOption().Network, this.GetOption().Address)
	default:
		return errors.Errorf("unknown network '%s'", this.GetOption().Network)
	}
	if err != nil {
		return err
	}

	this.codec, err = codec.NewCodec(this.GetOption().Codec)
	if err != nil {
		return err
	}

	var loggerServer = framework.ServerFromContext[*log.LoggerServer](ctx, "loggerServer")
	if loggerServer != nil {
		this.logger = loggerServer.Logger()
	} else {
		this.logger = log.NewNoneLogger()
	}

	this.mux = route.NewConnMux()

	var heartbeatServer = framework.ServerFromContext[*heartbeat.HeartbeatServer](ctx, "heartbeatServer")
	if heartbeatServer != nil {
		this.mux.MessageHandler(heartbeatServer.HandleHeartbeat())
	}

	this.connStore = store.New[*conn, struct{}]()

	this.closed = false
	return nil
}

func (this *NetServer) Start() error {
	var err error
	this.listener, err = net.Listen(this.address.Network(), this.address.String())
	if err != nil {
		return err
	}

	go func() {
		err = this.serve()
	}()
	return err
}

func (this *NetServer) serve() error {
	defer func() {
		if this.closed {
			return
		}

		var err = recover()
		if err != nil {
			this.logger.Error(err)
			this.logger.Error(string(debug.Stack()))
		}
	}()

	for {
		c, err := this.listener.Accept()
		if err != nil {
			return err
		}

		var cc = newConn(c, this.logger, this.codec)
		this.connStore.Put(cc, struct{}{})
		go func() {
			defer func() {
				this.mux.HandleDisconnect(cc)
				this.connStore.Remove(cc)
			}()

			this.mux.HandleConnect(cc)
			_ = cc.Serve(transport.MessageHandlerFunc(this.mux.HandleMessage))
		}()
	}
}

func (this *NetServer) Close() error {
	this.closed = true
	this.connStore.Foreach(func(cc *conn, _ struct{}) {
		_ = cc.Close()
	})
	if this.listener != nil {
		_ = this.listener.Close()
	}
	return nil
}

func (this *NetServer) ConnectHandler(handler func(transport.Conn)) {
	this.mux.ConnectHandler(handler)
}

func (this *NetServer) DisconnectHandler(handler func(transport.Conn)) {
	this.mux.DisconnectHandler(handler)
}

func (this *NetServer) DefaultHandler(handler transport.MessageHandlerFunc) {
	this.mux.DefaultHandler(handler)
}

func (this *NetServer) MessageHandler(modId uint16, msgId uint16, handler transport.MessageHandlerFunc) {
	this.mux.MessageHandler(modId, msgId, handler)
}
