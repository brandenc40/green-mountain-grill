package grillclient

import "go.uber.org/config"

type Config struct {
	GrillIP   string `yaml:"grill_ip"`
	GrillPort int    `yaml:"grill_port"`
}

func NewConfig(provider config.Provider) (*Config, error) {
	var c Config
	if err := provider.Get("grill_client").Populate(&c); err != nil {
		return nil, err
	}
	return &c, nil
}
