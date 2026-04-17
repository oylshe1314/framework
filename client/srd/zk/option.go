package zk

import (
	"github.com/oylshe1314/framework/server"
	"time"
)

type Option struct {
	Name string
	Guid string

	Servers  []string      `json:"servers"`
	Timeout  time.Duration `json:"timeout"`
	RootPath string        `json:"rootPath"`

	Components server.ComponentsOption `json:"components"`
}
