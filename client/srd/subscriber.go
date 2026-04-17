package srd

import "github.com/oylshe1314/framework/client"

type SubscribeCallback func(service string, nodes []*ServiceNode)

type Subscriber interface {
	client.AsyncClient

	Subscribe(service string, callback SubscribeCallback)
}
