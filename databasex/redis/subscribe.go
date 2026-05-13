package redis

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type Pong redis.Pong
type Message redis.Message
type Subscription redis.Subscription

type subConn struct {
	c *redis.PubSub
}

func (this *subConn) Close() (err error) {
	return this.c.Close()
}

func (this *subConn) Subscribe(channel ...string) error {
	return this.c.Subscribe(context.Background(), channel...)
}

func (this *subConn) Receive() (any, error) {
	var res, err = this.c.Receive(context.Background())
	if err != nil {
		return nil, err
	}
	switch rr := res.(type) {
	case *redis.Subscription:
		return (*Subscription)(rr), nil
	case *redis.Message:
		return (*Message)(rr), nil
	case *redis.Pong:
		return (*Pong)(rr), nil
	case error:
		return nil, rr
	default:
		return res, err
	}
}
