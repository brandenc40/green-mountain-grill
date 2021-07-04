package grillclient

import (
	"net"

	"github.com/brandenc40/green-mountain-grill/grillclient"
	"github.com/sirupsen/logrus"
	"go.uber.org/fx"
)

var Module = fx.Provide(New, NewConfig)

type Params struct {
	fx.In

	Config *Config
	Logger *logrus.Logger
}

func New(p Params) grillclient.Client {
	params := grillclient.Params{
		GrillIP:   net.ParseIP(p.Config.GrillIP),
		GrillPort: p.Config.GrillPort,
		Logger:    p.Logger,
	}
	return grillclient.New(params)
}
