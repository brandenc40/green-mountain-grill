package app

import (
	"github.com/brandenc40/green-mountain-grill/internal/config"
	"github.com/brandenc40/green-mountain-grill/internal/grillclient"
	"github.com/brandenc40/green-mountain-grill/internal/handler"
	"github.com/brandenc40/green-mountain-grill/internal/logger"
	"github.com/brandenc40/green-mountain-grill/internal/respository"
	"github.com/brandenc40/green-mountain-grill/internal/server"
	"go.uber.org/fx"
)

func App() *fx.App {
	return fx.New(
		config.Module,
		logger.Module,
		respository.Module,
		handler.Module,
		server.Module,
		grillclient.Module,
		fx.Invoke(server.RegisterRoutes),
	)
}