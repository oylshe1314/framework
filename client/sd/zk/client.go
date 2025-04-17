package zk

import (
	"context"
	"github.com/go-zookeeper/zk"
	"github.com/oylshe1314/framework/errors"
	"github.com/oylshe1314/framework/log"
	"time"
)

const DefaultTimeout = time.Millisecond * 30000

const (
	defaultRootPath    = "/sk.org/server"
	serviceServicePath = "/service"
	serviceNodesPath   = "/nodes"
)

type Client struct {
	ctx    context.Context
	cancel context.CancelFunc

	servers  []string
	timeout  time.Duration
	rootPath string

	logger log.Logger

	connectHandler func(conn *zk.Conn)
	closeHandler   func(conn *zk.Conn)
}

func (this *Client) WithServers(servers []string) {
	this.servers = servers
}

func (this *Client) WithTimeout(timeout time.Duration) {
	this.timeout = time.Millisecond * timeout
}

func (this *Client) WithRootPath(rootPath string) {
	this.rootPath = rootPath
}

func (this *Client) SetLogger(logger log.Logger) {
	this.logger = logger
}

func (this *Client) ConnectHandler(connectedHandler func(conn *zk.Conn)) {
	this.connectHandler = connectedHandler
}

func (this *Client) CloseHandler(closeHandler func(conn *zk.Conn)) {
	this.closeHandler = closeHandler
}

func (this *Client) Init() error {
	if this.logger == nil {
		this.logger = log.DefaultLogger
	}

	if len(this.servers) == 0 {
		return errors.Error("empty zookeeper server list")
	} else {
		for _, server := range this.servers {
			if len(server) == 0 {
				return errors.Error("wrong address of zookeeper server")
			}
		}
	}

	if this.rootPath == "" {
		this.rootPath = defaultRootPath
	} else {
		if this.rootPath[len(this.rootPath)-1] == '/' {
			this.rootPath = this.rootPath[:len(this.rootPath)-1]
		}
	}

	if this.timeout == 0 {
		this.timeout = DefaultTimeout
	}

	this.ctx, this.cancel = context.WithCancel(context.Background())
	return nil
}

func (this *Client) Close() error {
	if this.cancel != nil {
		this.cancel()
	}
	return nil
}

func (this *Client) work() error {
	var err error
	var conn *zk.Conn
	var eventChan <-chan zk.Event
	for {
		conn, eventChan, err = zk.Connect(this.servers, this.timeout, zk.WithLogger(this.logger))
		if err != nil {
			this.logger.Error(err)
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
					this.logger.Error(err)
				}
				if event.Type != zk.EventSession {
					continue
				}

				switch event.State {
				case zk.StateDisconnected:
					this.logger.Warn("Zookeeper server disconnected, will reconnect after")
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

func (this *Client) Work() error {
	if this.connectHandler == nil && this.closeHandler == nil {
		return errors.Error("at least one of 'connectedHandler' and 'closeHandler' is not nil")
	}
	return this.work()
}
