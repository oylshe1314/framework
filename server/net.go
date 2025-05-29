package server

import (
	"github.com/oylshe1314/framework/log"
	. "github.com/oylshe1314/framework/net"
	"runtime/debug"
)

type NetServer struct {
	Listener

	ConnMux

	running bool
	logger  log.Logger

	connMap map[*Conn]struct{}
}

func (this *NetServer) SetLogger(logger log.Logger) {
	this.logger = logger
}

func (this *NetServer) serve() error {
	defer func() {
		var err = recover()
		if err != nil {
			this.logger.Error(err)
			this.logger.Error(string(debug.Stack()))
		}
	}()

	for {
		cc, err := this.l.Accept()
		if err != nil {
			return err
		}

		conn := NewConn(cc, this.logger, &this.ConnMux)
		this.connMap[conn] = struct{}{}
		go func() {
			defer func() {
				delete(this.connMap, conn)
			}()
			_ = conn.Serve()
		}()
	}
}

func (this *NetServer) Init() (err error) {
	if this.logger == nil {
		this.logger = log.DefaultLogger
	}

	this.connMap = make(map[*Conn]struct{})
	return this.Listener.Init()
}

func (this *NetServer) Serve() (err error) {

	err = this.Listener.Listen()
	if err != nil {
		return err
	}

	this.logger.Info("NetServer is listening on ", this.Bind())

	this.running = true
	err = this.serve()
	if !this.running {
		return nil
	}

	this.running = false
	return err
}

func (this *NetServer) Close() error {
	this.running = false
	var err = this.Listener.Close()
	for conn := range this.connMap {
		_ = conn.Close()
	}
	return err
}
