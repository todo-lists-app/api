package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildLocal(t *testing.T) {
	t.Run("default values", func(t *testing.T) {
		os.Clearenv() // Clear all environment variables

		cfg := &Config{}
		err := BuildLocal(cfg)

		assert.NoError(t, err)
		assert.Equal(t, false, cfg.Local.KeepLocal)
		assert.Equal(t, false, cfg.Local.Development)
		assert.Equal(t, 80, cfg.Local.HTTPPort)
	})

	t.Run("custom values", func(t *testing.T) {
		os.Clearenv()
		_ = os.Setenv("BUGFIXES_LOCAL_ONLY", "true")
		_ = os.Setenv("DEVELOPMENT", "true")
		_ = os.Setenv("HTTP_PORT", "8080")

		cfg := &Config{}
		err := BuildLocal(cfg)

		assert.NoError(t, err)
		assert.Equal(t, true, cfg.Local.KeepLocal)
		assert.Equal(t, true, cfg.Local.Development)
		assert.Equal(t, 8080, cfg.Local.HTTPPort)
	})
}
