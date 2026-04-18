package websocketx

import "github.com/oylshe1314/framework/netx/httpx"

type Option struct {
	httpx.Option

	AllowOrigins []string `json:"allowOrigins"`
}
