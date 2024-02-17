package sd

import (
	"framework/client"
	"framework/log"
)

type SubscribeCallback func(service string, nodes []*ServiceNode)

type SubscribeClient interface {
	client.AsyncClient
	SetLogger(logger log.Logger)
	AddSubscribe(name string, callback SubscribeCallback)
}
