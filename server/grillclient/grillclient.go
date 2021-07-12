package grillclient

import (
	"net"

	"go.uber.org/zap"

	gmg "github.com/brandenc40/green-mountain-grill"
	"go.uber.org/fx"
)

type Params struct {
	fx.In

	Config *Config
	Logger *zap.Logger
}

func New(p Params) (gmg.Client, error) {
	params := gmg.Params{
		GrillIP:   net.ParseIP(p.Config.GrillIP),
		GrillPort: p.Config.GrillPort,
		Logger:    p.Logger,
	}
	return gmg.New(params)
}
