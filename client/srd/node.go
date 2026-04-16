package srd

type Node struct {
	Name string `json:"name"`
	Guid string `json:"guid"`

	Extra map[string]any `json:"extra"`
}

type NetNode struct {
	Node `json:",inline"`

	Network string `json:"network"`
	Address string `json:"address"`

	Codec string `json:"codec"`
}

type HttpNode struct {
	NetNode `json:",inline"`

	Secure   bool   `json:"secure"`
	BasePath string `json:"basePath"`
}

type WebsocketNode struct {
	HttpNode `json:",inline"`

	AllowOrigins []string `json:"allowOrigins"`
}
