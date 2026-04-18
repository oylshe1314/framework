package zk

import (
	"time"

	"github.com/oylshe1314/framework/component"
)

type Option struct {
	Name string `json:"name"`
	Guid string `json:"guid"`

	Servers  []string      `json:"servers"`
	Timeout  time.Duration `json:"timeout"`
	RootPath string        `json:"rootPath"`

	Components component.ComponentsOption `json:"components"`
}
