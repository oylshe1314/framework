package websocketx

import (
	"context"
	"framework"
	"framework/errors"
	"framework/log"
	"framework/netx/codec"
	"framework/netx/httpx"
	"framework/netx/route"
	"framework/netx/transport"
	"framework/option"
	"framework/store"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type WebsocketServer struct {
	httpx.HttpServer

	option.Optional[Option]

	codec  codec.Codec
	logger log.Logger

	upgrader *websocket.Upgrader

	connStore store.Store[*conn, struct{}]
}

func (this *WebsocketServer) Init(ctx context.Context) error {
	if this.GetOption() == nil {
		return errors.New("'WebsocketServer' option is nil")
	}

	var err error
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
	if len(this.GetOption().AllowOrigins) == 0 {
		return true
	}
	if len(this.GetOption().AllowOrigins) == 1 && this.GetOption().AllowOrigins[0] == "*" {
		return true
	}
	var origin = request.Header.Get("Origin")
	for _, allowOrigin := range this.GetOption().AllowOrigins {
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
