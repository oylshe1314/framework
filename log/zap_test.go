package log

import (
	"os"
	"testing"
	"time"
)

func TestNewZapLogger(t *testing.T) {
	var logger, err = newZapLogger(os.Stdout, &Option{
		Level:         TraceString,
		WithTimestamp: true,
		TimeFormat:    time.DateTime,
		WithCaller:    true,
	})
	if err != nil {
		t.Error(err)
		return
	}

	logger.WithField("name", "Lisa").WithField("age", 18).WithField("sexy", true).Error("error")
	logger.WithField("name", "Lisa").WithField("age", 18).WithField("sexy", true).Warn("warn")
	logger.WithField("name", "Lisa").WithField("age", 18).WithField("sexy", true).Info("info")
	logger.WithField("name", "Lisa").WithField("age", 18).WithField("sexy", true).Debug("debug")
	logger.WithField("name", "Lisa").WithField("age", 18).WithField("sexy", true).Trace("trace")
}
