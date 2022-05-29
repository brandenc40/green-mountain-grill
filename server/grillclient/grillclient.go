package grillclient

import (
	"net"

	gmg "github.com/brandenc40/green-mountain-grill"
	"github.com/brandenc40/green-mountain-grill/mocks"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type Params struct {
	fx.In

	Config *Config
	Logger *zap.Logger
}

func New(p Params) (gmg.Client, error) {
	if p.Config.IsMock {
		return newMockClient(), nil
	}
	return gmg.New(
		net.ParseIP(p.Config.GrillIP),
		p.Config.GrillPort,
		gmg.WithZapLogger(p.Logger),
	)
}

func newMockClient() gmg.Client {
	client := mocks.Client{}
	client.On("GetState").Return(&gmg.State{
		CurrentTemperature:      101,
		TargetTemperature:       150,
		Probe1Temperature:       110,
		Probe1TargetTemperature: 200,
		FireState:               gmg.FireStateRunning,
	}, nil)
	return &client
}
