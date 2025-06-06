package server

import (
	"github.com/gorilla/websocket"
	. "github.com/oylshe1314/framework/http/ws"
	"net/http"
)

type WebSocketServer struct {
	HttpServer
	ConnMux

	wsu websocket.Upgrader
}

func (this *WebSocketServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	wc, err := this.wsu.Upgrade(w, r, nil)
	if err != nil {
		this.server.Logger().Error("http request upgrade error, ", err)
		return
	}

	this.server.Logger().Debug("receive a websocket upgrade request, address: ", r.RemoteAddr)

	conn := NewConn(wc, this.server.Logger(), &this.ConnMux)
	go func() {
		_ = conn.Serve()
	}()
}

func (this *WebSocketServer) HandleUpgrade(pattern string) {
	this.HttpServer.sm.Handle(pattern, this)
}

func (this *WebSocketServer) errorHandle(w http.ResponseWriter, r *http.Request, status int, reason error) {
	http.Error(w, reason.Error(), status)
}

func (this *WebSocketServer) checkOrigin(r *http.Request) bool {
	return true
}

func (this *WebSocketServer) Init() (err error) {
	this.wsu.Error = this.errorHandle
	this.wsu.CheckOrigin = this.checkOrigin
	return this.HttpServer.Init()
}
