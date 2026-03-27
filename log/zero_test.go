package log

import (
	"os"
	"testing"
)

func TestNewZeroLogger(t *testing.T) {
	var logger = NewZeroLogger(LevelDebug, os.Stdout)

	logger.WithField("name", "Lisa").WithField("age", 18).WithField("sexy", true).Error("error")
	logger.WithField("name", "Lisa").WithField("age", 18).WithField("sexy", true).Warn("warn")
	logger.WithField("name", "Lisa").WithField("age", 18).WithField("sexy", true).Info("info")
	logger.WithField("name", "Lisa").WithField("age", 18).WithField("sexy", true).Debug("debug")
	logger.WithField("name", "Lisa").WithField("age", 18).WithField("sexy", true).Trace("trace")
}
