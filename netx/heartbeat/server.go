package heartbeat

import (
	"context"
	"framework/errors"
	"framework/netx/transport"
	"sync"
	"sync/atomic"
	"time"
)

type BeatServer struct {
	option *Option

	closed bool

	ticker *time.Ticker

	mutex sync.Mutex
	ticks atomic.Uint64
	slots []map[transport.Conn]struct{}

	heartbeatHandler transport.MessageHandler
}

func (server *BeatServer) SetOption(option *Option) {
	server.option = option
}

func (this *BeatServer) Init(ctx context.Context) error {
	if this.option == nil {
		return errors.New("'BeatServer' option is nil")
	}

	if this.option.Timeout == 0 {
		this.option.Timeout = time.Second * 120
	}

	this.slots = make([]map[transport.Conn]struct{}, this.option.Timeout)
	for i := range this.slots {
		this.slots[i] = make(map[transport.Conn]struct{})
	}
	return nil
}

func (this *BeatServer) tick() {
	var slot = int(this.ticks.Add(1)-1) % int(this.option.Timeout)

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

func (this *BeatServer) start() error {
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

func (this *BeatServer) Start() error {
	return this.start()
}

func (this *BeatServer) Close() error {
	if !this.closed {
		this.closed = true
		this.ticker.Stop()
	}
	return nil
}

func (this *BeatServer) Add(conn transport.Conn) {
	this.Remove(conn)
	var slot = (int(this.ticks.Load()) + int(this.option.Timeout) - 1) % int(this.option.Timeout)

	conn.Attribute().Put("slot", slot)

	this.mutex.Lock()
	defer this.mutex.Unlock()
	this.slots[slot][conn] = struct{}{}
}

func (this *BeatServer) Remove(conn transport.Conn) {
	slot, ok := conn.Attribute().Get("slot")
	if !ok {
		return
	}

	this.mutex.Lock()
	defer this.mutex.Unlock()
	delete(this.slots[slot.(int)], conn)
}

func (this *BeatServer) HeartbeatHandler(handler transport.MessageHandlerFunc) {
	this.heartbeatHandler = handler
}

func (this *BeatServer) HandleHeartbeat() (uint16, uint16, transport.MessageHandlerFunc) {
	return this.option.ModId, this.option.MsgId, func(conn transport.Conn, msg transport.Message) {
		this.Add(conn)
		if this.heartbeatHandler != nil {
			this.heartbeatHandler.Handle(conn, msg)
		} else {
			_ = msg.Reply(nil)
		}
	}
}
