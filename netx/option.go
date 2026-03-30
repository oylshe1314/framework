package netx

import "framework/server"

type Option struct {
	Network string `json:"network"`
	Address string `json:"address"`

	Codec string `json:"codec"`

	Components server.Components `json:"components"`
}
