package zk

import (
	"context"
	"encoding/json"
	"github.com/oylshe1314/framework"
	"github.com/oylshe1314/framework/client/srd"
	"github.com/oylshe1314/framework/errors"
	"github.com/oylshe1314/framework/netx"
	"github.com/oylshe1314/framework/netx/httpx"
	"github.com/oylshe1314/framework/netx/httpx/websocketx"
	"github.com/oylshe1314/framework/server"
	"net"
	"strings"
	"time"

	"github.com/go-zookeeper/zk"
	"github.com/google/uuid"
)

type RegisterClient struct {
	client

	node any

	version int32
	svcPath string
}

func (this *RegisterClient) extractNode(ctx context.Context) error {
	var ss []server.Server
	var ns = framework.ServersFromContext[*netx.NetServer](ctx)
	for i := range ns {
		ss = append(ss, ns[i])
	}

	var hs = framework.ServersFromContext[*httpx.HttpServer](ctx)
	for i := range hs {
		ss = append(ss, hs[i])
	}

	var ws = framework.ServersFromContext[*websocketx.WebsocketServer](ctx)
	for i := range ws {
		ss = append(ss, ws[i])
	}

	if len(ss) == 0 {
		return errors.New("can not find any server of 'NetServer' or 'HttpServer' or 'WebsocketServer' to register")
	}

	if len(ss) > 1 {
		return errors.New("find multiple server of 'NetServer' or 'HttpServer' or 'WebsocketServer' to register")
	}

	return this.createNode(ss[0])
}

func (this *RegisterClient) createNode(svr server.Server) error {

	var guid = this.GetOption().Guid
	if guid == "" {
		id, err := uuid.NewV7()
		if err != nil {
			return err
		}
		guid = id.String()
	}

	switch s := svr.(type) {
	case *netx.NetServer:
		var netOption = s.GetOption()
		var node = &srd.NetNode{
			Node: srd.Node{
				Type: "net",
				Name: this.GetOption().Name,
				Guid: guid,
			},
			Network: netOption.Network,
			Address: netOption.Address,
			Codec:   netOption.Codec,
		}
		this.node = node
	case *httpx.HttpServer:
		var httpOption = s.GetOption()
		var node = &srd.HttpNode{
			NetNode: srd.NetNode{
				Node: srd.Node{
					Type: "http",
					Name: this.GetOption().Name,
					Guid: guid,
				},
				Network: httpOption.Network,
				Address: httpOption.Address,
				Codec:   httpOption.Codec,
			},
			BasePath: httpOption.BasePath,
			Secure:   httpOption.Tls != nil,
		}
		this.node = node
	case *websocketx.WebsocketServer:
		var websocketOption = s.GetOption()
		var node = &srd.WebsocketNode{
			HttpNode: srd.HttpNode{
				NetNode: srd.NetNode{
					Node: srd.Node{
						Type: "websocket",
						Name: this.GetOption().Name,
						Guid: guid,
					},
					Network: websocketOption.Network,
					Address: websocketOption.Address,
					Codec:   websocketOption.Codec,
				},
				BasePath: websocketOption.BasePath,
				Secure:   websocketOption.Tls != nil,
			},
			AllowOrigins: websocketOption.AllowOrigins,
		}
		this.node = node
	}
	return nil
}

func (this *RegisterClient) Init(ctx context.Context) error {
	var err = this.client.Init(ctx)
	if err != nil {
		return err
	}

	var registerServerName = this.GetOption().Components.Get("registerServer")
	if registerServerName == "" {
		registerServerName = "registerServer"
	}

	var registerServer = framework.ServerFromContext[server.Server](ctx, registerServerName)
	if registerServer != nil {
		err = this.createNode(registerServer)
	} else {
		err = this.extractNode(ctx)
		if err != nil {
			return err
		}
	}

	this.client.zkDialer = this.dial
	this.client.connectHandler = this.handleConnect
	this.client.disconnectHandler = this.handleDisconnect
	return nil
}

func (this *RegisterClient) dial(network, address string, timeout time.Duration) (net.Conn, error) {
	conn, err := net.DialTimeout(network, address, timeout)
	if err != nil {
		return nil, err
	}

	this.node.(*srd.NetNode).Address = conn.LocalAddr().String()
	return conn, nil
}

func (this *RegisterClient) createParentNodes(conn *zk.Conn, path string) error {
	var strPath string
	var nodeNames = strings.Split(path, "/")

	var i = 0
	if len(nodeNames[i]) == 0 {
		i += 1
	}

	for ; i < len(nodeNames); i++ {
		strPath += "/" + nodeNames[i]

		_, err := conn.Create(strPath, []byte{}, 0, zk.WorldACL(zk.PermAll))

		if err != nil && !errors.Is(err, zk.ErrNodeExists) && !errors.Is(err, zk.ErrNoAuth) {
			return err
		}
	}
	return nil
}

func (this *RegisterClient) setServiceNode(conn *zk.Conn) (string, error) {
	var node = this.node.(*srd.Node)

	this.version = 0
	this.svcPath = ""

	data, err := json.Marshal(node)
	if err != nil {
		return "", err
	}

	var servicePath = this.GetOption().RootPath + "/" + node.Name + defaultServicePath
	var nodesPath = this.GetOption().RootPath + "/" + node.Name + defaultNodesPath
	var guidPath = servicePath + "/" + node.Guid

	var version int32
	var nodeExisted = false
	bs, stat, err := conn.Get(guidPath)
	if err != nil {
		if !errors.Is(err, zk.ErrNoNode) {
			return "", err
		}
	} else {
		_, _, err = conn.Get(string(bs))
		if err == nil {
			return "", errors.Errorf("service '%s:%s' is already existed", node.Name, node.Guid)
		}
		if !errors.Is(err, zk.ErrNoNode) {
			return "", err
		}
		version = stat.Version
		nodeExisted = true
	}

	if nodeExisted {
		var path string
		path, err = conn.CreateProtectedEphemeralSequential(nodesPath+"/", data, zk.WorldACL(zk.PermAll))
		if err != nil {
			return "", err
		}

		stat, err = conn.Set(guidPath, []byte(path), version)
		if err != nil {
			return "", err
		}

		this.svcPath = guidPath
		this.version = stat.Version
		return this.svcPath, nil
	} else {
		err = this.createParentNodes(conn, nodesPath)
		if err != nil {
			return "", err
		}

		var path string
		path, err = conn.CreateProtectedEphemeralSequential(nodesPath+"/", data, zk.WorldACL(zk.PermAll))
		if err != nil {
			return "", err
		}

		err = this.createParentNodes(conn, servicePath)
		if err != nil {
			return "", err
		}

		path, err = conn.Create(guidPath, []byte(path), 0, zk.WorldACL(zk.PermAll))
		if err != nil {
			return "", err
		}

		this.svcPath = path
		return this.svcPath, nil
	}
}

func (this *RegisterClient) register(conn *zk.Conn) {
	for {
		path, err := this.setServiceNode(conn)
		if err == nil {
			this.logger.Infof("Service register success, node: %s", path)
			break
		}
		this.logger.Error(err)
		time.Sleep(time.Second * 3)
	}
}

func (this *RegisterClient) handleConnect(conn *zk.Conn) {
	if conn != nil && conn.State() >= zk.StateConnected {
		this.register(conn)
	}
}

func (this *RegisterClient) deleteServiceNode(conn *zk.Conn) {
	if len(this.svcPath) > 0 {
		_ = conn.Delete(this.svcPath, this.version)
	}
}

func (this *RegisterClient) deregister(conn *zk.Conn) {
	this.deleteServiceNode(conn)
}

func (this *RegisterClient) handleDisconnect(conn *zk.Conn) {
	if conn != nil && conn.State() >= zk.StateConnected {
		this.deregister(conn)
	}
}
