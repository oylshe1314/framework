package sd

import (
	"framework/client"
	"framework/log"
	"framework/server"
	"framework/util"
)

type ServiceNetwork struct {
	Network string         `json:"network"`
	Address string         `json:"address"`
	Extra   map[string]any `json:"extra"`
}

type ServiceNode struct {
	Guid  string `json:"guid"`
	Name  string `json:"name"`
	AppId uint32 `json:"appId"`

	Inner *ServiceNetwork `json:"inner"`
	Exter *ServiceNetwork `json:"exter"`
}

func NewServiceNode(name string, appId uint32, inner, exter *server.Listener) *ServiceNode {
	var node = &ServiceNode{Guid: util.UUID(), Name: name, AppId: appId}

	if inner != nil {
		node.Inner = &ServiceNetwork{
			Network: inner.Network(),
			Address: inner.Address(),
			Extra:   inner.Extra(),
		}
	}

	if exter != nil {
		node.Exter = &ServiceNetwork{
			Network: exter.Network(),
			Address: exter.Address(),
			Extra:   exter.Extra(),
		}
	}

	return node
}

type RegisterClient interface {
	client.AsyncClient
	SetLogger(logger log.Logger)
	SetServiceNode(node *ServiceNode)
}
