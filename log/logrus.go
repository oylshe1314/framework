package log

import (
	"io"

	"github.com/sirupsen/logrus"
)

type logrusEntry struct {
	lv logrus.Level
	ll *logrus.Logger
	fl *fieldLogger
}

func (this *logrusEntry) log(msg string) {
	var le = logrus.NewEntry(this.ll)
	if this.fl != nil && len(this.fl.fields) != 0 {
		le.Data["fields"] = this.fl.fields
	}
	this.ll.Logln(this.lv, msg)
}

func (this *logrusEntry) WithField(name string, value any) entry {
	if this.fl == nil {
		this.fl = &fieldLogger{entry: this, fields: []*field{{name: name, value: value}}}
	} else {
		this.fl.WithField(name, value)
	}
	return this
}

func (this *logrusEntry) Panic(args ...any) {

}

func (this *logrusEntry) Panicf(format string, args ...any) {

}

func (this *logrusEntry) Fatal(args ...any) {

}

func (this *logrusEntry) Fatalf(format string, args ...any) {

}

func (this *logrusEntry) Error(args ...any) {

}

func (this *logrusEntry) Errorf(format string, args ...any) {

}

func (this *logrusEntry) Warn(args ...any) {

}

func (this *logrusEntry) Warnf(format string, args ...any) {

}

func (this *logrusEntry) Info(args ...any) {

}

func (this *logrusEntry) Infof(format string, args ...any) {

}

func (this *logrusEntry) Debug(args ...any) {

}

func (this *logrusEntry) Debugf(format string, args ...any) {

}

func (this *logrusEntry) Trace(args ...any) {

}

func (this *logrusEntry) Tracef(format string, args ...any) {

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

	return ll, nil
}

func (this *logrusLogger) Close() error {
	if closer, ok := this.writer.(io.Closer); ok {
		return closer.Close()
	}
	return nil
}

func (this *logrusLogger) entry() *logrusEntry {
	lv, err := logrus.ParseLevel(this.level.String())
	if err != nil {
		lv = logrus.InfoLevel
	}
	return &logrusEntry{lv: lv, ll: this.ll}
}

func (this *logrusLogger) WithField(name string, value any) entry {
	var ze = this.entry()
	ze.fl = &fieldLogger{entry: ze, fields: []*field{{name: name, value: value}}}
	return ze
}

func (this *logrusLogger) Panic(args ...any) {

}

func (this *logrusLogger) Panicf(format string, args ...any) {

}

func (this *logrusLogger) Fatal(args ...any) {

}

func (this *logrusLogger) Fatalf(format string, args ...any) {

}

func (this *logrusLogger) Error(args ...any) {

}

func (this *logrusLogger) Errorf(format string, args ...any) {

}

func (this *logrusLogger) Warn(args ...any) {

}

func (this *logrusLogger) Warnf(format string, args ...any) {

}

func (this *logrusLogger) Info(args ...any) {

}

func (this *logrusLogger) Infof(format string, args ...any) {

}

func (this *logrusLogger) Debug(args ...any) {

}

func (this *logrusLogger) Debugf(format string, args ...any) {

}

func (this *logrusLogger) Trace(args ...any) {

}

func (this *logrusLogger) Tracef(format string, args ...any) {

}
