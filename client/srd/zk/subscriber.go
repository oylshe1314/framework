package zk

import (
	"context"
	"encoding/json"
	"path"
	"strings"

	"github.com/go-zookeeper/zk"
	"github.com/oylshe1314/framework"
	"github.com/oylshe1314/framework/client/srd"
)

type SubscriberClient struct {
	client

	callbacks map[string]srd.SubscribeCallback
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

	return nil
}

func (this *SubscriberClient) readServiceData(conn *zk.Conn, nodesPath string, zkNodes []string) []*srd.ServiceNode {
	var nodes []*srd.ServiceNode
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

		var node = new(srd.ServiceNode)
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

func (this *SubscriberClient) Subscribe(service string, callback srd.SubscribeCallback) {

}
