package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildVault(t *testing.T) {
	t.Run("default values", func(t *testing.T) {
		os.Clearenv() // Clear all environment variables

		cfg := &Config{}
		err := BuildVault(cfg)

		assert.NoError(t, err)
		assert.Equal(t, "localhost", cfg.Vault.Host)
		assert.Equal(t, "root", cfg.Vault.Token)
		assert.Equal(t, "", cfg.Vault.Port)
		assert.Equal(t, "https://localhost", cfg.Vault.Address)
	})

	t.Run("with port", func(t *testing.T) {
		os.Clearenv()
		_ = os.Setenv("VAULT_PORT", "8080")

		cfg := &Config{}
		err := BuildVault(cfg)

		assert.NoError(t, err)
		assert.Equal(t, "localhost:8080", cfg.Vault.Address)
	})

	t.Run("with http prefix", func(t *testing.T) {
		os.Clearenv()
		_ = os.Setenv("VAULT_HOST", "http://localhost")

		cfg := &Config{}
		err := BuildVault(cfg)

		assert.NoError(t, err)
		assert.Equal(t, "http://localhost", cfg.Vault.Address)
	})
}
