package config

import (
	"github.com/bugfixes/go-bugfixes/logs"
	env "github.com/caarlos0/env/v8"
)

type Config struct {
	Local
	Mongo
	Vault
}

func Build() (*Config, error) {
	cfg := &Config{}

	if err := env.Parse(cfg); err != nil {
		return nil, logs.Error(err)
	}

	if err := BuildVault(cfg); err != nil {
		return nil, logs.Error(err)
	}

	if err := BuildMongo(cfg); err != nil {
		return nil, logs.Error(err)
	}

	if err := BuildLocal(cfg); err != nil {
		return nil, logs.Error(err)
	}

	return cfg, nil
}
