package netx

import (
	"github.com/oylshe1314/framework/component"
)

type Option struct {
	Network string `json:"network"`
	Address string `json:"address"`

	Codec string `json:"codec"`

	Components component.ComponentsOption `json:"components"`
}
