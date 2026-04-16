package zk

import "framework/log"

type logger struct {
	log.Logger
}

func (this *logger) Printf(format string, args ...interface{}) {
	this.Infof(format, args...)
}
