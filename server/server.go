package server

import (
	"framework/errors"
	"framework/log"
	"framework/util"
	"time"
)

type Server interface {
	OptionalServer

	Name() string
	AppId() uint32
	Close() error
	Serve() error
	Logger() log.Logger
}

type LoggerServer struct {
	name  string
	appId uint32

	logDir     string
	logLevel   log.Level
	logConsole bool

	logger log.Logger
}

func (this *LoggerServer) WithName(name string) {
	this.name = name
}

func (this *LoggerServer) WithAppId(appId uint32) {
	this.appId = appId
}

func (this *LoggerServer) WithLogDir(logDir string) {
	this.logDir = logDir
}

func (this *LoggerServer) WithLogLevel(logLevel string) {
	this.logLevel = log.LevelOf(logLevel)
}

func (this *LoggerServer) WithLogConsole(logConsole bool) {
	this.logConsole = logConsole
}

func (this *LoggerServer) Name() string {
	return this.name
}

func (this *LoggerServer) AppId() uint32 {
	return this.appId
}

func (this *LoggerServer) Logger() log.Logger {
	return this.logger
}

func (this *LoggerServer) Init() (err error) {
	if util.Unix() >= expiration {
		return errors.Error("the server was expired")
	}

	if len(this.name) == 0 {
		return errors.Error("'name' cannot be empty")
	}

	if this.appId == 0 {
		return errors.Error("'appId' cannot be 0")
	}

	if len(this.logDir) == 0 {
		return errors.Error("'logDir' cannot be empty")
	}

	if this.logLevel > log.LevelTrace {
		return errors.Error("incorrect 'logLevel' value")
	}

	this.logger, err = log.NewDailyLogger(this.Name(), this.AppId(), this.logDir, log.WithLevel(this.logLevel), log.WithConsole(this.logConsole))
	return err
}

func (this *LoggerServer) Close() (err error) {
	if this.logger != nil {
		time.Sleep(time.Second)
		err = this.logger.Close()
		this.logger = nil
	}
	return
}
