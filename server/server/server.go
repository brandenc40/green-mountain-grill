package server

import (
	"context"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type Params struct {
	fx.In

	Config     *Config
	Logger     *zap.Logger
	Lifecycle  fx.Lifecycle
	Shutdowner fx.Shutdowner
}

type Server struct {
	*fiber.App
}

func New(p Params) *Server {
	app := fiber.New(fiber.Config{
		ReadTimeout: 3 * time.Second,
		AppName:     "green-mountain-grill-server",
	})
	app.Use(recover.New())
	if p.Logger.Core().Enabled(zap.DebugLevel) {
		app.Use(logger.New(logger.Config{
			TimeFormat: time.RFC822,
		}))
	}
	s := &Server{App: app}
	p.Lifecycle.Append(
		fx.Hook{
			OnStart: func(context.Context) error {
				go func() {
					err := s.Listen(p.Config.ServerPort)
					if err != nil {
						p.Logger.Error("unable to start server", zap.Error(err))
						if err := p.Shutdowner.Shutdown(); err != nil {
							p.Logger.Error("could not shutdown, exiting with os.Exit(1)", zap.Error(err))
							os.Exit(1)
						}
					}
				}()
				return nil
			},
			OnStop: func(ctx context.Context) error {
				p.Logger.Info("shutting down server")
				return s.Shutdown()
			},
		},
	)
	return s
}
