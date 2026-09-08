// Package config loads and validates configuration from the environment.
package config

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/caarlos0/env/v11"
)

// Config is the complete runtime configuration. Every field comes from the environment.
type Config struct {
	Env             string        `env:"ENV" envDefault:"development"`
	Port            int           `env:"PORT" envDefault:"8080"`
	LogLevel        string        `env:"LOG_LEVEL" envDefault:"info"`
	DatabaseURL     string        `env:"DATABASE_URL,required"`
	APIKeys         string        `env:"API_KEYS"`
	CORSOrigins     []string      `env:"CORS_ORIGINS" envSeparator:","`
	TrustedProxies  []string      `env:"TRUSTED_PROXIES" envSeparator:","`
	RateLimitPerMin int           `env:"RATE_LIMIT_PER_MINUTE" envDefault:"120"`
	ShutdownTimeout time.Duration `env:"SHUTDOWN_TIMEOUT" envDefault:"10s"`
	MigrateOnStart  bool          `env:"MIGRATE_ON_START" envDefault:"false"`
	OTLPEndpoint    string        `env:"OTEL_EXPORTER_OTLP_ENDPOINT"`
}

// Load parses the environment and validates the result, returning every problem at once.
func Load() (Config, error) {
	var c Config
	if err := env.Parse(&c); err != nil {
		return c, fmt.Errorf("invalid environment: %w (see .env.example)", err)
	}
	var problems []error
	switch c.Env {
	case "development", "test", "production":
	default:
		problems = append(problems, fmt.Errorf("ENV must be development|test|production, got %q", c.Env))
	}
	if c.Port < 1 || c.Port > 65535 {
		problems = append(problems, fmt.Errorf("PORT must be 1-65535, got %d", c.Port))
	}
	if !strings.HasPrefix(c.DatabaseURL, "postgres://") && !strings.HasPrefix(c.DatabaseURL, "postgresql://") {
		problems = append(problems, errors.New("DATABASE_URL must start with postgres://"))
	}
	if c.RateLimitPerMin < 1 {
		problems = append(problems, errors.New("RATE_LIMIT_PER_MINUTE must be positive"))
	}
	if len(problems) > 0 {
		return c, fmt.Errorf("invalid environment (see .env.example):\n  - %w", errors.Join(problems...))
	}
	return c, nil
}

// APIKeyMap parses "key:principal,key2" into key → principal (principal defaults to the key).
func (c Config) APIKeyMap() map[string]string {
	out := map[string]string{}
	for pair := range strings.SplitSeq(c.APIKeys, ",") {
		pair = strings.TrimSpace(pair)
		if pair == "" {
			continue
		}
		key, principal, ok := strings.Cut(pair, ":")
		if !ok || principal == "" {
			principal = key
		}
		out[key] = principal
	}
	return out
}
