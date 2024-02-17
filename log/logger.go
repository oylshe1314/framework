package log

import (
	"bytes"
	"fmt"
	"framework/util"
	frl "github.com/lestrrat-go/file-rotatelogs"
	"github.com/sirupsen/logrus"
	"io"
	nativeLog "log"
	"strings"
)

type Level uint32

const (
	LevelPanic = Level(logrus.PanicLevel)
	LevelFatal = Level(logrus.FatalLevel)
	LevelError = Level(logrus.ErrorLevel)
	LevelWarn  = Level(logrus.WarnLevel)
	LevelInfo  = Level(logrus.InfoLevel)
	LevelDebug = Level(logrus.DebugLevel)
	LevelTrace = Level(logrus.TraceLevel)
)

func LevelOf(level string) Level {
	switch {
	case strings.EqualFold(level, "PANIC"):
		return LevelPanic
	case strings.EqualFold(level, "FATAL"):
		return LevelFatal
	case strings.EqualFold(level, "ERROR"):
		return LevelError
	case strings.EqualFold(level, "WARN"):
		return LevelWarn
	case strings.EqualFold(level, "WARNING"):
		return LevelWarn
	case strings.EqualFold(level, "INFO"):
		return LevelInfo
	case strings.EqualFold(level, "DEBUG"):
		return LevelDebug
	}
	return 0xFFFFFFFF
}

type Logger interface {
	io.Closer
	logrus.FieldLogger

	IsDebugEnabled() bool
	IsInfoEnabled() bool
	IsWarnEnabled() bool
	IsErrorEnabled() bool
	IsFatalEnabled() bool
	IsPanicEnabled() bool
}

type logFormatter struct {
}

func (this *logFormatter) relativePath(file string) string {
	var p = strings.Index(file, "ecs/")
	if p >= 0 {
		return file[p+4:]
	}
	return file
}

func (this *logFormatter) Format(entry *logrus.Entry) ([]byte, error) {

	var buffer = entry.Buffer
	if buffer == nil {
		buffer = &bytes.Buffer{}
	}

	var strLv = strings.ToUpper(entry.Level.String())
	switch strLv {
	case "WARNING":
		strLv = "WARN"
	case "UNKNOWN":
		strLv = "INFO"
	}

	buffer.WriteString("[")
	buffer.WriteString(entry.Time.Format("2006-01-02 15:04:05"))
	buffer.WriteString(" ")
	buffer.WriteString(fmt.Sprintf("%5s", strLv))
	buffer.WriteString("] ")
	buffer.WriteString(fmt.Sprintf("%-47s # ", fmt.Sprintf("%s:%d", this.relativePath(entry.Caller.File), entry.Caller.Line)))
	buffer.WriteString(entry.Message)
	if buffer.Bytes()[buffer.Len()-1] != '\n' {
		buffer.WriteByte('\n')
	}

	return buffer.Bytes(), nil
}

type Option interface {
	Name() string
	Value() interface{}
}

type logOption struct {
	name  string
	value interface{}
}

func (this *logOption) Name() string {
	return this.name
}
func (this *logOption) Value() interface{} {
	return this.value
}

func WithLevel(level Level) Option {
	return &logOption{name: "WithLevel", value: level}
}

func WithConsole(withConsole bool) Option {
	return &logOption{name: "WithConsole", value: withConsole}
}

type dailyLogger struct {
	*logrus.Logger
	rl *frl.RotateLogs
}

func NewDailyLogger(appName string, appId uint32, logDir string, opts ...Option) (Logger, error) {
	if logDir[len(logDir)-1] == '/' || logDir[len(logDir)-1] == '\\' {
		logDir = logDir[:len(logDir)-1]
	}

	var lv = LevelInfo
	var withConsole = false

	for _, opt := range opts {
		switch opt.Name() {
		case "WithLevel":
			lv = opt.Value().(Level)
		case "WithConsole":
			withConsole = opt.Value().(bool)
		}
	}

	rl, err := frl.New(fmt.Sprintf("%s/%s_%d_%%Y-%%m-%%d.log", logDir, appName, appId), frl.WithLocation(util.UTC8()))
	if err != nil {
		return nil, err
	}

	l := logrus.New()

	l.SetOutput(rl)
	l.SetReportCaller(true)
	l.SetLevel(logrus.Level(lv))
	l.SetFormatter(&logFormatter{})
	if withConsole {
		l.AddHook(&consoleHook{})
	}

	return &dailyLogger{Logger: l, rl: rl}, nil
}

func (this *dailyLogger) Close() error {
	return this.rl.Close()
}

func (this *dailyLogger) IsDebugEnabled() bool {
	return this.Logger.IsLevelEnabled(logrus.Level(LevelDebug))
}

func (this *dailyLogger) IsInfoEnabled() bool {
	return this.Logger.IsLevelEnabled(logrus.Level(LevelInfo))
}

func (this *dailyLogger) IsWarnEnabled() bool {
	return this.Logger.IsLevelEnabled(logrus.Level(LevelWarn))
}

func (this *dailyLogger) IsErrorEnabled() bool {
	return this.Logger.IsLevelEnabled(logrus.Level(LevelError))
}

func (this *dailyLogger) IsFatalEnabled() bool {
	return this.Logger.IsLevelEnabled(logrus.Level(LevelFatal))
}

func (this *dailyLogger) IsPanicEnabled() bool {
	return this.Logger.IsLevelEnabled(logrus.Level(LevelPanic))
}

type nativeLogWriter struct {
	level  Level
	logger Logger
}

func (this *nativeLogWriter) Write(buf []byte) (int, error) {
	switch this.level {
	case LevelPanic:
		this.logger.Panic(string(buf))
	case LevelFatal:
		this.logger.Fatal(string(buf))
	case LevelError:
		this.logger.Error(string(buf))
	case LevelWarn:
		this.logger.Warn(string(buf))
	case LevelInfo:
		this.logger.Info(string(buf))
	case LevelDebug:
		this.logger.Debug(string(buf))
	default:
		this.logger.Info(string(buf))
	}
	return len(buf), nil
}

func NewNativeLogger(logger Logger, level Level) *nativeLog.Logger {
	var writer = &nativeLogWriter{logger: logger, level: level}
	return nativeLog.New(writer, "NativeLog: ", 0)
}
