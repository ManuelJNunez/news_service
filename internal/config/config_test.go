package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadSuccess(t *testing.T) {
	t.Setenv("HTTP_PORT", "9090")
	t.Setenv("DB_DSN", "postgres://user:pass@localhost:5432/db?sslmode=disable")
	t.Setenv("MONGODB_URI", "mongodb://fake_host:27017")

	cfg, err := Load()

	assert.NoError(t, err)
	assert.Equal(t, "9090", cfg.HTTPPort)
	assert.NotEmpty(t, cfg.DB_DSN)
	assert.NotEmpty(t, cfg.MongoDB_URI)
}

func TestLoadUsesDefaultPort(t *testing.T) {
	t.Setenv("HTTP_PORT", "")
	_ = os.Unsetenv("HTTP_PORT") // make sure it is unset
	t.Setenv("DB_DSN", "dsn")
	t.Setenv("MONGODB_URI", "mongodb://fake_host:27017")

	cfg, err := Load()

	assert.NoError(t, err)
	assert.Equal(t, "8000", cfg.HTTPPort)
}

func TestLoadMissingDSN(t *testing.T) {
	t.Setenv("HTTP_PORT", "8081")
	t.Setenv("DB_DSN", "")

	_, err := Load()

	assert.Error(t, err)
}
