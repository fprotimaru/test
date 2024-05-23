package config

import "github.com/kelseyhightower/envconfig"

type Config struct {
	DatabaseURL string `envconfig:"DATABASE_URL"`
	GRPCAddress string `envconfig:"GRPC_ADDRESS"`
}

func New() (*Config, error) {
	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
