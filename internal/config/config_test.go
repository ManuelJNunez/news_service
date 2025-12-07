package config

import (
	"os"
	"testing"
)

func TestLoadSuccess(t *testing.T) {
	t.Setenv("HTTP_PORT", "9090")
	t.Setenv("DB_DSN", "postgres://user:pass@localhost:5432/db?sslmode=disable")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if cfg.HTTPPort != "9090" {
		t.Fatalf("expected HTTPPort=9090, got %s", cfg.HTTPPort)
	}
	if cfg.DB_DSN == "" {
		t.Fatalf("expected DB_DSN to be set")
	}
}

func TestLoadUsesDefaultPort(t *testing.T) {
	t.Setenv("HTTP_PORT", "")
	_ = os.Unsetenv("HTTP_PORT") // make sure it is unset
	t.Setenv("DB_DSN", "dsn")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if cfg.HTTPPort != "8000" {
		t.Fatalf("expected default port 8000, got %s", cfg.HTTPPort)
	}
}

func TestLoadMissingDSN(t *testing.T) {
	t.Setenv("HTTP_PORT", "8081")
	t.Setenv("DB_DSN", "")

	_, err := Load()
	if err == nil {
		t.Fatalf("expected error for missing DB_DSN, got nil")
	}
}
