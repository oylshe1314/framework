package sd

import (
	"github.com/oylshe1314/framework/errors"
	"time"
)

type Config struct {
	Servers []string
	Timeout time.Duration
	Extra   map[string]any
}

func (this *Config) WithServers(servers []string) {
	this.Servers = servers
}

func (this *Config) WithTimeout(timeout int) {
	this.Timeout = time.Duration(timeout)
}

func (this *Config) WithExtra(extra map[string]any) {
	this.Extra = extra
}

func (this *Config) Init() error {
	if len(this.Servers) == 0 {
		return errors.Error("empty register center server list")
	} else {
		for _, server := range this.Servers {
			if len(server) == 0 {
				return errors.Error("empty address of register center server")
			}
		}
	}
	return nil
}
