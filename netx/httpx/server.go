package httpx

import (
	"context"
	"framework"
	"framework/errors"
	"framework/log"
	"framework/netx/route"
	"framework/option"
	"net"
	"net/http"

	"github.com/gin-gonic/gin"
)

type HttpServer struct {
	gin.IRouter

	option.Optional[*Option]

	logger log.Logger

	address  net.Addr
	listener net.Listener

	httpMux *route.HttpMux

	httpServer *http.Server
}

func (this *HttpServer) Init(ctx context.Context) error {
	if this.GetOption() == nil {
		return errors.New("'HttpServer' option is nil")
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

	var loggerServerName = this.GetOption().Components.Logger
	if loggerServerName == "" {
		loggerServerName = "loggerServer"
	}

	var loggerServer = framework.ServerFromContext[*log.LoggerServer](ctx, loggerServerName)
	if loggerServer != nil {
		this.logger = loggerServer.Logger()
	} else {
		this.logger = log.NewNoneLogger()
	}

	this.httpMux = route.NewHttpMux()

	if this.GetOption().HtmlPath != "" {
		this.httpMux.LoadHTMLGlob(this.GetOption().HtmlPath)
	}

	this.httpServer = &http.Server{
		Handler: this.httpMux.Handler(),
	}

	return nil
}

func (this *HttpServer) Start() error {
	var err error
	this.listener, err = net.Listen(this.address.Network(), this.address.String())
	if err != nil {
		return err
	}

	go func() {
		_ = this.serve()
	}()
	return nil
}

func (this *HttpServer) serve() error {
	if this.GetOption().Tls != nil {
		return this.httpServer.ServeTLS(this.listener, this.GetOption().Tls.CertFile, this.GetOption().Tls.KeyFile)
	} else {
		return this.httpServer.Serve(this.listener)
	}
}

func (this *HttpServer) Close() error {
	return this.httpServer.Close()
}
