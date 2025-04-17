package client

import (
	"github.com/oylshe1314/framework/errors"
	"github.com/oylshe1314/framework/log"
	. "github.com/oylshe1314/framework/net"
	"net"
)

type NetClient struct {
	ConnMux

	network string
	address string

	logger log.Logger

	conn *Conn
}

func (this *NetClient) WithNetwork(network string) {
	this.network = network
}

func (this *NetClient) WithAddress(address string) {
	this.address = address
}

func (this *NetClient) Network() string {
	return this.network
}

func (this *NetClient) Address() string {
	return this.address
}

func (this *NetClient) SetLogger(logger log.Logger) {
	this.logger = logger
}

func (this *NetClient) Init() (err error) {
	if this.logger == nil {
		this.logger = log.DefaultLogger
	}

	if len(this.network) == 0 {
		return errors.Error("'network' cannot be empty")
	}

	if len(this.address) == 0 {
		return errors.Error("'address' cannot be empty")
	}

	var addr net.Addr
	switch this.network {
	case "tcp":
		addr, err = net.ResolveTCPAddr(this.network, this.address)
	case "udp":
		addr, err = net.ResolveUDPAddr(this.network, this.address)
	case "unix":
		addr, err = net.ResolveUnixAddr(this.network, this.address)
	default:
		return errors.Errorf("unknown network '%s'", this.network)
	}
	if err != nil {
		return err
	}

	this.network = addr.Network()
	this.address = addr.String()

	return nil
}

func (this *NetClient) Close() (err error) {
	if this.conn != nil {
		return this.conn.Close()
	}
	return
}

func (this *NetClient) Work() error {
	return this.conn.Serve()
}

func (this *NetClient) Dial() error {
	conn, err := net.Dial(this.network, this.address)
	if err != nil {
		return err
	}
	this.conn = NewConn(conn, this.logger, &this.ConnMux)
	return nil
}

func (this *NetClient) Send(modId, msgId uint16, v interface{}) error {
	if this.conn == nil {
		return errors.Error("please connect server first")
	}
	return this.conn.Send(modId, msgId, v)
}

func (this *NetClient) Read() (*Message, error) {
	if this.conn == nil {
		return nil, errors.Error("please connect server first")
	}
	return this.conn.Read()
}
