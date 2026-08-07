package heartbeat

import (
	"context"
	"sync"
	"sync/atomic"
	"time"

	"github.com/oylshe1314/framework/errors"
	"github.com/oylshe1314/framework/netx/transport"
	"github.com/oylshe1314/framework/option"
)

type HeartbeatServer struct {
	option.Optional[Option]

	closed bool

	ticker *time.Ticker

	mutex sync.Mutex
	ticks atomic.Uint64
	slots []map[transport.Conn]struct{}

	heartbeatHandler transport.MessageHandler
}

func (this *HeartbeatServer) Init(ctx context.Context) error {
	if this.GetOption() == nil {
		return errors.New("'HeartbeatServer' option is nil")
	}

	if this.GetOption().Timeout == 0 {
		this.GetOption().Timeout = time.Second * 120
	}

	this.slots = make([]map[transport.Conn]struct{}, this.GetOption().Timeout)
	for i := range this.slots {
		this.slots[i] = make(map[transport.Conn]struct{})
	}
	return nil
}

func (this *HeartbeatServer) tick() {
	var slot = int(this.ticks.Add(1)-1) % int(this.GetOption().Timeout)

	this.mutex.Lock()
	oldSlot := this.slots[slot]
	this.slots[slot] = make(map[transport.Conn]struct{})
	this.mutex.Unlock()

	if len(oldSlot) == 0 {
		return
	}

	for conn := range oldSlot {
		_ = conn.Close()
	}
	clear(oldSlot)
}

func (this *HeartbeatServer) start() error {
	this.closed = false
	this.ticks.Store(0)
	this.ticker = time.NewTicker(time.Second)
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

func (this *HeartbeatServer) Start() error {
	return this.start()
}

func (this *HeartbeatServer) Close() error {
	if !this.closed {
		this.closed = true
		this.ticker.Stop()
	}
	return nil
}

func (this *HeartbeatServer) add(conn transport.Conn) {
	this.remove(conn)
	var slot = (int(this.ticks.Load()) + int(this.GetOption().Timeout) - 1) % int(this.GetOption().Timeout)

	conn.Attribute().Put("slot", slot)

	this.mutex.Lock()
	defer this.mutex.Unlock()
	this.slots[slot][conn] = struct{}{}
}

func (this *HeartbeatServer) remove(conn transport.Conn) {
	slot, ok := conn.Attribute().Get("slot")
	if !ok {
		return
	}

	this.mutex.Lock()
	defer this.mutex.Unlock()
	delete(this.slots[slot.(int)], conn)
}

func (this *HeartbeatServer) HeartbeatHandler(handler transport.MessageHandler) {
	this.heartbeatHandler = handler
}

func (this *HeartbeatServer) HandleHeartbeat() (uint32, transport.MessageHandleFunc) {
	return this.GetOption().Command, func(conn transport.Conn, msg transport.Message) {
		this.add(conn)
		if this.heartbeatHandler != nil {
			this.heartbeatHandler.Handle(conn, msg)
		} else {
			_ = conn.Send(msg.Command(), nil)
		}
	}
}
