// Package config is used to build the configuration for the service.
package config

import (
	"github.com/bugfixes/go-bugfixes/logs"
	gc "github.com/keloran/go-config"
)

// Config is the main config
type Config struct {
	Services
	gc.Config
}

// Build is used to build the config, it will call BuildVault and BuildMongo
func Build() (*Config, error) {
	cfg := &Config{}

	gcc, err := gc.Build(gc.Vault, gc.Local)
	if err != nil {
		return nil, logs.Errorf("build config: %v", err)
	}
	cfg.Config = *gcc

	if err := BuildServices(cfg); err != nil {
		return nil, logs.Errorf("build services: %v", err)
	}

	return cfg, nil
}
