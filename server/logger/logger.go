package logger

import (
	"fmt"
	"log"
	"os"

	"go.uber.org/fx"
	"go.uber.org/zap"
)

var (
	Module = fx.Provide(New)

	FxOption     fx.Option
	GlobalLogger *zap.Logger
)

func init() {
	config := zap.NewDevelopmentConfig()
	if os.Getenv("ENVIRONMENT") == "production" {
		config.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	}
	logger, err := config.Build()
	if err != nil {
		log.Fatal(err)
	}
	GlobalLogger = logger
	FxOption = fx.Logger(fxLogger{logger})
}

func New() *zap.Logger {
	return GlobalLogger
}

// fxLogger implements the fx.Printer interface for use as the FX app logger
type fxLogger struct {
	*zap.Logger
}

func (l fxLogger) Printf(s string, args ...interface{}) {
	l.Info(fmt.Sprintf(s, args...))
}
