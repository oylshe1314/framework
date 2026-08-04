package log

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/rs/zerolog"
)

type zeroEntry struct {
	zl *zerolog.Logger
	fl *fieldLogger
}

func (this *zeroEntry) zeroLog(event *zerolog.Event, msg string) {
	if this.fl != nil {
		for _, it := range this.fl.fields {
			event = event.Any(it.name, it.value)
		}
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
	if this.fl == nil {
		this.fl = &fieldLogger{entry: this, fields: []*field{{name: name, value: value}}}
	} else {
		this.fl.WithField(name, value)
	}
	return this
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
		}
	}

	if option.WithCaller {
		zl.zl = zl.zl.With().CallerWithSkipFrameCount(4).Logger()
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
	ze.fl = &fieldLogger{entry: ze, fields: []*field{{name: name, value: value}}}
	return ze
}

func (this *zeroLogger) Panic(args ...any) {
	this.entry().log(LevelPanic, args...)
}

func (this *zeroLogger) Panicf(format string, args ...any) {
	this.entry().logf(LevelPanic, format, args...)
}

func (this *zeroLogger) Fatal(args ...any) {
	this.entry().log(LevelFatal, args...)
}

func (this *zeroLogger) Fatalf(format string, args ...any) {
	this.entry().logf(LevelFatal, format, args...)
}

func (this *zeroLogger) Error(args ...any) {
	this.entry().log(LevelError, args...)
}

func (this *zeroLogger) Errorf(format string, args ...any) {
	this.entry().logf(LevelError, format, args...)
}

func (this *zeroLogger) Warn(args ...any) {
	this.entry().log(LevelWarn, args...)
}

func (this *zeroLogger) Warnf(format string, args ...any) {
	this.entry().logf(LevelWarn, format, args...)
}

func (this *zeroLogger) Info(args ...any) {
	this.entry().log(LevelInfo, args...)
}

func (this *zeroLogger) Infof(format string, args ...any) {
	this.entry().logf(LevelInfo, format, args...)
}

func (this *zeroLogger) Debug(args ...any) {
	this.entry().log(LevelDebug, args...)
}

func (this *zeroLogger) Debugf(format string, args ...any) {
	this.entry().logf(LevelDebug, format, args...)
}

func (this *zeroLogger) Trace(args ...any) {
	this.entry().log(LevelTrace, args...)
}

func (this *zeroLogger) Tracef(format string, args ...any) {
	this.entry().logf(LevelTrace, format, args...)
}
