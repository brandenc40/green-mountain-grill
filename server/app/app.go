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

var FXModule = fx.Options(
	// setup fx logger
	fx.Logger(logrus.StandardLogger()),

	// build dependencies
	config.Module,
	grillclient.Module,
	handler.Module,
	logger.Module,
	respository.Module,
	server.Module,

	// invoke route registration for the server on startup
	fx.Invoke(server.RegisterRoutes),
)

func App() *fx.App {
	return fx.New(FXModule)
}
