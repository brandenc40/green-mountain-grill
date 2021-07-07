package logger

import (
	"os"

	"github.com/sirupsen/logrus"
	"go.uber.org/fx"
)

var Module = fx.Provide(New)

func New() *logrus.Logger {
	logger := logrus.StandardLogger()
	if os.Getenv("ENVIRONMENT") != "production" {
		logger.SetLevel(logrus.DebugLevel)
	}
	return logger
}
