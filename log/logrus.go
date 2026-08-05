package log

import (
	"fmt"
	"io"

	"github.com/oylshe1314/framework/store"
	"github.com/sirupsen/logrus"
)

type logrusEntry struct {
	ll *logrus.Logger

	fields []*store.Pair[string, any]
}

func (this *logrusEntry) logrusLog(level logrus.Level, msg string) {
	var le = logrus.NewEntry(this.ll)
	if len(this.fields) != 0 {
		le.Data["fields"] = this.fields
	}
	le.Logln(level, msg)
}

func (this *logrusEntry) log(level Level, args ...any) {
	this.logrusLog(level.logrusLevel(), fmt.Sprint(args...))
}

func (this *logrusEntry) logf(level Level, format string, args ...any) {
	this.logrusLog(level.logrusLevel(), fmt.Sprintf(format, args...))
}

func (this *logrusEntry) WithField(name string, value any) entry {
	this.fields = append(this.fields, store.NewPair[string, any](name, value))
	return this
}

func (this *logrusEntry) Panic(args ...any) {
	this.log(LevelPanic, fmt.Sprint(args...))
}

func (this *logrusEntry) Panicf(format string, args ...any) {
	this.log(LevelPanic, fmt.Sprintf(format, args...))
}

func (this *logrusEntry) Fatal(args ...any) {
	this.log(LevelFatal, fmt.Sprint(args...))
}

func (this *logrusEntry) Fatalf(format string, args ...any) {
	this.log(LevelFatal, fmt.Sprintf(format, args...))
}

func (this *logrusEntry) Error(args ...any) {
	this.log(LevelError, fmt.Sprint(args...))
}

func (this *logrusEntry) Errorf(format string, args ...any) {
	this.log(LevelError, fmt.Sprintf(format, args...))
}

func (this *logrusEntry) Warn(args ...any) {
	this.log(LevelWarn, fmt.Sprint(args...))
}

func (this *logrusEntry) Warnf(format string, args ...any) {
	this.log(LevelWarn, fmt.Sprintf(format, args...))
}

func (this *logrusEntry) Info(args ...any) {
	this.log(LevelInfo, fmt.Sprint(args...))
}

func (this *logrusEntry) Infof(format string, args ...any) {
	this.log(LevelInfo, fmt.Sprintf(format, args...))
}

func (this *logrusEntry) Debug(args ...any) {
	this.log(LevelDebug, fmt.Sprint(args...))
}

func (this *logrusEntry) Debugf(format string, args ...any) {
	this.log(LevelDebug, fmt.Sprintf(format, args...))
}

func (this *logrusEntry) Trace(args ...any) {
	this.log(LevelTrace, fmt.Sprint(args...))
}

func (this *logrusEntry) Tracef(format string, args ...any) {
	this.log(LevelTrace, fmt.Sprintf(format, args...))
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

	ll.ll = logrus.New()

	ll.ll.SetOutput(ll.writer)
	ll.ll.SetReportCaller(false)
	ll.ll.SetLevel(ll.level.logrusLevel())
	ll.ll.SetFormatter(&logrusJsonFormatter{option: option})

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
	var le = this.entry()
	le.WithField(name, value)
	return le
}

func (this *logrusLogger) Panic(args ...any) {
	this.entry().log(LevelPanic, args...)
}

func (this *logrusLogger) Panicf(format string, args ...any) {
	this.entry().logf(LevelPanic, format, args...)
}

func (this *logrusLogger) Fatal(args ...any) {
	this.entry().log(LevelFatal, args...)
}

func (this *logrusLogger) Fatalf(format string, args ...any) {
	this.entry().logf(LevelFatal, format, args...)
}

func (this *logrusLogger) Error(args ...any) {
	this.entry().log(LevelError, args...)
}

func (this *logrusLogger) Errorf(format string, args ...any) {
	this.entry().logf(LevelError, format, args...)
}

func (this *logrusLogger) Warn(args ...any) {
	this.entry().log(LevelWarn, args...)
}

func (this *logrusLogger) Warnf(format string, args ...any) {
	this.entry().logf(LevelWarn, format, args...)
}

func (this *logrusLogger) Info(args ...any) {
	this.entry().log(LevelInfo, args...)
}

func (this *logrusLogger) Infof(format string, args ...any) {
	this.entry().logf(LevelInfo, format, args...)
}

func (this *logrusLogger) Debug(args ...any) {
	this.entry().log(LevelDebug, args...)
}

func (this *logrusLogger) Debugf(format string, args ...any) {
	this.entry().logf(LevelDebug, format, args...)
}

func (this *logrusLogger) Trace(args ...any) {
	this.entry().log(LevelTrace, args...)
}

func (this *logrusLogger) Tracef(format string, args ...any) {
	this.entry().logf(LevelTrace, format, args...)
}
