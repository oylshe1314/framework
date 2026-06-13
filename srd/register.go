package srd

type Register interface {
	client.AsyncClient

	SetNode(node *ServiceNode)
}
