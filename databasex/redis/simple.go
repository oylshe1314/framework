package redis

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type simpleClient struct {
	client redis.UniversalClient
}

type Options redis.Options

func OpenRedis(address string, withOptions ...func(options *Options)) Redis {
	var options = &redis.Options{
		Addr: address,
	}

	for i := range withOptions {
		withOptions[i]((*Options)(options))
	}

	return &simpleClient{
		client: redis.NewClient(options),
	}
}

func (this *simpleClient) Close() error {
	return this.client.Close()
}

func (this *simpleClient) Exec(ctx context.Context, cmd string, args ...interface{}) error {
	args = append([]interface{}{cmd}, args...)
	var dr = this.client.Do(ctx, args...)
	return dr.Err()
}

func (this *simpleClient) String(ctx context.Context, cmd string, args ...interface{}) (string, error) {
	args = append([]interface{}{cmd}, args...)
	var c = redis.NewStringCmd(ctx, args...)
	var err = this.client.Process(ctx, c)
	if err != nil {
		return "", err
	}
	return c.Result()
}

func (this *simpleClient) Strings(ctx context.Context, cmd string, args ...interface{}) (Strings, error) {
	args = append([]interface{}{cmd}, args...)
	var c = redis.NewStringSliceCmd(ctx, args...)
	var err = this.client.Process(ctx, c)
	if err != nil {
		return nil, err
	}
	return c.Result()
}

func (this *simpleClient) StringMap(ctx context.Context, cmd string, args ...interface{}) (StringMap, error) {
	args = append([]interface{}{cmd}, args...)
	var c = redis.NewMapStringStringCmd(ctx, args...)
	var err = this.client.Process(ctx, c)
	if err != nil {
		return nil, err
	}
	return c.Result()
}

func (this *simpleClient) Subscribe(ctx context.Context) SubConn {
	return &subConn{c: this.client.Subscribe(ctx)}
}
