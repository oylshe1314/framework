package srd

import (
	"github.com/oylshe1314/framework/netx"
	"github.com/oylshe1314/framework/netx/httpx"
	"github.com/oylshe1314/framework/netx/httpx/websocketx"
)

type NodeBase struct {
	Type string `json:"type"`
	Name string `json:"name"`
	Guid string `json:"guid"`

	Network string `json:"network"`
	Address string `json:"address"`

	Extra map[string]any `json:"extra"`
}

type NetNode struct {
	Codec string `json:"codec"`
}

func NewNetNode(name, guid string, option *netx.Option) *ServiceNode {
	return &ServiceNode{
		NodeBase: NodeBase{
			Type:    "net",
			Name:    name,
			Guid:    guid,
			Network: option.Network,
			Address: option.Address,
		},
		NetNode: &NetNode{
			Codec: option.Codec,
		},
	}
}

type HttpNode struct {
	BasePath string `json:"basePath"`
	Secure   bool   `json:"secure"`
}

func NewHttpNode(name, guid string, option *httpx.Option) *ServiceNode {
	return &ServiceNode{
		NodeBase: NodeBase{
			Type:    "http",
			Name:    name,
			Guid:    guid,
			Network: option.Network,
			Address: option.Address,
		},
		HttpNode: &HttpNode{
			BasePath: option.BasePath,
			Secure:   option.Tls != nil,
		},
	}
}

type WebsocketNode struct {
	AllowOrigins []string `json:"allowOrigins"`
	Codec        string   `json:"codec"`
}

func NewWebsocketNode(name, guid string, option *websocketx.Option) *ServiceNode {
	return &ServiceNode{
		NodeBase: NodeBase{
			Type:    "websocket",
			Name:    name,
			Guid:    guid,
			Network: option.Network,
			Address: option.Address,
		},
		WebsocketNode: &WebsocketNode{
			AllowOrigins: option.AllowOrigins,
			Codec:        option.Codec,
		},
	}
}

type ServiceNode struct {
	NodeBase       `json:",inline"`
	*NetNode       `json:",inline"`
	*HttpNode      `json:",inline"`
	*WebsocketNode `json:",inline"`
}
