package netx

import "github.com/oylshe1314/framework/server"

type Option struct {
	Network string `json:"network"`
	Address string `json:"address"`

	Codec string `json:"codec"`

	Components server.ComponentsOption `json:"components"`
}
