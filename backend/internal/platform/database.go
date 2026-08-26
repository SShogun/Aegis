package platform

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Database manages the PostgreSQL connection pool.
type Database struct {
	Pool *pgxpool.Pool
}

// NewDatabase creates a new PostgreSQL connection pool using the provided configuration.
func NewDatabase(ctx context.Context, cfg *Config) (*Database, error) {
	poolConfig, err := pgxpool.ParseConfig(cfg.DatabaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse database configuration: %w", err)
	}

	// Configure reasonable pool limits for local development
	poolConfig.MaxConns = 10
	poolConfig.MinConns = 2
	poolConfig.MaxConnLifetime = time.Hour
	poolConfig.MaxConnIdleTime = 30 * time.Minute

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		// Do not log the raw connection error as it might leak the full connection string
		return nil, fmt.Errorf("failed to create connection pool")
	}

	return &Database{
		Pool: pool,
	}, nil
}

// HealthCheck verifies PostgreSQL connectivity.
func (db *Database) HealthCheck(ctx context.Context) error {
	if err := db.Pool.Ping(ctx); err != nil {
		return fmt.Errorf("database health check failed: %w", err)
	}
	return nil
}

// Close gracefully shuts down the connection pool.
func (db *Database) Close() {
	if db.Pool != nil {
		db.Pool.Close()
	}
}
