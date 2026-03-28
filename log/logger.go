package log

import "io"

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
	basic

	WithField(key string, value any) entry
}

type Logger interface {
	entry
	leveled

	io.Closer
}

func New(writer io.Writer, option *Option) (Logger, error) {
	switch option.Logger {
	case "zerolog":
		var options []zeroOption
		if option.WithConsole {
			options = append(options, zeroWithWithConsole())
		}
		if option.WithTime {
			options = append(options, zeroWithTime())
			if len(option.Timezone) > 0 {
				options = append(options, zeroWithTimezone(option.Timezone))
			}
			if len(option.TimeFormat) > 0 {
				options = append(options, zeroWithTimeFormat(option.TimeFormat))
			}
		}
		if option.WithCaller {
			options = append(options, zeroWithCaller())
			if option.CallerSkip > 0 {
				options = append(options, zeroWithCallerSkip(option.CallerSkip))
			}
		}

		return newZeroLogger(ParseLevel(option.Level), writer, options...)
	case "zap":
		return newZapLogger()
	default:
		return newLogrusLogger()
	}
}

func NewFileWriter(option *Option) (io.Writer, error) {
	return newDailyWriter(option)
}
