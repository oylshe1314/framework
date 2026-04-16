package server

type Component[T Server] interface {
	Name() string
	Server() T
}

type serverComponent[T Server] struct {
	name   string
	server T
}

func NewServerComponent[T Server](name string, server T) Component[T] {
	return &serverComponent[T]{
		name:   name,
		server: server,
	}
}

func (this *serverComponent[T]) Name() string {
	return this.name
}

func (this *serverComponent[T]) Server() T {
	return this.server
}

type ComponentsOption map[string]string

func (cso ComponentsOption) Get(name string) (string, bool) {
	c, ok := cso[name]
	return c, ok
}
