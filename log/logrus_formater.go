package log

import (
	"bytes"
	"fmt"
	"runtime"
	"time"

	"github.com/oylshe1314/framework/store"
	"github.com/sirupsen/logrus"
)

type logrusJsonFormatter struct {
	option *Option
}

func (this *logrusJsonFormatter) writeJson(buffer *bytes.Buffer, fields []*store.Pair[string, any]) {
	buffer.WriteString("{")
	buffer.WriteString("\"")
	buffer.WriteString(fields[0].Key)
	buffer.WriteString("\":")
	switch fields[0].Value.(type) {
	case string:
		buffer.WriteString("\"")
		buffer.WriteString(fields[0].Value.(string))
		buffer.WriteString("\"")
	default:
		buffer.WriteString(fmt.Sprint(fields[0].Value))
	}
	for _, field := range fields[1:] {
		buffer.WriteString(",")
		buffer.WriteString("\"")
		buffer.WriteString(field.Key)
		buffer.WriteString("\":")
		switch field.Value.(type) {
		case string:
			buffer.WriteString("\"")
			buffer.WriteString(field.Value.(string))
			buffer.WriteString("\"")
		default:
			buffer.WriteString(fmt.Sprint(field.Value))
		}
	}
	buffer.WriteString("}\r\n")
}

func (this *logrusJsonFormatter) getCaller() *runtime.Frame {
	pcs := make([]uintptr, 32)
	depth := runtime.Callers(10, pcs)
	frames := runtime.CallersFrames(pcs[:depth])
	frame, _ := frames.Next()
	return &frame
}

func (this *logrusJsonFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	var buffer = entry.Buffer
	if buffer == nil {
		buffer = &bytes.Buffer{}
	}

	var caller = this.getCaller()

	var fields = []*store.Pair[string, any]{
		store.NewPair[string, any]("level", entry.Level.String()),
		store.NewPair[string, any]("message", entry.Message),
	}

	if this.option.WithTimestamp {
		var timeFormat = this.option.TimeFormat
		if timeFormat == "" {
			timeFormat = time.RFC3339
		}
		fields = append(fields, store.NewPair[string, any]("time", entry.Time.Format(timeFormat)))
	}
	if this.option.WithCaller {
		fields = append(fields, store.NewPair[string, any]("caller", fmt.Sprint(caller.File, ":", caller.Line)))
	}

	fields = append(fields, entry.Data["fields"].([]*store.Pair[string, any])...)

	this.writeJson(buffer, fields)

	return buffer.Bytes(), nil
}
