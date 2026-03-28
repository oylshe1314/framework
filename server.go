package framework

import (
	"context"
)

type Server interface {
	Init(context.Context) error
	Start() error
	Close() error
}
