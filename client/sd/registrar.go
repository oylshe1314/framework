package sd

import (
	"github.com/oylshe1314/framework/client"
	"github.com/oylshe1314/framework/server"
	"github.com/oylshe1314/framework/util"
)

type ServerNetwork struct {
	Network string         `json:"network"`
	Address string         `json:"address"`
	Extra   map[string]any `json:"extra"`
}

type ServerNode struct {
	Guid  string `json:"guid"`
	Name  string `json:"name"`
	AppId uint32 `json:"appId"`

	Inner *ServerNetwork `json:"inner"`
	Exter *ServerNetwork `json:"exter"`
}

func NewServiceNode(name string, appId uint32, inner, exter *server.Listener) *ServerNode {
	var node = &ServerNode{Guid: util.UUID(), Name: name, AppId: appId}

	if inner != nil {
		node.Inner = &ServerNetwork{
			Network: inner.Network(),
			Address: inner.Address(),
			Extra:   inner.Extra(),
		}
	}

	if exter != nil {
		node.Exter = &ServerNetwork{
			Network: exter.Network(),
			Address: exter.Address(),
			Extra:   exter.Extra(),
		}
	}

	return node
}

type RegisterClient interface {
	client.AsyncClient
	SetServer(server server.Server)
	SetListener(inner, exter *server.Listener)
}
