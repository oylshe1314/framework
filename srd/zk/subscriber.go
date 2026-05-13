package zk

import (
	"context"
	"encoding/json"
	"path"
	"strings"
	"time"

	"github.com/go-zookeeper/zk"
	"github.com/oylshe1314/framework"
	"github.com/oylshe1314/framework/errors"
	srd2 "github.com/oylshe1314/framework/srd"
)

type SubscriberClient struct {
	client

	subscriptions map[string]srd2.SubscribeCallback
}

func (this *SubscriberClient) Init(ctx context.Context) error {
	var err = this.client.Init(ctx)
	if err != nil {
		return err
	}

	if this.GetOption().Name == "" {
		var registerClientName = this.GetOption().Components.Get("registerClient")
		if registerClientName == "" {
			registerClientName = "registerClient"
		}

		var registerClient = framework.ClientFromContext[*RegisterClient](ctx, registerClientName)
		if registerClient != nil {
			this.GetOption().Name = registerClient.GetOption().Name
			this.GetOption().Guid = registerClient.GetOption().Guid
		}
	}

	this.client.dialer = nil
	this.client.handler = this
	return nil
}

func (this *SubscriberClient) readServiceData(conn *zk.Conn, nodesPath string, zkNodes []string) []*srd2.ServiceNode {
	var nodes []*srd2.ServiceNode
	for _, zkNode := range zkNodes {
		if !strings.HasPrefix(zkNode, "_c_") {
			continue
		}

		var dataPath = path.Join(nodesPath, zkNode)

		data, _, err := conn.Get(dataPath)
		if err != nil {
			this.logger.Errorf("Get service node data error, path: %s, err: %v", dataPath, err)
			continue
		}

		if len(data) == 0 {
			continue
		}

		var node = new(srd2.ServiceNode)
		err = json.Unmarshal(data, node)
		if err != nil {
			this.logger.Errorf("Unmarshal service node data error, node: %s, data: %s, err: %v", dataPath, data, err)
			continue
		}

		if this.GetOption().Name != "" && node.Name == this.GetOption().Name && node.Guid == this.GetOption().Guid {
			continue
		}

		nodes = append(nodes, node)
	}
	return nodes
}

func (this *SubscriberClient) listenServiceData(conn *zk.Conn, service string, callback srd2.SubscribeCallback) {
	var nodesPath = this.GetOption().RootPath + "/" + service + defaultNodesPath
	for {
		zkNodes, _, eventChan, err := conn.ChildrenW(nodesPath)
		if err != nil {
			if errors.Is(err, zk.ErrNoNode) {
				this.logger.Warnf("The node of service '%s' was not exists, path: %s", service, nodesPath)
				time.Sleep(time.Second * 10)
				continue
			}

			this.logger.Errorf("Get service '%s' child nodes error, path: %s, err: %v", service, nodesPath, err)
			return
		}

		callback(service, this.readServiceData(conn, nodesPath, zkNodes))

		select {
		case event, ok := <-eventChan:
			if !ok {
				continue
			}
			if event.Err != nil {
				if errors.Is(event.Err, zk.ErrConnectionClosed) {
					return
				}
				continue
			}
		case <-this.ctx.Done():
			if errors.Is(this.ctx.Err(), context.Canceled) {
				return
			}
		}
	}
}

func (this *SubscriberClient) startServiceListen(conn *zk.Conn) {
	for service, callback := range this.subscriptions {
		go func() {
			this.listenServiceData(conn, service, callback)
		}()
	}
}

func (this *SubscriberClient) handleConnect(conn *zk.Conn) {
	if conn != nil && conn.State() >= zk.StateConnected {
		this.startServiceListen(conn)
	}
}

func (this *SubscriberClient) handleDisconnect(conn *zk.Conn) {

}

func (this *SubscriberClient) Subscribe(service string, callback srd2.SubscribeCallback) {

}
