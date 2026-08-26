package platform

import (
	"testing"
)

func TestLoadConfig(t *testing.T) {
	t.Run("Default values with valid DB URL", func(t *testing.T) {
		t.Setenv("AEGIS_SERVER_ADDR", "")
		t.Setenv("AEGIS_ENV", "")
		t.Setenv("AEGIS_DATABASE_URL", "postgres://user:pass@localhost:5432/db")

		cfg, err := LoadConfig()
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		if cfg.ServerAddr != ":8080" {
			t.Errorf("expected default ServerAddr ':8080', got %q", cfg.ServerAddr)
		}
		if cfg.Environment != EnvDevelopment {
			t.Errorf("expected default Environment 'development', got %q", cfg.Environment)
		}
	})

	t.Run("Missing database URL", func(t *testing.T) {
		t.Setenv("AEGIS_SERVER_ADDR", "")
		t.Setenv("AEGIS_ENV", "")
		t.Setenv("AEGIS_DATABASE_URL", "")

		_, err := LoadConfig()
		if err == nil {
			t.Fatal("expected error for missing database URL, got nil")
		}
		if err.Error() != "AEGIS_DATABASE_URL is required" {
			t.Errorf("unexpected error message: %v", err)
		}
	})

	t.Run("Invalid server address", func(t *testing.T) {
		t.Setenv("AEGIS_SERVER_ADDR", "invalid-address")
		t.Setenv("AEGIS_ENV", "")
		t.Setenv("AEGIS_DATABASE_URL", "postgres://user:pass@localhost:5432/db")

		_, err := LoadConfig()
		if err == nil {
			t.Fatal("expected error for invalid server address, got nil")
		}
	})

	t.Run("Invalid environment", func(t *testing.T) {
		t.Setenv("AEGIS_SERVER_ADDR", "")
		t.Setenv("AEGIS_ENV", "invalid-env")
		t.Setenv("AEGIS_DATABASE_URL", "postgres://user:pass@localhost:5432/db")

		_, err := LoadConfig()
		if err == nil {
			t.Fatal("expected error for invalid environment, got nil")
		}
	})

	t.Run("Valid configuration", func(t *testing.T) {
		t.Setenv("AEGIS_DATABASE_URL", "postgres://user:pass@localhost:5432/db")
		t.Setenv("AEGIS_SERVER_ADDR", "localhost:9090")
		t.Setenv("AEGIS_ENV", "production")

		cfg, err := LoadConfig()
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		if cfg.ServerAddr != "localhost:9090" {
			t.Errorf("expected ServerAddr 'localhost:9090', got %q", cfg.ServerAddr)
		}
		if cfg.Environment != EnvProduction {
			t.Errorf("expected Environment 'production', got %q", cfg.Environment)
		}
		if cfg.DatabaseURL != "postgres://user:pass@localhost:5432/db" {
			t.Errorf("expected DatabaseURL 'postgres://user:pass@localhost:5432/db', got %q", cfg.DatabaseURL)
		}
	})
}
