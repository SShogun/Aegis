package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/caarlos0/env/v11"
)

type AppEnv string

const (
	AppEnvDevelopment AppEnv = "development"
	AppEnvProduction  AppEnv = "production"
	AppEnvTest        AppEnv = "test"
)

type RateLimitFailureMode string

const (
	RateLimitFailureOpen   RateLimitFailureMode = "open"
	RateLimitFailureClosed RateLimitFailureMode = "closed"
)

type Config struct {
	AppEnv               AppEnv               `env:"APP_ENV" envDefault:"development"`
	ServiceName          string               `env:"SERVICE_NAME" envDefault:"aegis"`
	HTTPAddr             string               `env:"HTTP_ADDR" envDefault:":8080"`
	GRPCAddr             string               `env:"GRPC_ADDR" envDefault:":9090"`
	MetricsAddr          string               `env:"METRICS_ADDR" envDefault:":9091"`
	RedisAddrs           []string             `env:"REDIS_ADDRS" envSeparator:"," envDefault:"localhost:6379"`
	RedisPassword        string               `env:"REDIS_PASSWORD"`
	RedisDB              int                  `env:"REDIS_DB" envDefault:"0"`
	PostgresDSN          string               `env:"POSTGRES_DSN"`
	JWTIssuer            string               `env:"JWT_ISSUER"`
	JWTAudience          string               `env:"JWT_AUDIENCE"`
	JWKSURL              string               `env:"JWKS_URL"`
	LocalJWKSPath        string               `env:"LOCAL_JWKS_PATH"`
	RateLimitFailureMode RateLimitFailureMode `env:"RATE_LIMIT_FAILURE_MODE" envDefault:"open"`

	OTELExporterOTLPEndpoint string `env:"OTEL_EXPORTER_OTLP_ENDPOINT"`
	OTELServiceName          string `env:"OTEL_SERVICE_NAME" envDefault:"aegis"`
	OTELServiceVersion       string `env:"OTEL_SERVICE_VERSION" envDefault:"dev"`

	ReadTimeout       time.Duration `env:"HTTP_READ_TIMEOUT" envDefault:"5s"`
	ReadHeaderTimeout time.Duration `env:"HTTP_READ_HEADER_TIMEOUT" envDefault:"5s"`
	WriteTimeout      time.Duration `env:"HTTP_WRITE_TIMEOUT" envDefault:"10s"`
	IdleTimeout       time.Duration `env:"HTTP_IDLE_TIMEOUT" envDefault:"60s"`
	ShutdownTimeout   time.Duration `env:"SHUTDOWN_TIMEOUT" envDefault:"10s"`

	DefaultPolicyDecision string        `env:"DEFAULT_POLICY_DECISION" envDefault:"deny"`
	TokenClockSkew        time.Duration `env:"TOKEN_CLOCK_SKEW" envDefault:"30s"`
	CacheTTL              time.Duration `env:"CONFIG_CACHE_TTL" envDefault:"30s"`
}

func Load() (Config, error) {
	var cfg Config
	if err := env.Parse(&cfg); err != nil {
		return Config{}, fmt.Errorf("parse config: %w", err)
	}

	if err := cfg.Validate(); err != nil {
		return Config{}, err
	}

	return cfg, nil
}

func (c Config) Validate() error {
	if len(c.RedisAddrs) == 0 {
		return fmt.Errorf("REDIS_ADDRS must not be empty")
	}
	if strings.TrimSpace(c.PostgresDSN) == "" {
		return fmt.Errorf("POSTGRES_DSN is required")
	}
	if strings.TrimSpace(c.JWTIssuer) == "" {
		return fmt.Errorf("JWT_ISSUER is required")
	}
	if strings.TrimSpace(c.JWTAudience) == "" {
		return fmt.Errorf("JWT_AUDIENCE is required")
	}
	if strings.TrimSpace(c.JWKSURL) == "" && strings.TrimSpace(c.LocalJWKSPath) == "" {
		return fmt.Errorf("either JWKS_URL or LOCAL_JWKS_PATH is required")
	}

	switch c.AppEnv {
	case AppEnvDevelopment, AppEnvProduction, AppEnvTest:
	default:
		return fmt.Errorf("APP_ENV must be development, production, or test")
	}

	switch c.RateLimitFailureMode {
	case RateLimitFailureOpen, RateLimitFailureClosed:
	default:
		return fmt.Errorf("RATE_LIMIT_FAILURE_MODE must be open or closed")
	}

	switch strings.ToLower(strings.TrimSpace(c.DefaultPolicyDecision)) {
	case "allow", "deny":
	default:
		return fmt.Errorf("DEFAULT_POLICY_DECISION must be allow or deny")
	}

	return nil
}
