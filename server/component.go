package server

type Component interface {
	Name() string
	Server() Server
}

type ComponentsOption struct {
	Logger    string `json:"logger"`
	Heartbeat string `json:"heartbeat"`
}

type serverComponent struct {
	name   string
	server Server
}

func NewServerComponent(name string, server Server) Component {
	return &serverComponent{
		name:   name,
		server: server,
	}
}

func (this *serverComponent) Name() string {
	return this.name
}

func (this *serverComponent) Server() Server {
	return this.server
}
