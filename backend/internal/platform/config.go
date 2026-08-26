package platform

import (
	"errors"
	"fmt"
	"net"
	"os"
	"strings"
)

// Environment represents the runtime environment
type Environment string

const (
	// EnvDevelopment represents the development environment
	EnvDevelopment Environment = "development"
	// EnvTest represents the test environment
	EnvTest Environment = "test"
	// EnvProduction represents the production environment
	EnvProduction Environment = "production"
)

// Config holds the application configuration
type Config struct {
	ServerAddr  string
	DatabaseURL string
	Environment Environment
}

// LoadConfig loads configuration from environment variables with sensible defaults
func LoadConfig() (*Config, error) {
	cfg := &Config{
		ServerAddr:  ":8080",
		Environment: EnvDevelopment,
	}

	if addr := os.Getenv("AEGIS_SERVER_ADDR"); addr != "" {
		cfg.ServerAddr = addr
	}

	if env := os.Getenv("AEGIS_ENV"); env != "" {
		cfg.Environment = Environment(strings.ToLower(env))
	}

	cfg.DatabaseURL = os.Getenv("AEGIS_DATABASE_URL")

	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

// Validate checks if the configuration values are valid
func (c *Config) Validate() error {
	if c.DatabaseURL == "" {
		return errors.New("AEGIS_DATABASE_URL is required")
	}

	switch c.Environment {
	case EnvDevelopment, EnvTest, EnvProduction:
		// valid environment
	default:
		return fmt.Errorf("invalid environment: %s", c.Environment)
	}

	_, port, err := net.SplitHostPort(c.ServerAddr)
	if err != nil {
		return fmt.Errorf("invalid server address format: %w", err)
	}
	if port == "" {
		return errors.New("server address must include a port")
	}

	return nil
}

// String returns a string representation of the config, with secrets redacted
func (c *Config) String() string {
	return fmt.Sprintf("Config{ServerAddr: %s, Environment: %s, DatabaseURL: [REDACTED]}", c.ServerAddr, c.Environment)
}
