package logger

import (
	"os"

	"go.uber.org/zap/zapcore"

	"go.uber.org/fx"
	"go.uber.org/zap"
)

var Module = fx.Provide(New)

func New() (*zap.Logger, error) {
	logger, err := zap.NewDevelopment()
	if err != nil {
		return nil, err
	}
	if os.Getenv("ENVIRONMENT") == "production" {
		logger = logger.WithOptions(zap.IncreaseLevel(zapcore.InfoLevel))
	}
	return logger, nil
}
