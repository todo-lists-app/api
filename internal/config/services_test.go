package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildServices(t *testing.T) {
	t.Run("default values", func(t *testing.T) {
		os.Clearenv() // Clear all environment variables

		cfg := &Config{}
		err := BuildServices(cfg)

		assert.NoError(t, err)
		assert.Equal(t, "id-checker.todo-list:3000", cfg.Services.Identity)
		assert.Equal(t, "todo-service.todo-list:3000", cfg.Services.Todo)
		assert.Equal(t, "user-service.todo-list:3000", cfg.Services.User)
	})

	t.Run("custom values", func(t *testing.T) {
		os.Clearenv()
		_ = os.Setenv("IDENTITY_SERVICE", "custom-id-service:4000")
		_ = os.Setenv("TODO_SERVICE", "custom-todo-service:4000")
		_ = os.Setenv("USER_SERVICE", "custom-user-service:4000")

		cfg := &Config{}
		err := BuildServices(cfg)

		assert.NoError(t, err)
		assert.Equal(t, "custom-id-service:4000", cfg.Services.Identity)
		assert.Equal(t, "custom-todo-service:4000", cfg.Services.Todo)
		assert.Equal(t, "custom-user-service:4000", cfg.Services.User)
	})

	// ... Add more test cases as needed
}
