package config

import (
	"fmt"
	env "github.com/caarlos0/env/v8"
	"strings"
)

type Vault struct {
	Host    string `env:"VAULT_HOST" envDefault:"localhost"`
	Port    string `env:"VAULT_PORT" envDefault:""`
	Token   string `env:"VAULT_TOKEN" envDefault:"root"`
	Address string `env:"VAULT_ADDRESS" envDefault:""`
}

func BuildVault(cfg *Config) error {
	v := &Vault{}

	if err := env.Parse(v); err != nil {
		return err
	}

	if strings.HasPrefix(v.Host, "http") {
		v.Address = v.Host
	}

	if v.Port != "" {
		v.Address = fmt.Sprintf("%s:%s", v.Host, v.Port)
	}

	cfg.Vault = *v

	return nil
}
