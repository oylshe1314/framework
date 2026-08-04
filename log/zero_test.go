package log

import (
	"os"
	"testing"
)

func TestNewZeroLogger(t *testing.T) {
	var logger, err = newZeroLogger(os.Stdout, &Option{
		WithTimestamp: true,
		TimeFormat:    "2006-01-02 15:04:05",
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
