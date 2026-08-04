package log

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
)

var colors = [...]string{"31m", "31m", "31m", "33m", "34m", "32m", "32m"}

type logrusConsoleHook struct{}

func (this *logrusConsoleHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (this *logrusConsoleHook) Fire(entry *logrus.Entry) error {
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
