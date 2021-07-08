package grillclient

import (
	"net"

	gmg "github.com/brandenc40/green-mountain-grill"
	"github.com/sirupsen/logrus"
	"go.uber.org/fx"
)

type Params struct {
	fx.In

	Config *Config
	Logger *logrus.Logger
}

func New(p Params) gmg.Client {
	params := gmg.Params{
		GrillIP:   net.ParseIP(p.Config.GrillIP),
		GrillPort: p.Config.GrillPort,
		Logger:    p.Logger,
	}
	return gmg.New(params)
}
