// Package config is used to build the configuration for the service.
package config

import (
	"github.com/bugfixes/go-bugfixes/logs"
	"github.com/caarlos0/env/v8"
)

// Config is the main config
type Config struct {
	Local
	Mongo
	Vault
	Services
	Notifications
}

// Build is used to build the config, it will call BuildVault and BuildMongo
func Build() (*Config, error) {
	cfg := &Config{}

	if err := BuildVault(cfg); err != nil {
		return nil, logs.Errorf("build vault: %w", err)
	}

	if err := BuildMongo(cfg); err != nil {
		return nil, logs.Errorf("build mongo: %w", err)
	}

	if err := BuildServices(cfg); err != nil {
		return nil, logs.Errorf("build services: %w", err)
	}

	if err := BuildNotifications(cfg); err != nil {
		return nil, logs.Errorf("build notifications: %w", err)
	}

	if err := BuildLocal(cfg); err != nil {
		return nil, logs.Errorf("build local: %w", err)
	}

	if err := env.Parse(cfg); err != nil {
		return nil, logs.Errorf("parse config: %w", err)
	}

	return cfg, nil
}
