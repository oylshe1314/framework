package log

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
)

var colors = [...]string{"31m", "31m", "31m", "33m", "34m", "32m", "32m"}

type consoleHook struct{}

func (this *consoleHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (this *consoleHook) Fire(entry *logrus.Entry) error {
	line, err := entry.String()
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Unable to read entry, %v", err)
		return err
	}

	_, err = os.Stdout.WriteString(fmt.Sprintf("\x1b[%s", colors[entry.Level]))
	if err != nil {
		return err
	}

	_, err = os.Stdout.WriteString(line)
	if err != nil {
		return err
	}

	_, err = os.Stdout.WriteString("\x1b[0m")
	return err
}

type noneWriter struct {
}

func (noneWriter) Write(p []byte) (n int, err error) {
	return len(p), nil
}

type consoleLogger struct {
	*logrus.Logger
}

func (this *consoleLogger) Close() error {
	return nil
}

func (this *consoleLogger) IsDebugEnabled() bool {
	return this.Logger.IsLevelEnabled(logrus.Level(LevelDebug))
}

func (this *consoleLogger) IsInfoEnabled() bool {
	return this.Logger.IsLevelEnabled(logrus.Level(LevelInfo))
}

func (this *consoleLogger) IsWarnEnabled() bool {
	return this.Logger.IsLevelEnabled(logrus.Level(LevelWarn))
}

func (this *consoleLogger) IsErrorEnabled() bool {
	return this.Logger.IsLevelEnabled(logrus.Level(LevelError))
}

func (this *consoleLogger) IsFatalEnabled() bool {
	return this.Logger.IsLevelEnabled(logrus.Level(LevelFatal))
}

func (this *consoleLogger) IsPanicEnabled() bool {
	return this.Logger.IsLevelEnabled(logrus.Level(LevelPanic))
}

func newConsoleLogger() Logger {
	l := logrus.New()
	l.SetReportCaller(true)
	l.SetLevel(logrus.Level(LevelDebug))
	l.SetOutput(noneWriter{})
	l.SetFormatter(&logFormatter{})
	l.AddHook(&consoleHook{})
	return &consoleLogger{Logger: l}
}

var DefaultLogger = newConsoleLogger()
