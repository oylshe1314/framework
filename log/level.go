package log

import (
	"strings"
)

type Level int

const (
	LevelPanic Level = iota
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
