package config

import "github.com/caarlos0/env/v8"

type Identity struct {
	Service string `env:"IDENTITY_SERVICE" envDefault:"id-checker.todo-list:3000" json:"service,omitempty"`
}

func BuildIdentity(cfg *Config) error {
	identity := &Identity{}
	if err := env.Parse(identity); err != nil {
		return err
	}
	cfg.Identity = *identity

	return nil
}
