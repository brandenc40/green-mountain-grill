package server

import "go.uber.org/config"

type Config struct {
	ServerPort string `yaml:"server_port"`
}

func NewConfig(provider config.Provider) (*Config, error) {
	var c Config
	if err := provider.Get("server").Populate(&c); err != nil {
		return nil, err
	}
	return &c, nil
}
