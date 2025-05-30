package rpc

import (
	"github.com/oylshe1314/framework/client"
	"github.com/oylshe1314/framework/client/sd"
)

type WebSocketRpcNode struct {
	*sd.ServerNode
	*client.WebSocketClient
}

type WebSocketRpcClient struct {
}

func (this *WebSocketRpcClient) Init() error {
	panic("implement me")
}

func (this *WebSocketRpcClient) Close() error {
	panic("implement me")
}

func (this *WebSocketRpcClient) Work() error {
	panic("implement me")
}
