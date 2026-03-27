package websocketx

import (
	"context"
	"framework"
	"framework/log"
	"framework/netx/codec"
	"framework/netx/httpx"
	"framework/netx/route"
	"framework/netx/transport"
	"framework/store"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type WebsocketServer struct {
	httpx.HttpServer

	option *Option

	codec  codec.Codec
	logger log.Logger

	upgrader *websocket.Upgrader

	connStore store.Store[*conn, struct{}]
}

func (this *WebsocketServer) SetOption(option *Option) {
	this.HttpServer.SetOption(&option.Option)
	this.option = option
}

func (this *WebsocketServer) Init(ctx context.Context) error {
	var err = this.HttpServer.Init(ctx)
	if err != nil {
		return err
	}

	this.codec, err = codec.NewCodec(this.option.Codec)
	if err != nil {
		return err
	}

	var loggerServer = framework.ServerFromContext[*log.LoggerServer](ctx, "loggerServer")
	if loggerServer != nil {
		this.logger = loggerServer.Logger()
	} else {
		this.logger = log.NewNoneLogger()
	}

	this.upgrader = &websocket.Upgrader{
		Error:       this.handleError,
		CheckOrigin: this.checkOrigin,
	}

	this.connStore = store.New[*conn, struct{}]()
	return err
}

func (this *WebsocketServer) Close() error {
	this.connStore.Foreach(func(cc *conn, _ struct{}) {
		_ = cc.Close()
	})
	return this.HttpServer.Close()
}

func (this *WebsocketServer) handleError(w http.ResponseWriter, r *http.Request, status int, reason error) {
	http.Error(w, reason.Error(), status)
}

func (this *WebsocketServer) checkOrigin(request *http.Request) bool {
	if len(this.option.AllowOrigins) == 0 {
		return true
	}
	if len(this.option.AllowOrigins) == 1 && this.option.AllowOrigins[0] == "*" {
		return true
	}
	var origin = request.Header.Get("Origin")
	for _, allowOrigin := range this.option.AllowOrigins {
		if origin == allowOrigin {
			return true
		}
	}
	return false
}

func (this *WebsocketServer) upgradeHandlerFunc(handler route.ConnHandler) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		c, err := this.upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
		if err != nil {
			this.logger.Error(err)
			return
		}

		var cc = newConn(c, this.logger, this.codec)
		this.connStore.Put(cc, struct{}{})
		go func() {
			defer func() {
				handler.HandleDisconnect(cc)
				this.connStore.Remove(cc)
			}()

			handler.HandleConnect(cc)
			_ = cc.Serve(transport.MessageHandlerFunc(handler.HandleMessage))
		}()
	}
}

func (this *WebsocketServer) HandleUpgrade(pattern string, handler route.ConnHandler) {
	this.Any(pattern, this.upgradeHandlerFunc(handler))
}
