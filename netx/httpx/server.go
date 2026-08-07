package httpx

import (
	"context"
	"net"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/oylshe1314/framework"
	"github.com/oylshe1314/framework/errors"
	"github.com/oylshe1314/framework/log"
	"github.com/oylshe1314/framework/netx/route"
	"github.com/oylshe1314/framework/option"
)

type HttpServer struct {
	option.Optional[Option]

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
		return errors.Errorf("'HttpServer' unknown network '%s'", this.GetOption().Network)
	}
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

	this.httpMux = route.NewHttpMux(this.GetOption().BasePath, this.GetOption().HtmlPath)

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

func (this *HttpServer) Use(middleware ...gin.HandlerFunc) {
	this.httpMux.Use(middleware...)
}

func (this *HttpServer) Handle(pattern string, handler gin.HandlerFunc) {
	this.httpMux.Any(pattern, handler)
}

func (this *HttpServer) HandleGet(pattern string, handler gin.HandlerFunc) {
	this.httpMux.Handle("GET", pattern, handler)
}

func (this *HttpServer) HandlePost(pattern string, handler gin.HandlerFunc) {
	this.httpMux.Handle("POST", pattern, handler)
}

func (this *HttpServer) HandlePut(pattern string, handler gin.HandlerFunc) {
	this.httpMux.Handle("PUT", pattern, handler)
}

func (this *HttpServer) HandleDelete(pattern string, handler gin.HandlerFunc) {
	this.httpMux.Handle("DELETE", pattern, handler)
}

func (this *HttpServer) HandlePatch(pattern string, handler gin.HandlerFunc) {
	this.httpMux.Handle("PATCH", pattern, handler)
}

func (this *HttpServer) HandleOptions(pattern string, handler gin.HandlerFunc) {
	this.httpMux.Handle("OPTIONS", pattern, handler)
}

func (this *HttpServer) HandleHead(pattern string, handler gin.HandlerFunc) {
	this.httpMux.Handle("HEAD", pattern, handler)
}
