package zk

import (
	"context"
	"github.com/go-zookeeper/zk"
	"github.com/oylshe1314/framework/client/sd"
	"github.com/oylshe1314/framework/errors"
	"github.com/oylshe1314/framework/server"
	"time"
)

const DefaultTimeout = time.Millisecond * 30000

const (
	defaultRootPath    = "/sk.org/server"
	serviceServicePath = "/service"
	serviceNodesPath   = "/nodes"
)

type client struct {
	config *sd.Config

	server server.Server

	rootPath string
	timeout  time.Duration

	ctx    context.Context
	cancel context.CancelFunc

	connectHandler func(conn *zk.Conn)
	closeHandler   func(conn *zk.Conn)
}

func (this *client) SetServer(svr server.Server) {
	this.server = svr
}

func (this *client) ConnectHandler(connectedHandler func(conn *zk.Conn)) {
	this.connectHandler = connectedHandler
}

func (this *client) CloseHandler(closeHandler func(conn *zk.Conn)) {
	this.closeHandler = closeHandler
}

func (this *client) Init() error {
	if this.config == nil {
		return errors.Error("Service register-discovery client init 'config' can not be nil")
	}

	if this.server == nil {
		return errors.Error("Service register-discovery client init 'server' can not be nil")
	}

	var ok bool
	this.rootPath, ok = this.config.Extra["rootPath"].(string)
	if !ok || this.rootPath == "" {
		this.rootPath = defaultRootPath
	} else {
		if this.rootPath[len(this.rootPath)-1] == '/' {
			this.rootPath = this.rootPath[:len(this.rootPath)-1]
		}
	}

	if this.config.Timeout == 0 {
		this.timeout = DefaultTimeout
	} else {
		this.timeout = time.Millisecond * this.config.Timeout
	}

	this.ctx, this.cancel = context.WithCancel(context.Background())
	return nil
}

func (this *client) Close() error {
	if this.cancel != nil {
		this.cancel()
	}
	return nil
}

func (this *client) work() error {
	var err error
	var conn *zk.Conn
	var eventChan <-chan zk.Event
	for {
		conn, eventChan, err = zk.Connect(this.config.Servers, this.timeout, zk.WithLogger(this.server.Logger()))
		if err != nil {
			this.server.Logger().Error(err)
			time.Sleep(time.Second * 3)
			continue
		}

	eventLoop:
		for {
			select {
			case event, ok := <-eventChan:
				if !ok {
					break eventLoop
				}
				if event.Err != nil {
					this.server.Logger().Error(err)
				}
				if event.Type != zk.EventSession {
					continue
				}

				switch event.State {
				case zk.StateDisconnected:
					this.server.Logger().Warn("Zookeeper server disconnected, will reconnect after")
					if conn != nil {
						conn.Close()
						conn = nil
						if this.closeHandler != nil {
							this.closeHandler(nil)
						}
					}
					time.Sleep(time.Second * 3)
					break eventLoop
				case zk.StateAuthFailed:
					return errors.Errorf("zookeeper server '%s' authentication failed", conn.Server())
				case zk.StateConnectedReadOnly:
					return errors.Errorf("zookeeper server '%s' is connected but read only", conn.Server())
				case zk.StateHasSession:
					if this.connectHandler != nil {
						this.connectHandler(conn)
					}
					continue
				}
			case <-this.ctx.Done():
				if errors.Is(this.ctx.Err(), context.Canceled) {
					if this.closeHandler != nil {
						this.closeHandler(conn)
					}
					conn.Close()
					conn = nil
					return nil
				}
			}
		}
	}
}

func (this *client) Work() error {
	if this.connectHandler == nil && this.closeHandler == nil {
		return errors.Error("at least one of 'connectedHandler' and 'closeHandler' is not nil")
	}
	return this.work()
}
