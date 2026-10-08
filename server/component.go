package server

import (
	"context"

	"github.com/oylshe1314/framework/component"
)

type Component interface {
	component.Component

	Server() Server
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

func (this *serverComponent) Init(ctx context.Context) error {
	return this.server.Init(ctx)
}

func (this *serverComponent) Start() error {
	return this.server.Start()
}

func (this *serverComponent) Close() error {
	return this.server.Close()
}

func (this *serverComponent) Name() string {
	return this.name
}

func (this *serverComponent) Server() Server {
	return this.server
}
