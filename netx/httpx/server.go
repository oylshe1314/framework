package httpx

import (
	"context"
	"framework"
	"framework/errors"
	"framework/log"
	"net"
	"net/http"

	"github.com/gin-gonic/gin"
)

type HttpServer struct {
	gin.IRouter

	option *Option

	logger log.Logger

	address  net.Addr
	listener net.Listener

	ginEngine  *gin.Engine
	httpServer *http.Server
}

func (this *HttpServer) SetOption(option *Option) {
	this.option = option
}

func (this *HttpServer) Name() string {
	return "httpServer"
}

func (this *HttpServer) Init(ctx context.Context) error {
	if this.option == nil {
		return errors.New("'HttpServer' option is nil")
	}

	var err error
	switch this.option.Network {
	case "tcp":
		this.address, err = net.ResolveTCPAddr(this.option.Network, this.option.Address)
	case "udp":
		this.address, err = net.ResolveUDPAddr(this.option.Network, this.option.Address)
	case "unix":
		this.address, err = net.ResolveUnixAddr(this.option.Network, this.option.Address)
	default:
		return errors.Errorf("unknown network '%s'", this.option.Network)
	}
	if err != nil {
		return err
	}

	var loggerServer = framework.ServerFromContext[*log.LoggerServer](ctx, "loggerServer")
	if loggerServer != nil {
		this.logger = loggerServer.Logger()
	} else {
		this.logger = log.NewNoneLogger()
	}

	this.ginEngine = gin.New()
	this.httpServer = &http.Server{
		Handler: this.ginEngine.Handler(),
	}

	if this.option.BasePath == "" {
		this.option.BasePath = "/"
	}

	if this.option.BasePath != "/" {
		this.IRouter = this.ginEngine.Group(this.option.BasePath)
	} else {
		this.IRouter = this.ginEngine
	}

	if this.option.HtmlPath != "" {
		this.ginEngine.LoadHTMLGlob(this.option.HtmlPath)
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
	if this.option.Tls != nil {
		return this.httpServer.ServeTLS(this.listener, this.option.Tls.CertFile, this.option.Tls.KeyFile)
	} else {
		return this.httpServer.Serve(this.listener)
	}
}

func (this *HttpServer) Close() error {
	return this.httpServer.Close()
}
