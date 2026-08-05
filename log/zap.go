package log

import (
	"fmt"
	"io"
	"time"

	"github.com/oylshe1314/framework/store"
	"github.com/oylshe1314/framework/util"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type zapEntry struct {
	zl *zap.Logger

	fields []*store.Pair[string, any]
}

func (this *zapEntry) zapLog(level zapcore.Level, msg string) {
	var fields []zap.Field
	for _, field := range this.fields {
		fields = append(fields, zap.Any(field.Key, field.Value))
	}
	this.zl.Log(level, msg, fields...)
}

func (this *zapEntry) log(level Level, args ...any) {
	this.zapLog(level.zapLevel(), fmt.Sprint(args...))
}

func (this *zapEntry) logf(level Level, format string, args ...any) {
	this.zapLog(level.zapLevel(), fmt.Sprintf(format, args...))
}

func (this *zapEntry) WithField(name string, value any) entry {
	this.fields = append(this.fields, store.NewPair[string, any](name, value))
	return this
}

func (this *zapEntry) Panic(args ...any) {
	this.log(LevelPanic, fmt.Sprint(args...))
}

func (this *zapEntry) Panicf(format string, args ...any) {
	this.log(LevelPanic, fmt.Sprintf(format, args...))
}

func (this *zapEntry) Fatal(args ...any) {
	this.log(LevelFatal, fmt.Sprint(args...))
}

func (this *zapEntry) Fatalf(format string, args ...any) {
	this.log(LevelFatal, fmt.Sprintf(format, args...))
}

func (this *zapEntry) Error(args ...any) {
	this.log(LevelError, fmt.Sprint(args...))
}

func (this *zapEntry) Errorf(format string, args ...any) {
	this.log(LevelError, fmt.Sprintf(format, args...))
}

func (this *zapEntry) Warn(args ...any) {
	this.log(LevelWarn, fmt.Sprint(args...))
}

func (this *zapEntry) Warnf(format string, args ...any) {
	this.log(LevelWarn, fmt.Sprintf(format, args...))
}

func (this *zapEntry) Info(args ...any) {
	this.log(LevelInfo, fmt.Sprint(args...))
}

func (this *zapEntry) Infof(format string, args ...any) {
	this.log(LevelInfo, fmt.Sprintf(format, args...))
}

func (this *zapEntry) Debug(args ...any) {
	this.log(LevelDebug, fmt.Sprint(args...))
}

func (this *zapEntry) Debugf(format string, args ...any) {
	this.log(LevelDebug, fmt.Sprintf(format, args...))
}

func (this *zapEntry) Trace(args ...any) {
	this.log(LevelTrace, fmt.Sprint(args...))
}

func (this *zapEntry) Tracef(format string, args ...any) {
	this.log(LevelTrace, fmt.Sprintf(format, args...))
}

type zapLogger struct {
	*leveledLogger

	writer io.Writer

	zl *zap.Logger
}

func newZapLogger(writer io.Writer, option *Option) (Logger, error) {
	var zl = &zapLogger{leveledLogger: &leveledLogger{level: ParseLevel(option.Level)}}

	zl.writer = writer

	var ec = zapcore.EncoderConfig{
		MessageKey:       "message",
		LevelKey:         "level",
		TimeKey:          "timestamp",
		NameKey:          "logger",
		CallerKey:        "caller",
		StacktraceKey:    "stacktrace",
		SkipLineEnding:   false,
		LineEnding:       "\r\n",
		EncodeLevel:      zapcore.LowercaseLevelEncoder,
		EncodeTime:       zapcore.TimeEncoderOfLayout(util.If(option.TimeFormat != "", option.TimeFormat, time.RFC3339)),
		EncodeDuration:   zapcore.StringDurationEncoder,
		EncodeCaller:     zapcore.FullCallerEncoder,
		EncodeName:       zapcore.FullNameEncoder,
		ConsoleSeparator: "",
	}

	zl.zl = zap.New(zapcore.NewCore(zapcore.NewJSONEncoder(ec), writer.(zapcore.WriteSyncer), zl.level.zapLevel()))
	if option.WithCaller {
		zl.zl = zl.zl.WithOptions(zap.WithCaller(true))
		zl.zl = zl.zl.WithOptions(zap.AddCallerSkip(3))
	}

	return zl, nil
}

func (this *zapLogger) Close() error {
	if closer, ok := this.writer.(io.Closer); ok {
		return closer.Close()
	}
	return nil
}

func (this *zapLogger) entry() *zapEntry {
	return &zapEntry{zl: this.zl}
}

func (this *zapLogger) WithField(name string, value any) entry {
	var ze = this.entry()
	ze.WithField(name, value)
	return ze
}

func (this *zapLogger) Panic(args ...any) {
	this.entry().Panic(args...)
}

func (this *zapLogger) Panicf(format string, args ...any) {
	this.entry().Panicf(format, args...)
}

func (this *zapLogger) Fatal(args ...any) {
	this.entry().Fatal(args...)
}

func (this *zapLogger) Fatalf(format string, args ...any) {
	this.entry().Fatalf(format, args...)
}

func (this *zapLogger) Error(args ...any) {
	this.entry().Error(args...)
}

func (this *zapLogger) Errorf(format string, args ...any) {
	this.entry().Errorf(format, args...)
}

func (this *zapLogger) Warn(args ...any) {
	this.entry().Warn(args...)
}

func (this *zapLogger) Warnf(format string, args ...any) {
	this.entry().Warnf(format, args...)
}

func (this *zapLogger) Info(args ...any) {
	this.entry().Info(args...)
}

func (this *zapLogger) Infof(format string, args ...any) {
	this.entry().Infof(format, args...)
}

func (this *zapLogger) Debug(args ...any) {
	this.entry().Debug(args...)
}

func (this *zapLogger) Debugf(format string, args ...any) {
	this.entry().Debugf(format, args...)
}

func (this *zapLogger) Trace(args ...any) {
	this.entry().Trace(args...)
}

func (this *zapLogger) Tracef(format string, args ...any) {
	this.entry().Tracef(format, args...)
}
