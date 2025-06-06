package zk

import (
	"context"
	"github.com/go-zookeeper/zk"
	json "github.com/json-iterator/go"
	"github.com/oylshe1314/framework/client/sd"
	"github.com/oylshe1314/framework/errors"
	"strings"
	"time"
)

type subItem struct {
	ctx    context.Context
	cancel context.CancelFunc

	svrName  string
	callback sd.SubscribeCallback
}

func NewSubscribeClient(config *sd.Config) sd.SubscribeClient {
	return &subscribeClient{client: client{config: config}}
}

type subscribeClient struct {
	client

	subItems map[string]*subItem
}

func (this *subscribeClient) AddSubscribe(name string, callback sd.SubscribeCallback) {
	if this.subItems == nil {
		this.subItems = make(map[string]*subItem)
	}
	this.subItems[name] = &subItem{svrName: name, callback: callback}
}

func (this *subscribeClient) Init() error {
	if len(this.subItems) == 0 {
		return errors.Error("please add subscribe service name before init")
	}

	this.client.connectHandler = this.startItemLoop
	return this.client.Init()
}

func (this *subscribeClient) readServiceData(conn *zk.Conn, nodesPath string, zkNodes []string) ([]*sd.ServerNode, error) {
	var svrNodes []*sd.ServerNode
	for _, zkNode := range zkNodes {
		if !strings.HasPrefix(zkNode, "_c_") {
			continue
		}

		data, _, err := conn.Get(nodesPath + "/" + zkNode)
		if err != nil {
			this.logger.Errorf("Get service node data failed, %v, node: %s", err, zkNode)
			continue
		}

		if len(data) == 0 {
			continue
		}

		var svrNode = new(sd.ServerNode)
		err = json.Unmarshal(data, svrNode)
		if err != nil {
			this.logger.Errorf("Unmarshal service node data failed, %v, node: %s, data: %s", err, zkNode, data)
			continue
		}

		if this.server != nil {
			if svrNode.Name == this.server.Name() && svrNode.AppId == this.server.AppId() {
				continue
			}
		}

		svrNodes = append(svrNodes, svrNode)
	}
	return svrNodes, nil
}

func (this *subscribeClient) itemLoop(conn *zk.Conn, item *subItem) {
	var nodesPath = this.rootPath + "/" + item.svrName + serviceNodesPath
	for {
		zkNodes, _, eventChan, err := conn.ChildrenW(nodesPath)
		if err != nil {
			if errors.Is(err, zk.ErrNoNode) {
				this.logger.Warnf("Subscribe service '%s' node was not exists, path: %s", item.svrName, nodesPath)
				time.Sleep(time.Second * 10)
				continue
			}
			this.logger.Error(err, ", path: ", nodesPath)
			return
		}

		ss, err := this.readServiceData(conn, nodesPath, zkNodes)
		if err != nil {
			this.logger.Error(err, ", path: ", nodesPath)
			return
		}

		item.callback(item.svrName, ss)

		select {
		case event, ok := <-eventChan:
			if !ok {
				return
			}
			if event.Err != nil {
				this.logger.Error(event.Err)
				return
			}
		case <-this.ctx.Done():
			if errors.Is(this.ctx.Err(), context.Canceled) {
				return
			}
		}
	}
}

func (this *subscribeClient) startItemLoop(conn *zk.Conn) {
	for _, item := range this.subItems {
		item.ctx, item.cancel = context.WithCancel(this.ctx)
		go this.itemLoop(conn, item)
	}
}
