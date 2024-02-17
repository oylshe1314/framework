package redis

import (
	"context"
	std "github.com/redis/go-redis/v9"
)

type Pong *std.Pong
type Message *std.Message
type Subscription *std.Subscription

type SubConn struct {
	conn *std.PubSub
}

func (this *SubConn) Receive() (interface{}, error) {
	var res, err = this.conn.Receive(context.Background())
	if err != nil {
		return nil, err
	}
	switch rr := res.(type) {
	case *std.Subscription:
		return Subscription(rr), nil
	case *std.Message:
		return Message(rr), nil
	case *std.Pong:
		return Pong(rr), nil
	case error:
		return nil, rr
	default:
		return res, err
	}
}

func (this *SubConn) Close() (err error) {
	return this.conn.Close()
}

func (this *SubConn) Subscribe(channel string) error {
	return this.conn.Subscribe(context.Background(), channel)
}
