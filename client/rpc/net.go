package rpc

import (
	"github.com/oylshe1314/framework/client"
	"github.com/oylshe1314/framework/client/sd"
	"github.com/oylshe1314/framework/errors"
	"github.com/oylshe1314/framework/log"
	"github.com/oylshe1314/framework/message"
	"github.com/oylshe1314/framework/net"
	"github.com/oylshe1314/framework/util"
	"sync"
)

type NetRpcNode struct {
	*sd.ServerNode
	*client.NetClient
}

type NetRpcClient struct {
	closed bool

	logger log.Logger
	codec  message.Codec
	locker sync.RWMutex
	nodes  map[string]map[uint32]*NetRpcNode

	workChan chan *NetRpcNode

	connectHandlers    map[string]func(node *sd.ServerNode, conn *net.Conn)
	disconnectHandlers map[string]func(node *sd.ServerNode, conn *net.Conn)
	messageHandlers    map[string]map[uint32]func(node *sd.ServerNode, msg *net.Message)
	defaultHandlers    map[string]func(node *sd.ServerNode, msg *net.Message)
}

func (this *NetRpcClient) SetLogger(logger log.Logger) {
	this.logger = logger
}

func (this *NetRpcClient) SetCodec(codec message.Codec) {
	this.codec = codec
}

func (this *NetRpcClient) Init() error {
	if this.logger == nil {
		this.logger = log.DefaultLogger
	}

	this.nodes = map[string]map[uint32]*NetRpcNode{}

	this.connectHandlers = map[string]func(node *sd.ServerNode, conn *net.Conn){}
	this.disconnectHandlers = map[string]func(node *sd.ServerNode, conn *net.Conn){}
	this.messageHandlers = map[string]map[uint32]func(node *sd.ServerNode, msg *net.Message){}
	this.defaultHandlers = map[string]func(node *sd.ServerNode, msg *net.Message){}

	return nil
}

func (this *NetRpcClient) Close() error {
	this.closed = true

	if this.workChan != nil {
		close(this.workChan)
	}

	for _, nodes := range this.nodes {
		for _, node := range nodes {
			_ = node.Close()
		}
	}
	return nil
}

func (this *NetRpcClient) Work() error {
	this.workChan = make(chan *NetRpcNode, 8)
	for {
		if this.closed {
			return nil
		}

		var node, ok = <-this.workChan
		if !ok {
			return nil
		}

		var gotoWork = false

		var connectHandler = this.connectHandlers[node.Name]
		if connectHandler != nil {
			gotoWork = true
			node.ConnectHandler(func(conn *net.Conn) {
				connectHandler(node.ServerNode, conn)
			})
		}

		var disconnectHandler = this.connectHandlers[node.Name]
		if disconnectHandler != nil {
			gotoWork = true
			node.DisconnectHandler(func(conn *net.Conn) {
				disconnectHandler(node.ServerNode, conn)
			})
		}

		var messageHandlers = this.messageHandlers[node.Name]
		if messageHandlers != nil {
			gotoWork = true
			for id, handler := range messageHandlers {
				modId, msgId := util.Split2uint16(id)
				node.MessageHandler(modId, msgId, func(msg *net.Message) {
					handler(node.ServerNode, msg)
				})
			}
		}

		var defaultHandler = this.defaultHandlers[node.Name]
		if defaultHandler != nil {
			gotoWork = true
			node.DefaultHandler(func(msg *net.Message) {
				defaultHandler(node.ServerNode, msg)
			})
		}

		if gotoWork {
			go func(node *NetRpcNode) {
				var err = node.Work()
				if err != nil {
					return
				}
			}(node)
		}
	}
}

func (this *NetRpcClient) ConnectHandler(service string, handler func(node *sd.ServerNode, conn *net.Conn)) {
	this.connectHandlers[service] = handler
}

func (this *NetRpcClient) DisconnectHandler(service string, handler func(node *sd.ServerNode, conn *net.Conn)) {
	this.disconnectHandlers[service] = handler
}

func (this *NetRpcClient) MessageHandler(service string, modId, msgId uint16, handler func(node *sd.ServerNode, msg *net.Message)) {
	var messageHandlers = this.messageHandlers[service]
	if messageHandlers == nil {
		messageHandlers = map[uint32]func(node *sd.ServerNode, msg *net.Message){}
		this.messageHandlers[service] = messageHandlers
	}
	messageHandlers[util.Compose2uint16(modId, msgId)] = handler
}

func (this *NetRpcClient) DefaultHandler(service string, handler func(node *sd.ServerNode, msg *net.Message)) {
	this.defaultHandlers[service] = handler
}

func (this *NetRpcClient) SubscribeCallback(service string, nodes []*sd.ServerNode) {
	if len(nodes) == 0 {
		this.locker.Lock()
		delete(this.nodes, service)
		this.locker.Unlock()

		this.logger.Warn("The service subscribe callback received an empty nodes list, service: ", service)
	} else {
		this.locker.Lock()
		var oldNodes = this.nodes[service]
		this.locker.Unlock()

		var clients = make(map[uint32]*NetRpcNode)
		for _, node := range nodes {
			if node.Inner == nil {
				this.logger.Warnf("The inner network information of the service node is nil, service: %s, appId: %d", service, node.AppId)
				continue
			}

			if oldNodes != nil {
				oldNode := oldNodes[node.AppId]
				if oldNode != nil && oldNode.Inner.Network == node.Inner.Network && oldNode.Inner.Address == node.Inner.Address {
					clients[node.AppId] = &NetRpcNode{ServerNode: node, NetClient: oldNode.NetClient}
					continue
				}
			}

			var netClient = &client.NetClient{}
			netClient.WithNetwork(node.Inner.Network)
			netClient.WithAddress(node.Inner.Address)
			netClient.SetLogger(this.logger)
			netClient.SetCodec(this.codec)

			var err = netClient.Init()
			if err != nil {
				this.logger.Errorf("Init the service node failed, service: %s, appId: %d, error: %v", service, node.AppId, err)
				continue
			}

			err = netClient.Dial()
			if err != nil {
				this.logger.Errorf("Dial the service node failed, service: %s, appId: %d, error: %v", service, node.AppId, err)
				continue
			}

			this.logger.Infof("Init the service node succeed, service: %s, appId: %d, address: %s", service, node.AppId, node.Inner.Address)

			var rpcNode = &NetRpcNode{ServerNode: node, NetClient: netClient}

			clients[node.AppId] = rpcNode

			if this.workChan != nil {
				this.workChan <- rpcNode
			}
		}

		this.locker.Lock()
		this.nodes[service] = clients
		this.locker.Unlock()

		var removed []*NetRpcNode
		for appId, oldNode := range oldNodes {
			if clients[appId] == nil {
				removed = append(removed, oldNode)
			}
		}

		for _, rn := range removed {
			_ = rn.Close()
			this.logger.Infof("The service node was closed, service: %s, appId: %d, address: %s", service, rn.AppId, rn.Inner.Address)
		}
	}
}

func (this *NetRpcClient) Servers() []string {
	this.locker.RLock()
	defer this.locker.RUnlock()
	return util.MapKeys(this.nodes)
}

func (this *NetRpcClient) Nodes(service string) map[uint32]*NetRpcNode {
	this.locker.RLock()
	defer this.locker.RUnlock()
	return this.nodes[service]
}

func (this *NetRpcClient) Node(service string, appId uint32) *NetRpcNode {
	var nodes = this.Nodes(service)
	if nodes == nil {
		return nil
	}
	return nodes[appId]
}

func (this *NetRpcClient) RandNode(service string) *NetRpcNode {
	var nodes = util.MapValues(this.Nodes(service))
	if len(nodes) == 0 {
		return nil
	}

	return nodes[util.NewRand().IntN(len(nodes))]
}

func (this *NetRpcClient) nodesSend(nodes map[uint32]*NetRpcNode, modId, msgId uint16, v interface{}) (MultiResults[any], error) {
	var result = MultiResults[any]{}
	for _, node := range nodes {
		result[node.AppId] = &MultiResult[any]{Err: node.Send(modId, msgId, v)}
	}
	return result, nil
}

func (this *NetRpcClient) AllSend(service string, modId, msgId uint16, v interface{}) (MultiResults[any], error) {
	var nodes = this.Nodes(service)
	if nodes == nil {
		return nil, errors.Error("the node is unavailable")
	}
	return this.nodesSend(nodes, modId, msgId, v)
}

func (this *NetRpcClient) MultiSend(service string, appIds []uint32, modId, msgId uint16, v interface{}) (MultiResults[any], error) {
	var nodes = this.Nodes(service)
	if nodes == nil {
		return nil, errors.Error("the node is unavailable")
	}

	var selectNodes = map[uint32]*NetRpcNode{}
	for _, appId := range appIds {
		var node = nodes[appId]
		if node != nil {
			selectNodes[appId] = node
		}
	}
	return this.nodesSend(this.Nodes(service), modId, msgId, v)
}

func (this *NetRpcClient) RandSend(service string, modId, msgId uint16, v interface{}) error {
	var node = this.RandNode(service)
	if node == nil {
		return errors.Error("the node is unavailable")
	}

	return node.Send(modId, msgId, v)
}

func (this *NetRpcClient) AppIdSend(service string, appId uint32, modId, msgId uint16, v interface{}) error {
	var node = this.Node(service, appId)
	if node == nil {
		return errors.Error("the node is unavailable")
	}

	return node.Send(modId, msgId, v)
}

func (this *NetRpcClient) nodesRead(nodes map[uint32]*NetRpcNode) MultiResults[*net.Message] {
	var fs []func() error
	var ars = MultiResults[*net.Message]{}
	for _, node := range nodes {
		var curNode = node
		var ar = &MultiResult[*net.Message]{}

		fs = append(fs, func() error {
			ar.Res, ar.Err = curNode.Read()
			return ar.Err
		})
	}

	util.WaitAll(fs...)
	return ars
}

func (this *NetRpcClient) AllRead(service string) (MultiResults[*net.Message], error) {
	var nodes = this.Nodes(service)
	if nodes == nil {
		return nil, errors.Error("the node is unavailable")
	}
	return this.nodesRead(nodes), nil
}

func (this *NetRpcClient) MultiRead(service string, appIds []uint32) (MultiResults[*net.Message], error) {
	var nodes = this.Nodes(service)
	if nodes == nil {
		return nil, errors.Error("the node is unavailable")
	}

	var selectNodes = map[uint32]*NetRpcNode{}
	for _, appId := range appIds {
		var node = nodes[appId]
		if node != nil {
			selectNodes[appId] = node
		}
	}
	return this.nodesRead(nodes), nil
}

func (this *NetRpcClient) anyRead(nodes map[uint32]*NetRpcNode) (*net.Message, error) {
	var fs []func() (*net.Message, error)
	for _, node := range nodes {
		var curNode = node
		fs = append(fs, func() (*net.Message, error) {
			return curNode.Read()
		})
	}

	return util.WaitAnySucceed(fs...)
}

func (this *NetRpcClient) AnyRead(service string) (MultiResults[*net.Message], error) {
	var nodes = this.Nodes(service)
	if nodes == nil {
		return nil, errors.Error("the node is unavailable")
	}
	return this.nodesRead(nodes), nil
}

func (this *NetRpcClient) AnyOfMultiRead(service string, appIds []uint32) (*net.Message, error) {
	var nodes = this.Nodes(service)
	if nodes == nil {
		return nil, errors.Error("the node is unavailable")
	}

	var selectNodes = map[uint32]*NetRpcNode{}
	for _, appId := range appIds {
		var node = nodes[appId]
		if node != nil {
			selectNodes[appId] = node
		}
	}
	return this.anyRead(nodes)
}

func (this *NetRpcClient) RandRead(service string) (*net.Message, error) {
	var node = this.RandNode(service)
	if node == nil {
		return nil, errors.Error("the node is unavailable")
	}

	return node.Read()
}

func (this *NetRpcClient) AppIdRead(service string, appId uint32) (*net.Message, error) {
	var node = this.Node(service, appId)
	if node == nil {
		return nil, errors.Error("the node is unavailable")
	}
	return node.Read()
}
