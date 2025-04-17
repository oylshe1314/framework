package sd

import (
	"github.com/oylshe1314/framework/client"
	"github.com/oylshe1314/framework/log"
)

type SubscribeCallback func(service string, nodes []*ServiceNode)

type SubscribeClient interface {
	client.AsyncClient
	SetLogger(logger log.Logger)
	AddSubscribe(name string, callback SubscribeCallback)
}
