package logger

import (
	"os"

	"github.com/sirupsen/logrus"
	"go.uber.org/fx"
)

var Module = fx.Provide(New)

func New() *logrus.Logger {
	level := logrus.DebugLevel
	if os.Getenv("ENVIRONMENT") == "production" {
		level = logrus.InfoLevel
	}
	logger := &logrus.Logger{
		Out:       os.Stderr,
		Formatter: new(logrus.TextFormatter),
		Hooks:     make(logrus.LevelHooks),
		Level:     level,
	}
	logger.Debug("logging at DEBUG level")
	return logger
}
