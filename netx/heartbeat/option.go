package heartbeat

import "time"

type Option struct {
	ModId uint16 `json:"modId"`
	MsgId uint16 `json:"msgId"`

	Timeout time.Duration `json:"timeout"`
}
