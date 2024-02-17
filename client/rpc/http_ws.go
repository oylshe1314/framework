package rpc

import (
	"framework/client"
	"framework/client/sd"
)

type WebSocketRpcNode struct {
	*sd.ServiceNode
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
