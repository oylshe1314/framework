package netx

import (
	"context"
	"net"

	"github.com/oylshe1314/framework"
	"github.com/oylshe1314/framework/errors"
	"github.com/oylshe1314/framework/log"
	"github.com/oylshe1314/framework/netx/codec"
	"github.com/oylshe1314/framework/netx/heartbeat"
	"github.com/oylshe1314/framework/netx/route"
	"github.com/oylshe1314/framework/netx/transport"
	"github.com/oylshe1314/framework/option"
)

type NetClient struct {
	option.Optional[Option]

	codec  codec.Codec
	logger log.Logger

	address net.Addr

	connMux *route.ConnMux

	conn *conn
}

func (this *NetClient) Init(ctx context.Context) error {
	if this.GetOption() == nil {
		return errors.New("option is nil")
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

	var loggerServerName = this.GetOption().Components.Get("loggerServer")
	if loggerServerName == "" {
		loggerServerName = "loggerServer"
	}

	var loggerServer = framework.ServerFromContext[*log.LoggerServer](ctx, loggerServerName)
	if loggerServer != nil {
		this.logger = loggerServer.Logger()
	} else {
		this.logger = log.NewNoneLogger()
	}

	this.connMux = route.NewConnMux()

	var heartbeatClientName = this.GetOption().Components.Get("heartbeatClient")
	if heartbeatClientName != "" {
		heartbeatClientName = "heartbeatClient"
	}

	var heartbeatClient = framework.ClientFromContext[*heartbeat.HeartbeatClient](ctx, heartbeatClientName)
	if heartbeatClient != nil {
		this.connMux.ConnectHandler(heartbeatClient.HandleConnect())
		this.connMux.MessageHandler(heartbeatClient.HandleHeartbeat())
	}

	return nil
}

func (this *NetClient) Close() error {
	if this.conn != nil {
		this.connMux.HandleDisconnect(this.conn)
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
	this.connMux.HandleConnect(this.conn)
	return nil
}

func (this *NetClient) Read() (transport.Message, error) {
	return this.conn.Read()
}

func (this *NetClient) Send(command uint32, v any) error {
	return this.conn.Send(command, v)
}

func (this *NetClient) work() error {
	if this.conn == nil {
		if err := this.Dial(); err != nil {
			return err
		}
	}
	return this.conn.Serve(this.connMux.HandleMessage)
}

func (this *NetClient) Work() error {
	return this.work()
}

func (this *NetClient) ConnectHandler(handler func(transport.Conn)) {
	this.connMux.ConnectHandler(handler)
}

func (this *NetClient) DisconnectHandler(handler func(transport.Conn)) {
	this.connMux.DisconnectHandler(handler)
}

func (this *NetClient) DefaultHandler(handler transport.MessageHandleFunc) {
	this.connMux.DefaultHandler(handler)
}

func (this *NetClient) MessageHandler(command uint32, handler transport.MessageHandleFunc) {
	this.connMux.MessageHandler(command, handler)
}
