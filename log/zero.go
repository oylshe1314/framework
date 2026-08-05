package log

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/oylshe1314/framework/store"
	"github.com/rs/zerolog"
)

type zeroEntry struct {
	zl *zerolog.Logger

	fields []*store.Pair[string, any]
}

func (this *zeroEntry) zeroLog(event *zerolog.Event, msg string) {
	for _, field := range this.fields {
		event = event.Any(field.Key, field.Value)
	}
	event.Msg(msg)
}

func (this *zeroEntry) zeroEvent(level Level) *zerolog.Event {
	switch level {
	case LevelPanic:
		return this.zl.Panic()
	case LevelFatal:
		return this.zl.Fatal()
	case LevelError:
		return this.zl.Error()
	case LevelWarn:
		return this.zl.Warn()
	case LevelInfo:
		return this.zl.Info()
	case LevelDebug:
		return this.zl.Debug()
	case LevelTrace:
		return this.zl.Trace()
	default:
		return this.zl.Info()
	}
}

func (this *zeroEntry) log(level Level, args ...any) {
	this.zeroLog(this.zeroEvent(level), fmt.Sprint(args...))
}

func (this *zeroEntry) logf(level Level, format string, args ...any) {
	this.zeroLog(this.zeroEvent(level), fmt.Sprintf(format, args...))
}

func (this *zeroEntry) WithField(name string, value any) entry {
	this.fields = append(this.fields, store.NewPair[string, any](name, value))
	return this
}

func (this *zeroEntry) Panic(args ...any) {
	this.log(LevelPanic, fmt.Sprint(args...))
}

func (this *zeroEntry) Panicf(format string, args ...any) {
	this.log(LevelPanic, fmt.Sprintf(format, args...))
}

func (this *zeroEntry) Fatal(args ...any) {
	this.log(LevelFatal, fmt.Sprint(args...))
}

func (this *zeroEntry) Fatalf(format string, args ...any) {
	this.log(LevelFatal, fmt.Sprintf(format, args...))
}

func (this *zeroEntry) Error(args ...any) {
	this.log(LevelError, fmt.Sprint(args...))
}

func (this *zeroEntry) Errorf(format string, args ...any) {
	this.log(LevelError, fmt.Sprintf(format, args...))
}

func (this *zeroEntry) Warn(args ...any) {
	this.log(LevelWarn, fmt.Sprint(args...))
}

func (this *zeroEntry) Warnf(format string, args ...any) {
	this.log(LevelWarn, fmt.Sprintf(format, args...))
}

func (this *zeroEntry) Info(args ...any) {
	this.log(LevelInfo, fmt.Sprint(args...))
}

func (this *zeroEntry) Infof(format string, args ...any) {
	this.log(LevelInfo, fmt.Sprintf(format, args...))
}

func (this *zeroEntry) Debug(args ...any) {
	this.log(LevelDebug, fmt.Sprint(args...))
}

func (this *zeroEntry) Debugf(format string, args ...any) {
	this.log(LevelDebug, fmt.Sprintf(format, args...))
}

func (this *zeroEntry) Trace(args ...any) {
	this.log(LevelTrace, fmt.Sprint(args...))
}

func (this *zeroEntry) Tracef(format string, args ...any) {
	this.log(LevelTrace, fmt.Sprintf(format, args...))
}

type zeroLogger struct {
	*leveledLogger

	writer io.Writer

	zl zerolog.Logger
}

func newZeroLogger(writer io.Writer, option *Option) (Logger, error) {
	var zl = &zeroLogger{leveledLogger: &leveledLogger{level: ParseLevel(option.Level)}}

	zl.writer = writer
	if option.WithConsole {
		zl.writer = io.MultiWriter(os.Stdout, zl.writer)
	}

	zl.zl = zerolog.New(zl.writer)
	if option.WithTimestamp {
		zl.zl = zl.zl.With().Timestamp().Logger()
		if option.Timezone != "" {
			location, err := time.LoadLocation(option.Timezone)
			if err != nil {
				return nil, err
			}

			zerolog.TimestampFunc = func() time.Time {
				return time.Now().In(location)
			}
		}

		if option.TimeFormat != "" {
			zerolog.TimeFieldFormat = option.TimeFormat
		} else {
			zerolog.TimeFieldFormat = time.RFC3339
		}
	}

	if option.WithCaller {
		zl.zl = zl.zl.With().CallerWithSkipFrameCount(5).Logger()
	}

	return zl, nil
}

func (this *zeroLogger) Close() error {
	if closer, ok := this.writer.(io.Closer); ok {
		return closer.Close()
	}
	return nil
}

func (this *zeroLogger) entry() *zeroEntry {
	return &zeroEntry{zl: &this.zl}
}

func (this *zeroLogger) WithField(name string, value any) entry {
	var ze = this.entry()
	ze.WithField(name, value)
	return ze
}

func (this *zeroLogger) Panic(args ...any) {
	this.entry().Panic(args...)
}

func (this *zeroLogger) Panicf(format string, args ...any) {
	this.entry().Panicf(format, args...)
}

func (this *zeroLogger) Fatal(args ...any) {
	this.entry().Fatal(args...)
}

func (this *zeroLogger) Fatalf(format string, args ...any) {
	this.entry().Fatalf(format, args...)
}

func (this *zeroLogger) Error(args ...any) {
	this.entry().Error(args...)
}

func (this *zeroLogger) Errorf(format string, args ...any) {
	this.entry().Errorf(format, args...)
}

func (this *zeroLogger) Warn(args ...any) {
	this.entry().Warn(args...)
}

func (this *zeroLogger) Warnf(format string, args ...any) {
	this.entry().Warnf(format, args...)
}

func (this *zeroLogger) Info(args ...any) {
	this.entry().Info(args...)
}

func (this *zeroLogger) Infof(format string, args ...any) {
	this.entry().Infof(format, args...)
}

func (this *zeroLogger) Debug(args ...any) {
	this.entry().Debug(args...)
}

func (this *zeroLogger) Debugf(format string, args ...any) {
	this.entry().Debugf(format, args...)
}

func (this *zeroLogger) Trace(args ...any) {
	this.entry().Trace(args...)
}

func (this *zeroLogger) Tracef(format string, args ...any) {
	this.entry().Tracef(format, args...)
}
