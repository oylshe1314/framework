package heartbeat

import (
	"context"
	"time"

	"github.com/oylshe1314/framework/errors"
	"github.com/oylshe1314/framework/netx/transport"
	"github.com/oylshe1314/framework/option"
)

type HeartbeatClient struct {
	option.Optional[Option]

	closed bool

	ticker *time.Ticker

	conn transport.Conn

	heartbeatBuilder transport.MessageBuilder
	heartbeatHandler transport.MessageHandler
}

func (this *HeartbeatClient) Init(ctx context.Context) error {
	if this.GetOption() == nil {
		return errors.New("option is nil")
	}

	if this.GetOption().Interval == 0 {
		this.GetOption().Interval = time.Second * 90
	}
	return nil
}

func (this *HeartbeatClient) Close() error {
	if !this.closed {
		this.closed = true
		if this.ticker != nil {
			this.ticker.Stop()
		}
	}
	return nil
}

func (this *HeartbeatClient) tick() {
	if this.heartbeatBuilder == nil {
		_ = this.conn.Send(this.GetOption().Command, nil)
	} else {
		_ = this.conn.Send(this.GetOption().Command, this.heartbeatBuilder.Build())
	}
}

func (this *HeartbeatClient) work() error {
	this.closed = false
	this.ticker = time.NewTicker(this.GetOption().Interval)
	go func() {
		for {
			if this.closed {
				return
			}

			<-this.ticker.C
			this.tick()
		}
	}()
	return nil
}

func (this *HeartbeatClient) Work() error {
	return this.work()
}

func (this *HeartbeatClient) HeartbeatBuilder(builder transport.MessageBuilder) {
	this.heartbeatBuilder = builder
}

func (this *HeartbeatClient) HeartbeatHandler(handler transport.MessageHandler) {
	this.heartbeatHandler = handler
}

func (this *HeartbeatClient) HandleConnect() func(conn transport.Conn) {
	return func(conn transport.Conn) {
		this.conn = conn
	}
}

func (this *HeartbeatClient) HandleHeartbeat() (uint32, transport.MessageHandleFunc) {
	return this.GetOption().Command, func(conn transport.Conn, msg transport.Message) {
		if this.heartbeatHandler != nil {
			this.heartbeatHandler.Handle(conn, msg)
		}
	}
}
