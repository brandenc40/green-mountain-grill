package app

import (
	"github.com/brandenc40/green-mountain-grill/server/config"
	"github.com/brandenc40/green-mountain-grill/server/grillclient"
	"github.com/brandenc40/green-mountain-grill/server/handler"
	"github.com/brandenc40/green-mountain-grill/server/logger"
	"github.com/brandenc40/green-mountain-grill/server/respository"
	"github.com/brandenc40/green-mountain-grill/server/server"
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
