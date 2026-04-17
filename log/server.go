package log

import (
	"context"
	"github.com/oylshe1314/framework/errors"
	"github.com/oylshe1314/framework/option"
)

type LoggerServer struct {
	option.Optional[Option]

	logger Logger
}

func (this *LoggerServer) Init(ctx context.Context) error {
	if this.GetOption() == nil {
		return errors.New("option is nil")
	}

	writer, err := NewFileWriter(this.GetOption())
	if err != nil {
		return err
	}

	logger, err := NewLogger(writer, this.GetOption())
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
