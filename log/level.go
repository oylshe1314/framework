package log

import (
	"strings"

	"github.com/rs/zerolog"
	"github.com/sirupsen/logrus"
	"go.uber.org/zap/zapcore"
)

type Level int

const (
	LevelNone Level = iota
	LevelPanic
	LevelFatal
	LevelError
	LevelWarn
	LevelInfo
	LevelDebug
	LevelTrace
)

const (
	PanicString = "PANIC"
	FatalString = "FATAL"
	ErrorString = "ERROR"
	WarnString  = "WARN"
	InfoString  = "INFO"
	DebugString = "DEBUG"
	TraceString = "TRACE"
)

func (level Level) String() string {
	switch level {
	case LevelPanic:
		return PanicString
	case LevelFatal:
		return FatalString
	case LevelError:
		return ErrorString
	case LevelWarn:
		return WarnString
	case LevelInfo:
		return InfoString
	case LevelDebug:
		return DebugString
	case LevelTrace:
		return TraceString
	default:
		return "UNKNOWN"
	}
}

func ParseLevel(level string) Level {
	switch strings.ToUpper(level) {
	case PanicString:
		return LevelPanic
	case FatalString:
		return LevelFatal
	case ErrorString:
		return LevelError
	case WarnString:
		return LevelWarn
	case InfoString:
		return LevelInfo
	case DebugString:
		return LevelDebug
	case TraceString:
		return LevelTrace
	default:
		return LevelInfo
	}
}

func (level Level) logrusLevel() logrus.Level {
	switch level {
	case LevelPanic:
		return logrus.PanicLevel
	case LevelFatal:
		return logrus.FatalLevel
	case LevelError:
		return logrus.ErrorLevel
	case LevelWarn:
		return logrus.WarnLevel
	case LevelInfo:
		return logrus.InfoLevel
	case LevelDebug:
		return logrus.DebugLevel
	case LevelTrace:
		return logrus.TraceLevel
	default:
		return logrus.InfoLevel
	}
}

func (level Level) zeroLevel() zerolog.Level {
	switch level {
	case LevelPanic:
		return zerolog.PanicLevel
	case LevelFatal:
		return zerolog.FatalLevel
	case LevelError:
		return zerolog.ErrorLevel
	case LevelWarn:
		return zerolog.WarnLevel
	case LevelInfo:
		return zerolog.InfoLevel
	case LevelDebug:
		return zerolog.DebugLevel
	case LevelTrace:
		return zerolog.TraceLevel
	default:
		return zerolog.InfoLevel
	}
}

func (level Level) zapLevel() zapcore.Level {
	switch level {
	case LevelPanic:
		return zapcore.PanicLevel
	case LevelFatal:
		return zapcore.FatalLevel
	case LevelError:
		return zapcore.ErrorLevel
	case LevelWarn:
		return zapcore.WarnLevel
	case LevelInfo:
		return zapcore.InfoLevel
	case LevelDebug:
		return zapcore.DebugLevel
	case LevelTrace:
		return zapcore.DebugLevel
	default:
		return zapcore.InfoLevel
	}
}

type leveledLogger struct {
	level Level
}

func (this *leveledLogger) IsLevelEnabled(level Level) bool {
	return this.level >= level
}

func (this *leveledLogger) IsPanicEnabled() bool {
	return this.IsLevelEnabled(LevelPanic)
}

func (this *leveledLogger) IsFatalEnabled() bool {
	return this.IsLevelEnabled(LevelFatal)
}

func (this *leveledLogger) IsErrorEnabled() bool {
	return this.IsLevelEnabled(LevelError)
}

func (this *leveledLogger) IsWarnEnabled() bool {
	return this.IsLevelEnabled(LevelWarn)
}

func (this *leveledLogger) IsInfoEnabled() bool {
	return this.IsLevelEnabled(LevelInfo)
}

func (this *leveledLogger) IsDebugEnabled() bool {
	return this.IsLevelEnabled(LevelDebug)
}

func (this *leveledLogger) IsTraceEnabled() bool {
	return this.IsLevelEnabled(LevelTrace)
}
