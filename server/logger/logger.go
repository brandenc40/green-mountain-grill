package logger

import (
	"fmt"
	"log"
	"os"

	"go.uber.org/fx"
	"go.uber.org/zap"
)

func BuildFxOptions() fx.Option {
	config := zap.NewDevelopmentConfig()
	if os.Getenv("ENVIRONMENT") == "production" {
		config.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	}
	logger, err := config.Build()
	if err != nil {
		log.Fatal(err)
	}
	return fx.Options(
		fx.Provide(func() *zap.Logger { return logger }),
		fx.Logger(fxLogger{logger}),
	)
}

// fxLogger implements the fx.Printer interface for use as the FX app logger
type fxLogger struct {
	*zap.Logger
}

func (l fxLogger) Printf(s string, args ...interface{}) {
	l.Info(fmt.Sprintf(s, args...))
}
