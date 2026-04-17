package zk

import (
	"context"
	"net"
	"time"

	"github.com/oylshe1314/framework"
	"github.com/oylshe1314/framework/errors"
	"github.com/oylshe1314/framework/log"
	"github.com/oylshe1314/framework/option"

	"github.com/go-zookeeper/zk"
)

const (
	defaultRootPath    = "/sk.org/server"
	defaultServicePath = "/services"
	defaultNodesPath   = "/nodes"

	defaultTimeout = time.Millisecond * 30000
)

type connHandler interface {
	handleConnect(conn *zk.Conn)
	handleDisconnect(conn *zk.Conn)
}

type client struct {
	option.Optional[Option]

	ctx    context.Context
	cancel context.CancelFunc

	logger log.Logger

	dialer zk.Dialer

	handler connHandler
}

func (this *client) Init(ctx context.Context) error {
	if this.GetOption() == nil {
		return errors.New("option is nil")
	}

	if len(this.GetOption().Servers) == 0 {
		return errors.New("option 'servers' is empty")
	}

	if this.GetOption().Timeout == 0 {
		this.GetOption().Timeout = defaultTimeout
	}

	if this.GetOption().RootPath == "" {
		this.GetOption().RootPath = defaultRootPath
	}

	var loggerServerName = this.GetOption().Components.Get("loggerServer")
	if loggerServerName == "" {
		loggerServerName = "loggerServer"
	}

	var loggerServer = framework.ServerFromContext[*log.LoggerServer](ctx, loggerServerName)
	if loggerServer != nil {
		this.logger = loggerServer.Logger()
	} else {
		this.logger = log.NewNoneLogger()
	}

	this.ctx, this.cancel = context.WithCancel(ctx)
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

	var dialer = this.dialer
	if dialer == nil {
		dialer = net.DialTimeout
	}

	var logger = &internalLogger{Logger: this.logger}
	for {
		conn, eventChan, err = zk.Connect(this.GetOption().Servers, this.GetOption().Timeout, zk.WithDialer(dialer), zk.WithLogger(logger))
		if err != nil {
			this.logger.Error("connect to zookeeper server failed, ", err)
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
						this.handler.handleDisconnect(nil)
						conn.Close()
						conn = nil
					}
					time.Sleep(time.Second * 3)
					break eventLoop
				case zk.StateAuthFailed:
					return errors.Errorf("zookeeper server '%s' authentication failed", conn.Server())
				case zk.StateConnectedReadOnly:
					return errors.Errorf("zookeeper server '%s' is connected but read only", conn.Server())
				case zk.StateHasSession:
					this.handler.handleConnect(conn)
					continue
				}
			case <-this.ctx.Done():
				if errors.Is(this.ctx.Err(), context.Canceled) {
					this.handler.handleConnect(conn)
					conn.Close()
					conn = nil
					return nil
				}
			}
		}
	}
}

func (this *client) Work() error {
	if this.handler == nil {
		return errors.Error("the connection handler is nil")
	}
	return this.work()
}
