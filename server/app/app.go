package app

import (
	"github.com/brandenc40/green-mountain-grill/server/config"
	"github.com/brandenc40/green-mountain-grill/server/grillclient"
	"github.com/brandenc40/green-mountain-grill/server/handler"
	"github.com/brandenc40/green-mountain-grill/server/logger"
	"github.com/brandenc40/green-mountain-grill/server/respository"
	"github.com/brandenc40/green-mountain-grill/server/scheduler"
	"github.com/brandenc40/green-mountain-grill/server/server"
	"go.uber.org/fx"
)

var Options = fx.Options(
	// build single logger object to be used by fx and dep injection
	logger.BuildFxOptions(),

	// build dependency modules
	config.Module,
	grillclient.Module,
	handler.Module,
	respository.Module,
	scheduler.Module,
	server.Module,

	// invoke route registration for the server on startup
	fx.Invoke(server.RegisterRoutes),
)

func App() *fx.App {
	return fx.New(Options)
}
