package main

import (
	"github.com/oylshe1314/framework/client/sd"
	"github.com/oylshe1314/framework/client/sd/zk"
	"github.com/oylshe1314/framework/server"
	"github.com/oylshe1314/framework/util"
)

type testNetServer struct {
	server.LoggerServer

	NetServer server.NetServer

	SdConfig       sd.Config
	RegisterClient sd.RegisterClient
}

func (this *testNetServer) Init() error {
	var err error
	err = this.LoggerServer.Init()
	if err != nil {
		return err
	}

	this.NetServer.SetServer(this)
	err = this.NetServer.Init()
	if err != nil {
		return err
	}

	this.RegisterClient = zk.NewRegisterClient(&this.SdConfig)
	this.RegisterClient.SetServer(this)
	this.RegisterClient.SetListener(&this.NetServer.Listener, nil)
	err = this.RegisterClient.Init()
	if err != nil {
		return err
	}

	return err
}

func (this *testNetServer) Serve() error {
	return util.WaitAny(this.NetServer.Serve, this.RegisterClient.Work)
}

func (this *testNetServer) Close() error {
	_ = this.RegisterClient.Close()
	_ = this.NetServer.Close()
	return nil
}

func newTestNetServer() server.Server {
	return &testNetServer{}
}

func main() {
	server.Start(newTestNetServer())
}
