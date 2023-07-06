package config

import (
	"github.com/bugfixes/go-bugfixes/logs"
	env "github.com/caarlos0/env/v8"
)

// Local is the local config
type Local struct {
	KeepLocal   bool `env:"BUGFIXES_LOCAL_ONLY" envDefault:"false"`
	Development bool `env:"DEVELOPMENT" envDefault:"false"`
	HTTPPort    int  `env:"HTTP_PORT" envDefault:"80"`
}

// BuildLocal builds the local config
func BuildLocal(cfg *Config) error {
	local := &Local{}
	if err := env.Parse(local); err != nil {
		return logs.Errorf("failed to parse local config: %v", err)
	}
	cfg.Local = *local

	return nil
}
