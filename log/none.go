package log

import (
	"fmt"
	"os"
)

type noneLogger struct {
}

func NewNoneLogger() Logger {
	return &noneLogger{}
}

func (this *noneLogger) Panic(args ...any) {
	panic(fmt.Sprint(args...))
}

func (this *noneLogger) Panicf(format string, args ...any) {
	panic(fmt.Sprintf(format, args...))
}

func (this *noneLogger) Fatal(args ...any) {
	os.Exit(1)
}

func (this *noneLogger) Fatalf(format string, args ...any) {
	os.Exit(1)
}

func (this *noneLogger) Error(args ...any) {

}

func (this *noneLogger) Errorf(format string, args ...any) {

}

func (this *noneLogger) Warn(args ...any) {

}

func (this *noneLogger) Warnf(format string, args ...any) {

}

func (this *noneLogger) Info(args ...any) {

}

func (this *noneLogger) Infof(format string, args ...any) {

}

func (this *noneLogger) Debug(args ...any) {

}

func (this *noneLogger) Debugf(format string, args ...any) {

}

func (this *noneLogger) Trace(args ...any) {

}

func (this *noneLogger) Tracef(format string, args ...any) {

}

func (this *noneLogger) WithField(key string, value any) entry {
	return this
}

func (this *noneLogger) IsLevelEnabled(level Level) bool {
	return false
}

func (this *noneLogger) IsPanicEnabled() bool {
	return false
}

func (this *noneLogger) IsFatalEnabled() bool {
	return false
}

func (this *noneLogger) IsErrorEnabled() bool {
	return false
}

func (this *noneLogger) IsWarnEnabled() bool {
	return false
}

func (this *noneLogger) IsInfoEnabled() bool {
	return false
}

func (this *noneLogger) IsDebugEnabled() bool {
	return false
}

func (this *noneLogger) IsTraceEnabled() bool {
	return false
}

func (this *noneLogger) Close() error {
	return nil
}
