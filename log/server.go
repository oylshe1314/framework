package log

import (
	"context"
	"framework"
	"framework/errors"
)

type LoggerServer struct {
	option *Option

	logger Logger
}

func (this *LoggerServer) SetOption(option *Option) {
	this.option = option
}

func (this *LoggerServer) Init(ctx context.Context) error {
	if this.option == nil {
		return errors.New("'LoggerServer' option is nil")
	}

	var name = "server"

	var namedServer = framework.ServerFromContext[*framework.NamedServer](ctx, "namedServer")
	if namedServer != nil {
		name = namedServer.Name()
	}

	writer, err := NewFileWriter(name, this.option)
	if err != nil {
		return err
	}

	logger, err := New(writer, this.option)
	if err != nil {
		return err
	}

	this.logger = logger
	return nil
}

func (this *LoggerServer) Start() error {
	return nil
}

func (this *LoggerServer) Close() error {
	if this.logger != nil {
		return this.logger.Close()
	}
	return nil
}

func (this *LoggerServer) Logger() Logger {
	return this.logger
}
