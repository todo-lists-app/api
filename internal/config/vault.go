package config

import (
	"fmt"
	"github.com/bugfixes/go-bugfixes/logs"
	"strings"

	"github.com/caarlos0/env/v8"
)

// Vault is the vault config
type Vault struct {
	Host    string `env:"VAULT_HOST" envDefault:"localhost"`
	Port    string `env:"VAULT_PORT" envDefault:""`
	Token   string `env:"VAULT_TOKEN" envDefault:"root"`
	Address string
}

// BuildVault builds the vault config
func BuildVault(cfg *Config) error {
	v := &Vault{}

	if err := env.Parse(v); err != nil {
		return logs.Errorf("vault: %v", err)
	}

	if strings.HasPrefix(v.Host, "http") {
		v.Address = v.Host
	}

	if v.Port != "" {
		v.Address = fmt.Sprintf("%s:%s", v.Host, v.Port)
	}

	if v.Address == "" {
		v.Address = fmt.Sprintf("https://%s", v.Host)
	}

	cfg.Vault = *v

	return nil
}
