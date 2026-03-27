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
		for name, value := range this.fl.fields {
			event = event.Any(name, value)
		}
	}
	event.Msg(msg)
}

func (this *zeroEntry) WithField(name string, value any) entry {
	this.fl.fields[name] = value
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

	options map[string]any
}

type zeroOption func(*zeroLogger)

func zeroWithWithConsole() zeroOption {
	return func(zl *zeroLogger) {
		zl.options["withConsole"] = true
	}
}

func zeroWithTime() zeroOption {
	return func(zl *zeroLogger) {
		zl.options["time"] = true
	}
}

func zeroWithTimezone(timezone string) zeroOption {
	return func(zl *zeroLogger) {
		zl.options["timezone"] = timezone
	}
}

func zeroWithTimeFormat(timeFormat string) zeroOption {
	return func(zl *zeroLogger) {
		zl.options["timeFormat"] = timeFormat
	}
}

func zeroWithCaller() zeroOption {
	return func(zl *zeroLogger) {
		zl.options["caller"] = true
	}
}

func zeroWithCallerSkip(callerSkip int) zeroOption {
	return func(zl *zeroLogger) {
		zl.options["callerSkip"] = callerSkip
	}
}

func newZeroLogger(level Level, writer io.Writer, options ...zeroOption) (Logger, error) {
	var zl = &zeroLogger{}
	for _, option := range options {
		option(zl)
	}

	zl.leveledLogger = &leveledLogger{level: level}

	zl.writer = writer
	if _, ok := zl.options["withConsole"]; ok {
		zl.writer = io.MultiWriter(os.Stdout, zl.writer)
	}

	zl.zl = zerolog.New(zl.writer)
	if _, ok := zl.options["time"]; ok {
		zl.zl = zl.zl.With().Timestamp().Logger()
		if timezone, ok := zl.options["timezone"].(string); ok && timezone != "" {
			location, err := time.LoadLocation(timezone)
			if err != nil {
				return nil, err
			}

			zerolog.TimestampFunc = func() time.Time {
				return time.Now().In(location)
			}
		}

		if timeFormat, ok := zl.options["timeFormat"].(string); ok && timeFormat != "" {
			zerolog.TimeFieldFormat = timeFormat
		}
	}

	if _, ok := zl.options["caller"]; ok {
		if callerSkip, ok := zl.options["callerSkip"].(int); ok && callerSkip > 0 {
			zl.zl = zl.zl.With().CallerWithSkipFrameCount(callerSkip).Logger()
		} else {
			zl.zl = zl.zl.With().Caller().Logger()
		}
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
	ze.fl = &fieldLogger{entry: ze, fields: map[string]any{name: value}}
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
