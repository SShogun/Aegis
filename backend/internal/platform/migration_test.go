package platform_test

import (
	"os"
	"testing"
	"path/filepath"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func TestMigrationExecution_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	dbURL := os.Getenv("AEGIS_DATABASE_URL")
	if dbURL == "" {
		t.Skip("skipping integration test: AEGIS_DATABASE_URL is not set")
	}

	// Assuming migrations are in the backend/migrations directory.
	// Since tests are run from internal/platform, the path is ../../migrations
	migrationsDir, err := filepath.Abs("../../migrations")
	if err != nil {
		t.Fatalf("failed to resolve migrations directory: %v", err)
	}
	
	sourceURL := "file://" + migrationsDir

	m, err := migrate.New(sourceURL, dbURL)
	if err != nil {
		t.Fatalf("failed to create migrate instance: %v", err)
	}
	defer m.Close()

	// Apply migrations (Up)
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		t.Fatalf("failed to apply migrations (Up): %v", err)
	}

	// Rollback migrations (Down)
	if err := m.Down(); err != nil && err != migrate.ErrNoChange {
		t.Fatalf("failed to rollback migrations (Down): %v", err)
	}
}
