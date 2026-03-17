package config

import (
	"crypto/rand"
	"encoding/hex"
	"os"
	"strconv"
)

// Config holds all application configuration.
type Config struct {
	Port         int
	DatabasePath string
	JWTSecret    string
	LogLevel     string
}

// Load reads configuration from environment variables with sensible defaults.
func Load() *Config {
	cfg := &Config{
		Port:         envInt("PORT", 8080),
		DatabasePath: envStr("DATABASE_PATH", "presence.db"),
		JWTSecret:    envStr("JWT_SECRET", ""),
		LogLevel:     envStr("LOG_LEVEL", "info"),
	}

	if cfg.JWTSecret == "" {
		cfg.JWTSecret = generateSecret()
	}

	return cfg
}

func envStr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func envInt(key string, fallback int) int {
	if v := os.Getenv(key); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	return fallback
}

func generateSecret() string {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		panic("failed to generate JWT secret: " + err.Error())
	}
	return hex.EncodeToString(b)
}
