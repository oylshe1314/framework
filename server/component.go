package server

type Component[T Server] interface {
	Name() string
	Server() T
}

type ComponentsOption map[string]string

func (cs ComponentsOption) Get(name string) string {
	c, ok := cs[name]
	if !ok {
		return ""
	}
	return c
}

type serverComponent struct {
	name   string
	server Server
}

func NewServerComponent(name string, server Server) Component[Server] {
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
