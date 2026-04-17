package server

import (
	"context"

	"github.com/oylshe1314/framework/component"
)

type Component[T Server] interface {
	component.Component

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

func (this *serverComponent[T]) Init(ctx context.Context) error {
	return this.server.Init(ctx)
}

func (this *serverComponent[T]) Start() error {
	return this.server.Start()
}

func (this *serverComponent[T]) Close() error {
	return this.server.Close()
}

func (this *serverComponent[T]) Name() string {
	return this.name
}

func (this *serverComponent[T]) Server() T {
	return this.server
}
