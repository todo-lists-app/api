// Package config is used to build the configuration for the service.
package config

import (
	"github.com/bugfixes/go-bugfixes/logs"
	env "github.com/caarlos0/env/v8"
)

// Config is the main config
type Config struct {
	Local
	Mongo
	Vault
	Identity
	Notifications
}

// Build is used to build the config, it will call BuildVault and BuildMongo
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

	if err := BuildIdentity(cfg); err != nil {
		return nil, logs.Error(err)
	}

	if err := BuildNotifications(cfg); err != nil {
		return nil, logs.Error(err)
	}

	if err := BuildLocal(cfg); err != nil {
		return nil, logs.Error(err)
	}

	return cfg, nil
}
