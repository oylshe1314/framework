package zk

import "github.com/oylshe1314/framework/log"

type internalLogger struct {
	log.Logger
}

func (this *internalLogger) Printf(format string, args ...interface{}) {
	this.Infof(format, args...)
}
