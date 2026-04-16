package zk

import (
	"framework/server"
	"time"
)

type Option struct {
	Servers  []string      `json:"servers"`
	Timeout  time.Duration `json:"timeout"`
	RootPath string        `json:"rootPath"`

	Components server.ComponentsOption `json:"components"`
}
