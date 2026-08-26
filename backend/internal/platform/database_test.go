package platform_test

import (
	"context"
	"os"
	"testing"

	"github.com/sshogun/Aegis/backend/internal/platform"
)

func TestNewDatabase_InvalidURL(t *testing.T) {
	cfg := &platform.Config{
		DatabaseURL: "postgres://%gh&%ij", // invalid URL to trigger parse error
	}

	_, err := platform.NewDatabase(context.Background(), cfg)
	if err == nil {
		t.Fatal("expected error parsing invalid DB URL, got nil")
	}
}

// TestDatabase_Integration requires a real PostgreSQL database to run.
// Set AEGIS_DATABASE_URL to a valid connection string to run this test.
func TestDatabase_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	dbURL := os.Getenv("AEGIS_DATABASE_URL")
	if dbURL == "" {
		t.Skip("skipping integration test: AEGIS_DATABASE_URL is not set")
	}

	cfg := &platform.Config{
		DatabaseURL: dbURL,
		Environment: platform.EnvTest,
	}

	ctx := context.Background()
	db, err := platform.NewDatabase(ctx, cfg)
	if err != nil {
		t.Fatalf("failed to connect to test database: %v", err)
	}
	defer db.Close()

	if err := db.HealthCheck(ctx); err != nil {
		t.Fatalf("expected health check to pass, got: %v", err)
	}
}
