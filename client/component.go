package client

import (
	"context"

	"github.com/oylshe1314/framework/component"
)

type Component interface {
	component.Component

	Client() Client
}

type clientComponent struct {
	name   string
	client Client
}

func NewClientComponent(name string, client Client) Component {
	return &clientComponent{
		name:   name,
		client: client,
	}
}

func (this *clientComponent) Name() string {
	return this.name
}

func (this *clientComponent) Init(ctx context.Context) error {
	return this.client.Init(ctx)
}

func (this *clientComponent) Start() error {
	var ac, ok = any(this.client).(AsyncClient)
	if !ok {
		return nil
	}
	go func() {
		_ = ac.Work()
	}()
	return nil
}

func (this *clientComponent) Close() error {
	return this.client.Close()
}

func (this *clientComponent) Client() Client {
	return this.client
}
