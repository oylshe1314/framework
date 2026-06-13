package srd

type SubscribeCallback func(service string, nodes []*ServiceNode)

type Subscriber interface {
	client.AsyncClient

	Subscribe(service string, callback SubscribeCallback)
}
