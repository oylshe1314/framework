package sd

import (
	"github.com/oylshe1314/framework/client"
	"github.com/oylshe1314/framework/server"
)

type SubscribeCallback func(service string, nodes []*ServerNode)

type SubscribeClient interface {
	client.AsyncClient
	SetServer(server server.Server)
	AddSubscribe(name string, callback SubscribeCallback)
}
