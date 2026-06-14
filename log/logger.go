package log

import (
	"io"

	"github.com/oylshe1314/framework/errors"
)

type basic interface {
	Panic(args ...any)
	Panicf(format string, args ...any)

	Fatal(args ...any)
	Fatalf(format string, args ...any)

	Error(args ...any)
	Errorf(format string, args ...any)

	Warn(args ...any)
	Warnf(format string, args ...any)

	Info(args ...any)
	Infof(format string, args ...any)

	Debug(args ...any)
	Debugf(format string, args ...any)

	Trace(args ...any)
	Tracef(format string, args ...any)
}

type leveled interface {
	IsLevelEnabled(level Level) bool
	IsPanicEnabled() bool
	IsFatalEnabled() bool
	IsErrorEnabled() bool
	IsWarnEnabled() bool
	IsInfoEnabled() bool
	IsDebugEnabled() bool
	IsTraceEnabled() bool
}

type entry interface {
	//basic

	WithField(key string, value any) entry
}

type Logger interface {
	entry
	leveled

	io.Closer
}

func NewLogger(writer io.Writer, option *Option) (Logger, error) {
	switch option.Logger {
	case "logrus":
		return newLogrusLogger(writer, option)
	case "zap":
		return newZapLogger()
	case "zerolog":
		return newZeroLogger(writer, option)
	default:
		return nil, errors.Errorf("unknown logger type '%s'", option.Logger)
	}
}

func NewFileWriter(option *Option) (io.Writer, error) {
	return newDailyWriter(option)
}
