package client

import (
	"context"

	"github.com/oylshe1314/framework/component"
)

type Component[T Client] interface {
	component.Component

	Client() T
}

type clientComponent[T Client] struct {
	name   string
	client T
}

func NewClientComponent[T Client](name string, client T) Component[T] {
	return &clientComponent[T]{
		name:   name,
		client: client,
	}
}

func (this *clientComponent[T]) Name() string {
	return this.name
}

func (this *clientComponent[T]) Init(ctx context.Context) error {
	return this.client.Init(ctx)
}

func (this *clientComponent[T]) Start() error {
	var ac, ok = any(this.client).(AsyncClient)
	if !ok {
		return nil
	}
	go func() {
		_ = ac.Work()
	}()
	return nil
}

func (this *clientComponent[T]) Close() error {
	return this.client.Close()
}

func (this *clientComponent[T]) Client() T {
	return this.client
}
