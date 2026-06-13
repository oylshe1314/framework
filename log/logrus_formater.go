package log

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"
)

type logrusJsonFormatter struct {
	option *Option
}

func (this *logrusJsonFormatter) relativePath(file string) string {
	var p = strings.Index(file, "ecs/")
	if p >= 0 {
		return file[p+4:]
	}
	return file
}

func (this *logrusJsonFormatter) Format(entry *logrus.Entry) ([]byte, error) {
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
