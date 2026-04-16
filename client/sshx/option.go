package sshx

type Option struct {
	Address  string `json:"address"`
	User     string `json:"user"`
	Password string `json:"password"`
	Key      string `json:"key"`
	KeyPath  string `json:"keyPath"`
}
