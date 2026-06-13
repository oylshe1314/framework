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

func (this *zeroEntry) log(event *zerolog.Event, msg string) {
	if this.fl != nil {
		for _, it := range this.fl.fields {
			event = event.Any(it.name, it.value)
		}
	}
	event.Msg(msg)
}

func (this *zeroEntry) WithField(name string, value any) entry {
	if this.fl == nil {
		this.fl = &fieldLogger{entry: this, fields: []*field{{name: name, value: value}}}
	} else {
		this.fl.WithField(name, value)
	}
	return this
}

func (this *zeroEntry) Panic(args ...any) {
	this.log(this.zl.Panic(), fmt.Sprint(args...))
}

func (this *zeroEntry) Panicf(format string, args ...any) {
	this.log(this.zl.Panic(), fmt.Sprintf(format, args...))
}

func (this *zeroEntry) Fatal(args ...any) {
	this.log(this.zl.Fatal(), fmt.Sprint(args...))
}

func (this *zeroEntry) Fatalf(format string, args ...any) {
	this.log(this.zl.Fatal(), fmt.Sprintf(format, args...))
}

func (this *zeroEntry) Error(args ...any) {
	this.log(this.zl.Error(), fmt.Sprint(args...))
}

func (this *zeroEntry) Errorf(format string, args ...any) {
	this.log(this.zl.Error(), fmt.Sprintf(format, args...))
}

func (this *zeroEntry) Warn(args ...any) {
	this.log(this.zl.Warn(), fmt.Sprint(args...))
}

func (this *zeroEntry) Warnf(format string, args ...any) {
	this.log(this.zl.Warn(), fmt.Sprintf(format, args...))
}

func (this *zeroEntry) Info(args ...any) {
	this.log(this.zl.Info(), fmt.Sprint(args...))
}

func (this *zeroEntry) Infof(format string, args ...any) {
	this.log(this.zl.Info(), fmt.Sprintf(format, args...))
}

func (this *zeroEntry) Debug(args ...any) {
	this.log(this.zl.Debug(), fmt.Sprint(args...))
}

func (this *zeroEntry) Debugf(format string, args ...any) {
	this.log(this.zl.Debug(), fmt.Sprintf(format, args...))
}

func (this *zeroEntry) Trace(args ...any) {
	this.log(this.zl.Trace(), fmt.Sprint(args...))
}

func (this *zeroEntry) Tracef(format string, args ...any) {
	this.log(this.zl.Trace(), fmt.Sprintf(format, args...))
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
