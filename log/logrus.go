package log

import (
	"fmt"
	"io"

	"github.com/sirupsen/logrus"
)

type logrusEntry struct {
	ll *logrus.Logger
	fl *fieldLogger
}

func (this *logrusEntry) logrusLog(level logrus.Level, msg string) {
	var le = logrus.NewEntry(this.ll)
	if this.fl != nil && len(this.fl.fields) != 0 {
		le.Data["fields"] = this.fl.fields
	}
	this.ll.Logln(level, msg)
}

func (this *logrusEntry) log(level Level, args ...any) {
	this.logrusLog(level.logrusLevel(), fmt.Sprint(args...))
}

func (this *logrusEntry) logf(level Level, format string, args ...any) {
	this.logrusLog(level.logrusLevel(), fmt.Sprintf(format, args...))
}

func (this *logrusEntry) WithField(name string, value any) entry {
	if this.fl == nil {
		this.fl = &fieldLogger{entry: this, fields: []*field{{name: name, value: value}}}
	} else {
		this.fl.WithField(name, value)
	}
	return this
}

type logrusLogger struct {
	*leveledLogger

	writer io.Writer

	ll *logrus.Logger
}

func newLogrusLogger(writer io.Writer, option *Option) (Logger, error) {
	var ll = &logrusLogger{leveledLogger: &leveledLogger{level: ParseLevel(option.Level)}}

	ll.writer = writer

	if option.WithConsole {
		ll.ll.AddHook(&logrusConsoleHook{})
	}

	var l = logrus.New()

	l.SetOutput(ll.writer)
	l.SetReportCaller(option.WithCaller)
	l.SetLevel(ll.level.logrusLevel())
	l.SetFormatter(&logrusJsonFormatter{option: option})

	return ll, nil
}

func (this *logrusLogger) Close() error {
	if closer, ok := this.writer.(io.Closer); ok {
		return closer.Close()
	}
	return nil
}

func (this *logrusLogger) entry() *logrusEntry {
	return &logrusEntry{ll: this.ll}
}

func (this *logrusLogger) WithField(name string, value any) entry {
	var ze = this.entry()
	ze.fl = &fieldLogger{entry: ze, fields: []*field{{name: name, value: value}}}
	return ze
}

func (this *logrusLogger) Panic(args ...any) {
	this.entry().log(LevelPanic, args)
}

func (this *logrusLogger) Panicf(format string, args ...any) {
	this.entry().logf(LevelPanic, format, args)
}

func (this *logrusLogger) Fatal(args ...any) {
	this.entry().log(LevelFatal, args)
}

func (this *logrusLogger) Fatalf(format string, args ...any) {
	this.entry().logf(LevelFatal, format, args)
}

func (this *logrusLogger) Error(args ...any) {
	this.entry().log(LevelError, args)
}

func (this *logrusLogger) Errorf(format string, args ...any) {
	this.entry().logf(LevelError, format, args)
}

func (this *logrusLogger) Warn(args ...any) {
	this.entry().log(LevelWarn, args)
}

func (this *logrusLogger) Warnf(format string, args ...any) {
	this.entry().logf(LevelWarn, format, args)
}

func (this *logrusLogger) Info(args ...any) {
	this.entry().log(LevelInfo, args)
}

func (this *logrusLogger) Infof(format string, args ...any) {
	this.entry().logf(LevelInfo, format, args)
}

func (this *logrusLogger) Debug(args ...any) {
	this.entry().log(LevelDebug, args)
}

func (this *logrusLogger) Debugf(format string, args ...any) {
	this.entry().logf(LevelDebug, format, args)
}

func (this *logrusLogger) Trace(args ...any) {
	this.entry().log(LevelTrace, args)
}

func (this *logrusLogger) Tracef(format string, args ...any) {
	this.entry().logf(LevelTrace, format, args)
}
