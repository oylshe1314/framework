package redis

import (
	"context"
	std "github.com/redis/go-redis/v9"
	"strings"
)

type simple struct {
	client std.UniversalClient
}

func OpenRedis(address, username, password string, db int) Redis {
	var addresses = strings.Split(address, ",")
	if len(addresses) > 1 {
		return OpenCluster(addresses, username, password)
	}

	return &simple{client: std.NewClient(&std.Options{Addr: address, Username: username, Password: password, DB: db})}
}

func (this *simple) Close() error {
	return this.client.Close()
}

func (this *simple) Exec(ctx context.Context, cmd string, args ...interface{}) error {
	args = append([]interface{}{cmd}, args...)
	var dr = this.client.Do(ctx, args...)
	return dr.Err()
}

func (this *simple) String(ctx context.Context, cmd string, args ...interface{}) (string, error) {
	args = append([]interface{}{cmd}, args...)
	var c = std.NewStringCmd(ctx, args...)
	var err = this.client.Process(ctx, c)
	if err != nil {
		return "", err
	}
	return c.Result()
}

func (this *simple) Strings(ctx context.Context, cmd string, args ...interface{}) (Strings, error) {
	args = append([]interface{}{cmd}, args...)
	var c = std.NewStringSliceCmd(ctx, args...)
	var err = this.client.Process(ctx, c)
	if err != nil {
		return nil, err
	}
	return c.Result()
}

func (this *simple) StringMap(ctx context.Context, cmd string, args ...interface{}) (StringMap, error) {
	args = append([]interface{}{cmd}, args...)
	var c = std.NewMapStringStringCmd(ctx, args...)
	var err = this.client.Process(ctx, c)
	if err != nil {
		return nil, err
	}
	return c.Result()
}

func (this *simple) Subscribe(ctx context.Context) *SubConn {
	return &SubConn{conn: this.client.Subscribe(ctx)}
}
