package client

import (
	"framework/errors"
	. "framework/http/ws"
	"github.com/gorilla/websocket"
	"net/url"
)

type WebSocketClient struct {
	HttpClient
	ConnMux

	conn *Conn
}

func (this *WebSocketClient) Close() (err error) {
	if this.conn != nil {
		return this.conn.Close()
	}
	return
}

func (this *WebSocketClient) Work() error {
	return this.conn.Serve()
}

func (this *WebSocketClient) Dial(pattern string) error {
	var u, err = url.Parse(this.address)
	if err != nil {
		return err
	}

	u.Path = pattern

	wc, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		return err
	}

	this.conn = NewConn(wc, this.logger, &this.ConnMux)
	return nil
}

func (this *WebSocketClient) Send(modId, msgId uint16, v interface{}) error {
	if this.conn == nil {
		return errors.Error("please connect server first")
	}
	return this.conn.Send(modId, msgId, v)
}

func (this *WebSocketClient) Read() (*Message, error) {
	if this.conn == nil {
		return nil, errors.Error("please connect server first")
	}
	return this.conn.Read()
}
