package config

import "github.com/caarlos0/env/v8"

type Services struct {
	Identity string `env:"IDENTITY_SERVICE" envDefault:"id-checker.todo-list:3000"`
	Todo     string `env:"TODO_SERVICE" envDefault:"todo-service.todo-list:3000"`
}

func BuildServices(cfg *Config) error {
	services := &Services{}
	if err := env.Parse(services); err != nil {
		return err
	}
	cfg.Services = *services

	return nil
}
