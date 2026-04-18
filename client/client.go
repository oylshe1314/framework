package client

import "context"

type Client interface {
	Init(ctx context.Context) error
	Close() error
}

type AsyncClient interface {
	Client

	Work() error
}
