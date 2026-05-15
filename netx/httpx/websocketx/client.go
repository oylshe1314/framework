package websocketx

import (
	"context"
	"crypto/tls"
	"net/http"
	"net/url"

	"github.com/gorilla/websocket"
	"github.com/oylshe1314/framework"
	"github.com/oylshe1314/framework/errors"
	"github.com/oylshe1314/framework/log"
	"github.com/oylshe1314/framework/netx/codec"
	"github.com/oylshe1314/framework/netx/route"
	"github.com/oylshe1314/framework/netx/transport"
	"github.com/oylshe1314/framework/option"
)

type WebsocketClient struct {
	option.Optional[DialOption]

	codec  codec.Codec
	logger log.Logger

	url *url.URL

	dialer *websocket.Dialer

	mux *route.ConnMux

	conn *conn
}

func (this *WebsocketClient) Init(ctx context.Context) error {
	if this.GetOption() == nil {
		return errors.New("option is nil")
	}

	var err error
	this.url, err = url.Parse(this.GetOption().Url)
	if err != nil {
		return err
	}

	this.codec, err = codec.NewCodec(this.GetOption().Codec)
	if err != nil {
		return err
	}

	var tlsConfig *tls.Config
	if this.GetOption().Tls != nil {
		var certificates []tls.Certificate
		certificates = make([]tls.Certificate, 1)
		certificates[0], err = tls.LoadX509KeyPair(this.GetOption().Tls.CertFile, this.GetOption().Tls.KeyFile)
		if err != nil {
			return err
		}
		tlsConfig = &tls.Config{
			Certificates: certificates,
		}
		this.url.Scheme = "wss"
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

	this.dialer = &websocket.Dialer{
		TLSClientConfig: tlsConfig,
		Subprotocols:    this.GetOption().Protocols,
	}

	this.mux = route.NewConnMux()
	return nil
}

func (this *WebsocketClient) Close() error {
	if this.conn != nil {
		this.mux.HandleDisconnect(this.conn)
		return this.conn.Close()
	}
	return nil
}

func (this *WebsocketClient) Dial() error {
	var header http.Header
	if this.GetOption().Origin != "" {
		header = http.Header{
			"Origin": []string{this.GetOption().Origin},
		}
	}
	c, _, err := this.dialer.Dial(this.url.String(), header)
	if err != nil {
		return err
	}
	this.conn = newConn(c, this.logger, this.codec)
	this.mux.HandleConnect(this.conn)
	return nil
}

func (this *WebsocketClient) Read() (transport.Message, error) {
	return this.conn.Read()
}

func (this *WebsocketClient) Send(command uint32, v any) error {
	return this.conn.Send(command, v)
}

func (this *WebsocketClient) work() error {
	if this.conn == nil {
		if err := this.Dial(); err != nil {
			return err
		}
	}
	return this.conn.Serve(this.mux.HandleMessage)
}

func (this *WebsocketClient) Work() error {
	return this.work()
}

func (this *WebsocketClient) ConnectHandler(handler func(transport.Conn)) {
	this.mux.ConnectHandler(handler)
}

func (this *WebsocketClient) DisconnectHandler(handler func(transport.Conn)) {
	this.mux.DisconnectHandler(handler)
}

func (this *WebsocketClient) DefaultHandler(handler transport.MessageHandleFunc) {
	this.mux.DefaultHandler(handler)
}

func (this *WebsocketClient) MessageHandler(command uint32, handler transport.MessageHandleFunc) {
	this.mux.MessageHandler(command, handler)
}
