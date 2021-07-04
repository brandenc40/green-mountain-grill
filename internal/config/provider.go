package config

import (
	"go.uber.org/config"
	"go.uber.org/fx"
)

// Module -
var Module = fx.Provide(NewProvider)

func NewProvider() (config.Provider, error) {
	return config.NewYAML(config.File("internal/config/config.yaml"))
}
