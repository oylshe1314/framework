package main

import (
	"fmt"
	"github.com/oylshe1314/framework/client/rpc"
	"github.com/oylshe1314/framework/client/sd"
	"github.com/oylshe1314/framework/client/sd/zk"
	"github.com/oylshe1314/framework/options"
)

type testNetClient struct {
	SdConfig        sd.Config
	SubscribeClient sd.SubscribeClient
	NetRpcClient    rpc.NetRpcClient
}

func (this *testNetClient) Init() error {
	var err error
	this.SubscribeClient = zk.NewSubscribeClient(&this.SdConfig)
	this.SubscribeClient.AddSubscribe("test", this.NetRpcClient.SubscribeCallback)
	err = this.SubscribeClient.Init()
	if err != nil {
		return err
	}
	return nil
}

func (this *testNetClient) Close() error {
	_ = this.SubscribeClient.Close()
	_ = this.NetRpcClient.Close()
	return nil
}

func main() {
	var client = &testNetClient{}
	opts, err := options.ReadOptions("./conf/net.json")
	if err != nil {
		fmt.Println(err)
		return
	}

	err = opts.Init(client)
	if err != nil {
		fmt.Println(err)
		return
	}
}
