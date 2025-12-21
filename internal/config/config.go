package config

import (
	"fmt"
	"os"
)

type Config struct {
	HTTPPort    string
	DB_DSN      string
	MongoDB_URI string
}

func Load() (*Config, error) {
	cfg := &Config{
		HTTPPort:    getEnv("HTTP_PORT", "8000"),
		DB_DSN:      getEnv("DB_DSN", ""),
		MongoDB_URI: getEnv("MONGODB_URI", ""),
	}

	if cfg.DB_DSN == "" {
		return nil, fmt.Errorf("missing DB_DSN environment variable")
	}

	if cfg.MongoDB_URI == "" {
		return nil, fmt.Errorf("missing MONGODB_URI environment variable")
	}

	return cfg, nil
}

func getEnv(key string, defaultVal string) string {
	if val, ok := os.LookupEnv(key); ok && val != "" {
		return val
	}
	return defaultVal
}
