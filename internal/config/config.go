// Package config is used to build the configuration for the service.
package config

import (
	"github.com/bugfixes/go-bugfixes/logs"
	"github.com/caarlos0/env/v8"
)

// Config is the main config
type Config struct {
	Local
	Vault
	Services
}

// Build is used to build the config, it will call BuildVault and BuildMongo
func Build() (*Config, error) {
	cfg := &Config{}

	if err := BuildVault(cfg); err != nil {
		return nil, logs.Errorf("build vault: %v", err)
	}

	if err := BuildServices(cfg); err != nil {
		return nil, logs.Errorf("build services: %v", err)
	}

	if err := BuildLocal(cfg); err != nil {
		return nil, logs.Errorf("build local: %v", err)
	}

	if err := env.Parse(cfg); err != nil {
		return nil, logs.Errorf("parse config: %v", err)
	}

	return cfg, nil
}
