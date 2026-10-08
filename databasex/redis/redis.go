package redis

import "context"

type Strings []string

type StringMap map[string]string

type SubConn interface {
	Close() error
	Subscribe(channels ...string) error
	Receive() (any, error)
}

type Redis interface {
	Close() error
	Exec(ctx context.Context, cmd string, args ...any) error
	String(ctx context.Context, cmd string, args ...any) (string, error)
	Strings(ctx context.Context, cmd string, args ...any) (Strings, error)
	StringMap(ctx context.Context, cmd string, args ...any) (StringMap, error)
	Subscribe(ctx context.Context) SubConn
}
