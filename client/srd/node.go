package srd

type Node struct {
	Type string `json:"type"`
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

	BasePath string `json:"basePath"`
	Secure   bool   `json:"secure"`
}

type WebsocketNode struct {
	HttpNode `json:",inline"`

	AllowOrigins []string `json:"allowOrigins"`
}
