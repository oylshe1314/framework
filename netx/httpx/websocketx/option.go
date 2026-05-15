package websocketx

import (
	"github.com/oylshe1314/framework/component"
	"github.com/oylshe1314/framework/netx/httpx"
)

type Option struct {
	httpx.Option

	AllowOrigins []string `json:"allowOrigins"`
}

type DialOption struct {
	Url       string   `json:"url"`
	Protocols []string `json:"protocols"`
	Origin    string   `json:"origin"`

	Codec string `json:"codec"`

	Tls *httpx.TlsOption `json:"tls"`

	Components component.ComponentsOption `json:"components"`
}
