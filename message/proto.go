package message

type MsgServerDetectAck struct {
	ProgramHash string  `json:"programHash"`
	DataHash    string  `json:"dataHash"`
	ConfigHash  string  `json:"configHash"`
	Pid         int     `json:"pid"`
	CPU         float64 `json:"cpu"`
	Memory      float64 `json:"memory"`
	Coroutine   int     `json:"coroutine"`
	Info        string  `json:"info"`
}
