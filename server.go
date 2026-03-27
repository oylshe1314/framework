package framework

import (
	"context"
	"framework/errors"
)

type Server interface {
	Init(context.Context) error
	Start() error
	Close() error
}

type NamedServer struct {
	name string
}

func (this *NamedServer) SetName(name string) {
	this.name = name
}

func (this *NamedServer) Init(ctx context.Context) error {
	if this.name == "" {
		return errors.New("name is empty")
	}
	return nil
}

func (this *NamedServer) Start() error {
	return nil
}

func (this *NamedServer) Close() error {
	return nil
}

func (this *NamedServer) Name() string {
	return this.name
}
