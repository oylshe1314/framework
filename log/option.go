package log

type Option struct {
	Dir           string `json:"dir"`
	Logger        string `json:"LoggerType"`
	Level         string `json:"level"`
	WithConsole   bool   `json:"console"`
	WithTimestamp bool   `json:"withTimestamp"`
	Timezone      string `json:"timezone"`
	TimeFormat    string `json:"timeFormat"`
	WithCaller    bool   `json:"withCaller"`
}
