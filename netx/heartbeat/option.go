package heartbeat

import "time"

type Option struct {
	Command uint32 `json:"command"`

	Timeout  time.Duration `json:"timeout"`
	Interval time.Duration `json:"interval"`
}
