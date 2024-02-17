package zk

import (
	"framework/client/sd"
	"framework/errors"
	"framework/util"
	"github.com/go-zookeeper/zk"
	json "github.com/json-iterator/go"
	"strconv"
	"strings"
	"time"
)

type RegisterClient struct {
	Client

	version int32
	svcPath string
	svcNode *sd.ServiceNode
}

func (this *RegisterClient) SetServiceNode(node *sd.ServiceNode) {
	this.svcNode = node
}

func (this *RegisterClient) Init() (err error) {
	if this.svcNode == nil {
		return errors.Error("please set service node before init")
	}

	this.Client.connectHandler = this.register
	this.Client.closeHandler = this.deregister
	return this.Client.Init()
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
	var node = this.svcNode

	this.version = 0
	this.svcPath = ""

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

		this.svcPath = appIdPath
		this.version = stat.Version
		return this.svcPath, nil
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
		this.svcPath = newPath
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

func (this *RegisterClient) deleteServiceNode(conn *zk.Conn) {
	if len(this.svcPath) > 0 {
		_ = conn.Delete(this.svcPath, this.version)
	}
}

func (this *RegisterClient) deregister(conn *zk.Conn) {
	if conn != nil && conn.State() >= zk.StateConnected {
		this.deleteServiceNode(conn)
	}
}
