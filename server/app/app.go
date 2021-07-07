package app

import (
	"github.com/brandenc40/green-mountain-grill/server/config"
	"github.com/brandenc40/green-mountain-grill/server/grillclient"
	"github.com/brandenc40/green-mountain-grill/server/handler"
	"github.com/brandenc40/green-mountain-grill/server/logger"
	"github.com/brandenc40/green-mountain-grill/server/respository"
	"github.com/brandenc40/green-mountain-grill/server/server"
	"github.com/sirupsen/logrus"
	"go.uber.org/fx"
)

func App() *fx.App {
	return fx.New(
		config.Module,
		grillclient.Module,
		handler.Module,
		logger.Module,
		respository.Module,
		server.Module,
		fx.Logger(logrus.StandardLogger()),
		fx.Invoke(server.RegisterRoutes),
	)
}
