package httpx

import "github.com/oylshe1314/framework/netx"

type Option struct {
	netx.Option

	BasePath string `json:"basePath"`
	HtmlPath string `json:"htmlPath"`

	Tls *TlsOption `json:"tls"`
}

type TlsOption struct {
	CertFile string `json:"certFile"`
	KeyFile  string `json:"keyFile"`
}
