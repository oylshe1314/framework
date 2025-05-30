package rpc

import (
	"github.com/oylshe1314/framework/client"
	"github.com/oylshe1314/framework/client/sd"
	"github.com/oylshe1314/framework/errors"
	"github.com/oylshe1314/framework/log"
	"github.com/oylshe1314/framework/message"
	"github.com/oylshe1314/framework/util"
	"net/url"
	"sync"
)

type HttpRpcNode struct {
	*sd.ServerNode
	*client.HttpClient
}

type HttpRpcClient struct {
	logger log.Logger
	locker sync.RWMutex
	nodes  map[string]map[uint32]*HttpRpcNode
}

func (this *HttpRpcClient) SetLogger(logger log.Logger) {
	this.logger = logger
}

func (this *HttpRpcClient) Init() error {
	if this.logger == nil {
		this.logger = log.DefaultLogger
	}

	this.nodes = make(map[string]map[uint32]*HttpRpcNode)
	return nil
}

func (this *HttpRpcClient) Close() error {
	for _, nodes := range this.nodes {
		for _, node := range nodes {
			_ = node.Close()
		}
	}
	return nil
}

func (this *HttpRpcClient) SubscribeCallback(service string, nodes []*sd.ServerNode) {
	if len(nodes) == 0 {
		this.locker.Lock()
		delete(this.nodes, service)
		this.locker.Unlock()

		this.logger.Warn("The service subscribe callback received an empty nodes list, service: ", service)
	} else {
		this.locker.Lock()
		var oldNodes = this.nodes[service]
		this.locker.Unlock()

		var clients = make(map[uint32]*HttpRpcNode)
		for _, node := range nodes {
			if node.Inner == nil {
				this.logger.Warnf("The inner network information of the service node is nil, service: %s, appId: %d", service, node.AppId)
				continue
			}

			if oldNodes != nil {
				oldNode := oldNodes[node.AppId]
				if oldNode != nil && oldNode.Inner.Network == node.Inner.Network && oldNode.Inner.Address == node.Inner.Address {
					clients[node.AppId] = &HttpRpcNode{ServerNode: node, HttpClient: oldNode.HttpClient}
					continue
				}
			}

			var httpClient = &client.HttpClient{}
			httpClient.WithNetwork(node.Inner.Network)
			httpClient.WithAddress(node.Inner.Address)
			httpClient.SetLogger(this.logger)

			var err = httpClient.Init()
			if err != nil {
				this.logger.Errorf("Init the service node failed, service: %s, appId: %d, error: %v", service, node.AppId, err)
				continue
			}

			this.logger.Infof("Init the service node succeed, service: %s, appId: %d, address: %s", service, node.AppId, node.Inner.Address)

			clients[node.AppId] = &HttpRpcNode{ServerNode: node, HttpClient: httpClient}
		}

		this.locker.Lock()
		this.nodes[service] = clients
		this.locker.Unlock()

		var removed []*HttpRpcNode
		for appId, on := range oldNodes {
			if clients[appId] == nil {
				removed = append(removed, on)
			}
		}

		for _, rn := range removed {
			_ = rn.Close()
			this.logger.Infof("The service node was closed, service: %s, appId: %d, address: %s", service, rn.AppId, rn.Inner.Address)
		}
	}
}

func (this *HttpRpcClient) Servers() []string {
	this.locker.RLock()
	defer this.locker.RUnlock()
	return util.MapKeys(this.nodes)
}

func (this *HttpRpcClient) Nodes(service string) map[uint32]*HttpRpcNode {
	this.locker.RLock()
	defer this.locker.RUnlock()
	return this.nodes[service]
}

func (this *HttpRpcClient) Node(service string, appId uint32) *HttpRpcNode {
	var nodes = this.Nodes(service)
	if nodes == nil {
		return nil
	}
	return nodes[appId]
}

func (this *HttpRpcClient) RandNode(service string) *HttpRpcNode {
	var nodes = util.MapValues(this.Nodes(service))
	if len(nodes) == 0 {
		return nil
	}

	return nodes[util.NewRand().IntN(len(nodes))]
}

func (this *HttpRpcClient) nodesGet(nodes map[uint32]*HttpRpcNode, path string, query url.Values, res interface{}, headers ...client.HttpHeader) MultiResults[*message.Reply] {
	var fs []func() error
	var ars = MultiResults[*message.Reply]{}
	for _, node := range nodes {
		var curNode = node //Don't delete! Don't delete! Don't delete! Say three times for important things. Otherwise, You guess why I wrote the code like this.
		var ar = &MultiResult[*message.Reply]{}

		ars[node.AppId] = ar
		fs = append(fs, func() error {
			ar.Res, ar.Err = curNode.Get(path, query, util.New(res), headers...)
			return ar.Err
		})
	}

	util.WaitAll(fs...)
	return ars
}

func (this *HttpRpcClient) AllGet(service, path string, query url.Values, res interface{}, headers ...client.HttpHeader) (MultiResults[*message.Reply], error) {
	var nodes = this.Nodes(service)
	if nodes == nil {
		return nil, errors.Error("the node is unavailable")
	}

	return this.nodesGet(nodes, path, query, res, headers...), nil
}

func (this *HttpRpcClient) MultiGet(service string, appIds []uint32, path string, query url.Values, res interface{}, headers ...client.HttpHeader) (MultiResults[*message.Reply], error) {
	var nodes = this.Nodes(service)
	if nodes == nil {
		return nil, errors.Error("the node is unavailable")
	}

	var selectNodes = map[uint32]*HttpRpcNode{}
	for _, appId := range appIds {
		var node = nodes[appId]
		if node != nil {
			selectNodes[appId] = node
		}
	}

	return this.nodesGet(selectNodes, path, query, res, headers...), nil
}

func (this *HttpRpcClient) RandGet(service, path string, query url.Values, res interface{}, headers ...client.HttpHeader) (*message.Reply, error) {
	var node = this.RandNode(service)
	if node == nil {
		return nil, errors.Error("do not have any available node")
	}

	return node.Get(path, query, res, headers...)
}

func (this *HttpRpcClient) AppIdGet(service string, appId uint32, path string, query url.Values, res interface{}, headers ...client.HttpHeader) (*message.Reply, error) {
	var node = this.Node(service, appId)
	if node == nil {
		return nil, errors.Error("the node is unavailable")
	}

	return node.Get(path, query, res, headers...)
}

func (this *HttpRpcClient) nodesPost(nodes map[uint32]*HttpRpcNode, path string, query url.Values, req, res interface{}, headers ...client.HttpHeader) MultiResults[*message.Reply] {
	var fs []func() error
	var ars = MultiResults[*message.Reply]{}
	for id, node := range nodes {
		var curNode = node //Don't delete! Don't delete! Don't delete! Say three times for important things. Otherwise, You guess why I wrote the code like this.
		var ar = &MultiResult[*message.Reply]{}

		ars[id] = ar
		fs = append(fs, func() error {
			ar.Res, ar.Err = curNode.Post(path, query, req, util.New(res), headers...)
			return ar.Err
		})
	}

	util.WaitAll(fs...)
	return ars
}

func (this *HttpRpcClient) AllPost(service, path string, query url.Values, req, res interface{}, headers ...client.HttpHeader) (MultiResults[*message.Reply], error) {
	var nodes = this.Nodes(service)
	if nodes == nil {
		return nil, errors.Error("the node is unavailable")
	}

	return this.nodesPost(nodes, path, query, req, res, headers...), nil
}

func (this *HttpRpcClient) MultiPost(service string, appIds []uint32, path string, query url.Values, req, res interface{}, headers ...client.HttpHeader) (MultiResults[*message.Reply], error) {
	var nodes = this.Nodes(service)
	if nodes == nil {
		return nil, errors.Error("the node is unavailable")
	}

	var selectNodes = map[uint32]*HttpRpcNode{}
	for _, appId := range appIds {
		var node = nodes[appId]
		if node != nil {
			selectNodes[appId] = node
		}
	}

	return this.nodesPost(selectNodes, path, query, req, res, headers...), nil
}

func (this *HttpRpcClient) RandPost(service, path string, query url.Values, req, res interface{}, headers ...client.HttpHeader) (*message.Reply, error) {
	var node = this.RandNode(service)
	if node == nil {
		return nil, errors.Error("do not have any available node")
	}

	return node.Post(path, query, req, res, headers...)
}

func (this *HttpRpcClient) AppIdPost(service string, appId uint32, path string, query url.Values, req, res interface{}, headers ...client.HttpHeader) (*message.Reply, error) {
	var node = this.Node(service, appId)
	if node == nil {
		return nil, errors.Error("the node is unavailable")
	}

	return node.Post(path, query, req, res, headers...)
}
