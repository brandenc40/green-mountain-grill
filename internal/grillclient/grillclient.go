package grillclient

import (
	"net"

	"github.com/brandenc40/green-mountain-grill/client"
	"github.com/sirupsen/logrus"
	"go.uber.org/fx"
)

var Module = fx.Provide(New, NewConfig)

type Params struct {
	fx.In

	Config *Config
	Logger *logrus.Logger
}

func New(p Params) client.Client {
	params := client.Params{
		GrillIP:   net.ParseIP(p.Config.GrillIP),
		GrillPort: p.Config.GrillPort,
		Logger:    p.Logger,
	}
	return client.New(params)
}
