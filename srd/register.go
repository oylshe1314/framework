package srd

import "github.com/oylshe1314/framework/client"

type Register interface {
	client.AsyncClient

	SetNode(node *ServiceNode)
}
