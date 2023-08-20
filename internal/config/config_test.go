package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuild(t *testing.T) {
	t.Run("default values", func(t *testing.T) {
		os.Clearenv() // Clear all environment variables

		cfg, err := Build()

		assert.NoError(t, err)
		assert.Equal(t, false, cfg.Local.KeepLocal)
		assert.Equal(t, false, cfg.Local.Development)
		assert.Equal(t, 80, cfg.Local.HTTPPort)
		assert.Equal(t, "https://localhost", cfg.Vault.Address)
		assert.Equal(t, "todo-service.todo-list:3000", cfg.Services.Todo)
	})

	t.Run("custom values", func(t *testing.T) {
		os.Clearenv()
		os.Setenv("BUGFIXES_LOCAL_ONLY", "true")
		os.Setenv("DEVELOPMENT", "true")
		os.Setenv("HTTP_PORT", "8080")
		os.Setenv("VAULT_HOST", "https://vault.example.com")
		os.Setenv("TODO_SERVICE", "todo-service.example.com")

		cfg, err := Build()

		assert.NoError(t, err)
		assert.Equal(t, true, cfg.Local.KeepLocal)
		assert.Equal(t, true, cfg.Local.Development)
		assert.Equal(t, 8080, cfg.Local.HTTPPort)
		assert.Equal(t, "https://vault.example.com", cfg.Vault.Address)
		assert.Equal(t, "todo-service.example.com", cfg.Services.Todo)
	})
}
