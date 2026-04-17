package srd

type SubscribeCallback func(service string, nodes []*ServiceNode)

type Subscriber interface {
	Subscribe(service string, callback SubscribeCallback)
}
