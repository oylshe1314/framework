package zk

import (
	"github.com/go-zookeeper/zk"
	json "github.com/json-iterator/go"
	"github.com/oylshe1314/framework/client/sd"
	"github.com/oylshe1314/framework/errors"
	"github.com/oylshe1314/framework/server"
	"github.com/oylshe1314/framework/util"
	"strconv"
	"strings"
	"time"
)

func NewRegisterClient(config *sd.Config) sd.RegisterClient {
	return &registerClient{client: client{config: config}}
}

type registerClient struct {
	client

	version int32
	svrPath string
	svrNode *sd.ServerNode
}

func (this *registerClient) SetListener(inner, exter *server.Listener) {
	if inner == nil && exter == nil {
		return
	}
	this.svrNode = sd.NewServiceNode(this.server.Name(), this.server.AppId(), inner, exter)
}

func (this *registerClient) Init() (err error) {
	if this.server == nil {
		return errors.Error("Service register-discovery client init 'server' can not be nil")
	}

	if this.svrNode == nil {
		return errors.Error("please set service node before init")
	}

	this.client.connectHandler = this.register
	this.client.closeHandler = this.deregister
	return this.client.Init()
}

func (this *registerClient) createParentNodes(conn *zk.Conn, path string) error {
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

func (this *registerClient) setServiceNode(conn *zk.Conn) (string, error) {
	var node = this.svrNode

	this.version = 0
	this.svrPath = ""

	if len(node.Guid) == 0 {
		node.Guid = util.UUID()
	}

	data, err := json.Marshal(node)
	if err != nil {
		return "", err
	}

	var servicePath = this.rootPath + "/" + node.Name + serviceServicePath
	var nodesPath = this.rootPath + "/" + node.Name + serviceNodesPath
	var appIdPath = servicePath + "/" + strconv.Itoa(int(node.AppId))

	var version int32
	var nodeExisted = false
	bs, stat, err := conn.Get(appIdPath)
	if err != nil {
		if !errors.Is(err, zk.ErrNoNode) {
			return "", err
		}
	} else {
		_, _, err = conn.Get(string(bs))
		if err == nil {
			return "", errors.Errorf("service '%s:%d' is already existed", node.Name, node.AppId)
		}
		if !errors.Is(err, zk.ErrNoNode) {
			return "", err
		}
		version = stat.Version
		nodeExisted = true
	}

	if nodeExisted {
		tmpPath, err := conn.CreateProtectedEphemeralSequential(nodesPath+"/", data, zk.WorldACL(zk.PermAll))
		if err != nil {
			return "", err
		}

		stat, err = conn.Set(appIdPath, []byte(tmpPath), version)
		if err != nil {
			return "", err
		}

		this.svrPath = appIdPath
		this.version = stat.Version
		return this.svrPath, nil
	} else {
		err = this.createParentNodes(conn, nodesPath)
		if err != nil {
			return "", err
		}

		tmpPath, err := conn.CreateProtectedEphemeralSequential(nodesPath+"/", data, zk.WorldACL(zk.PermAll))
		if err != nil {
			return "", err
		}

		err = this.createParentNodes(conn, servicePath)
		if err != nil {
			return "", err
		}

		newPath, err := conn.Create(appIdPath, []byte(tmpPath), 0, zk.WorldACL(zk.PermAll))
		if err != nil {
			return "", err
		}
		this.svrPath = newPath
		return this.svrPath, nil
	}
}

func (this *registerClient) register(conn *zk.Conn) {
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

func (this *registerClient) deleteServiceNode(conn *zk.Conn) {
	if len(this.svrPath) > 0 {
		_ = conn.Delete(this.svrPath, this.version)
	}
}

func (this *registerClient) deregister(conn *zk.Conn) {
	if conn != nil && conn.State() >= zk.StateConnected {
		this.deleteServiceNode(conn)
	}
}
