package websocketx

import "framework/netx/httpx"

type Option struct {
	httpx.Option

	AllowOrigins []string `json:"allowOrigins"`
}
