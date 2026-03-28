package netx

import (
	"context"
	"framework"
	"framework/errors"
	"framework/log"
	"framework/netx/codec"
	"framework/netx/route"
	"framework/netx/transport"
	"framework/option"
	"net"
)

type NetClient struct {
	option.Optional[*Option]

	codec  codec.Codec
	logger log.Logger

	address net.Addr

	mux *route.ConnMux

	conn *conn
}

func (this *NetClient) Init(ctx context.Context) error {
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
	return nil
}

func (this *NetClient) Close() error {
	if this.conn != nil {
		return this.conn.Close()
	}
	return nil
}

func (this *NetClient) Dial() error {
	c, err := net.Dial(this.address.Network(), this.address.String())
	if err != nil {
		return err
	}
	this.conn = newConn(c, this.logger, this.codec)
	return nil
}

func (this *NetClient) Read() (transport.Message, error) {
	return this.conn.Read()
}

func (this *NetClient) Write(modId, msgId uint16, v any) error {
	return this.conn.Write(modId, msgId, v)
}

func (this *NetClient) work() error {
	return this.conn.Serve(transport.MessageHandlerFunc(this.mux.HandleMessage))
}

func (this *NetClient) Work() error {
	return this.work()
}

func (this *NetClient) ConnectHandler(handler func(transport.Conn)) {
	this.mux.ConnectHandler(handler)
}

func (this *NetClient) DisconnectHandler(handler func(transport.Conn)) {
	this.mux.DisconnectHandler(handler)
}

func (this *NetClient) DefaultHandler(handler transport.MessageHandlerFunc) {
	this.mux.DefaultHandler(handler)
}

func (this *NetClient) MessageHandler(modId uint16, msgId uint16, handler transport.MessageHandlerFunc) {
	this.mux.MessageHandler(modId, msgId, handler)
}
