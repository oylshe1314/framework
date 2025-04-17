package server

import (
	"github.com/oylshe1314/framework/errors"
	"github.com/oylshe1314/framework/util"
	"net"
)

type Listener struct {
	network string
	bind    string
	address string
	extra   map[string]any

	l net.Listener
}

func (this *Listener) WithNetwork(network string) {
	this.network = network
}

func (this *Listener) WithBind(bind string) {
	this.bind = bind
}

func (this *Listener) WithAddress(address string) {
	this.address = address
}

func (this *Listener) WithExtra(extra map[string]any) {
	this.extra = extra
}

func (this *Listener) Network() string {
	return this.network
}

func (this *Listener) Address() string {
	return this.address
}

func (this *Listener) Bind() string {
	return this.bind
}

func (this *Listener) Extra() map[string]any {
	return this.extra
}

func (this *Listener) Init() (err error) {
	if util.Unix() >= expiration {
		return errors.Error("the server was expired")
	}

	if len(this.network) == 0 {
		return errors.Error("'network' cannot be empty")
	}

	if len(this.bind) == 0 {
		return errors.Error("'bind' cannot be empty")
	}

	if len(this.address) == 0 {
		return errors.Error("'address' cannot be empty")
	}

	var addr net.Addr
	switch this.network {
	case "tcp":
		addr, err = net.ResolveTCPAddr(this.network, this.bind)
	case "udp":
		addr, err = net.ResolveUDPAddr(this.network, this.bind)
	case "unix":
		addr, err = net.ResolveUnixAddr(this.network, this.bind)
	default:
		return errors.Errorf("unknown network '%s'", this.network)
	}
	if err != nil {
		return err
	}

	this.network = addr.Network()
	this.bind = addr.String()

	return nil
}

func (this *Listener) Listen() (err error) {
	this.l, err = net.Listen(this.network, this.bind)
	if err != nil {
		return err
	}

	this.bind = this.l.Addr().String()
	return
}

func (this *Listener) Close() (err error) {
	if this.l != nil {
		err = this.l.Close()
		this.l = nil
	}
	return
}
