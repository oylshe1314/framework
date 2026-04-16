package zk

import (
	"context"
	"framework"
	"framework/errors"
	"framework/log"
	"framework/option"
	"time"

	"github.com/go-zookeeper/zk"
)

const (
	defaultRootPath    = "/sk.org/server"
	defaultServicePath = "/services"
	defaultNodesPath   = "/nodes"

	defaultTimeout = time.Millisecond * 30000
)

type client struct {
	option.Optional[Option]

	ctx    context.Context
	cancel context.CancelFunc

	logger log.Logger

	zkDialer zk.Dialer

	connectHandler    func(*zk.Conn)
	disconnectHandler func(*zk.Conn)
}

func (this *client) Init(ctx context.Context) error {
	if this.GetOption() == nil {
		return errors.New("'client' option is nil")
	}

	if len(this.GetOption().Servers) == 0 {
		return errors.New("'client' server list is empty")
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

	var loggerServer = framework.ComponentFromContext[*log.LoggerServer](ctx, loggerServerName).Server()
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
	for {
		conn, eventChan, err = zk.Connect(this.GetOption().Servers, this.GetOption().Timeout, zk.WithDialer(this.zkDialer), zk.WithLogger(&logger{Logger: this.logger}))
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
						conn.Close()
						conn = nil
						if this.disconnectHandler != nil {
							this.disconnectHandler(nil)
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
					if this.disconnectHandler != nil {
						this.disconnectHandler(conn)
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
	if this.connectHandler == nil && this.disconnectHandler == nil {
		return errors.Error("at least one of 'connectHandler' and 'disconnectHandler' is not nil")
	}
	return this.work()
}
