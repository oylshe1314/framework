package redis

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type simpleClient struct {
	c redis.UniversalClient
}

func OpenRedis(address string, password string, db int) Redis {
	return nil
}

func (this *simpleClient) Close() error {
	return this.c.Close()
}

func (this *simpleClient) Exec(ctx context.Context, cmd string, args ...interface{}) error {
	args = append([]interface{}{cmd}, args...)
	var dr = this.c.Do(ctx, args...)
	return dr.Err()
}

func (this *simpleClient) String(ctx context.Context, cmd string, args ...interface{}) (string, error) {
	args = append([]interface{}{cmd}, args...)
	var c = redis.NewStringCmd(ctx, args...)
	var err = this.c.Process(ctx, c)
	if err != nil {
		return "", err
	}
	return c.Result()
}

func (this *simpleClient) Strings(ctx context.Context, cmd string, args ...interface{}) (Strings, error) {
	args = append([]interface{}{cmd}, args...)
	var c = redis.NewStringSliceCmd(ctx, args...)
	var err = this.c.Process(ctx, c)
	if err != nil {
		return nil, err
	}
	return c.Result()
}

func (this *simpleClient) StringMap(ctx context.Context, cmd string, args ...interface{}) (StringMap, error) {
	args = append([]interface{}{cmd}, args...)
	var c = redis.NewMapStringStringCmd(ctx, args...)
	var err = this.c.Process(ctx, c)
	if err != nil {
		return nil, err
	}
	return c.Result()
}

func (this *simpleClient) Subscribe(ctx context.Context) SubConn {
	return &subConn{c: this.c.Subscribe(ctx)}
}
