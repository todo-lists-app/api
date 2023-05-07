package config

import (
	env "github.com/caarlos0/env/v8"
)

type Local struct {
	KeepLocal   bool `env:"LOCAL_ONLY" envDefault:"false" json:"keep_local,omitempty"`
	Development bool `env:"DEVELOPMENT" envDefault:"false" json:"development,omitempty"`
	HTTPPort    int  `env:"HTTP_PORT" envDefault:"80" json:"port,omitempty"`
}

func BuildLocal(cfg *Config) error {
	local := &Local{}
	if err := env.Parse(local); err != nil {
		return err
	}
	cfg.Local = *local

	return nil
}