package platform_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/sshogun/Aegis/backend/auth"
	"github.com/sshogun/Aegis/backend/internal/health"
	"github.com/sshogun/Aegis/backend/internal/platform"
)

// mockDB implements health.Database for testing server routing.
type mockDB struct {
	err error
}

func (m *mockDB) HealthCheck(ctx context.Context) error {
	return m.err
}

func TestServerRoutes(t *testing.T) {
	cfg := &platform.Config{
		ServerAddr: ":8080",
	}

	db := &mockDB{}
	healthHandler := health.NewHandler(db)
	server := platform.NewServer(cfg, healthHandler, nil)
	handler := server.Handler()

	t.Run("Health routes are registered and liveness succeeds", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/health/live", nil)
		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("expected status %d for /health/live, got %d", http.StatusOK, rr.Code)
		}
	})

	t.Run("Readiness uses the database dependency successfully", func(t *testing.T) {
		db.err = nil // ensure success
		req := httptest.NewRequest(http.MethodGet, "/health/ready", nil)
		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("expected status %d for /health/ready on success, got %d", http.StatusOK, rr.Code)
		}
	})

	t.Run("Readiness returns error when database dependency fails", func(t *testing.T) {
		db.err = errors.New("db error")
		req := httptest.NewRequest(http.MethodGet, "/health/ready", nil)
		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)

		if rr.Code != http.StatusServiceUnavailable {
			t.Errorf("expected status %d for /health/ready on failure, got %d", http.StatusServiceUnavailable, rr.Code)
		}
	})

	t.Run("Unknown routes return 404 Not Found", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/unknown-endpoint", nil)
		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)

		if rr.Code != http.StatusNotFound {
			t.Errorf("expected status %d for unknown route, got %d", http.StatusNotFound, rr.Code)
		}
	})
}

func TestAuthMeRouteRequiresAuthentication(t *testing.T) {
	cfg := &platform.Config{
		ServerAddr: ":8080",
	}

	db := &mockDB{}
	healthHandler := health.NewHandler(db)
	authHandler := auth.NewHandler(nil, nil, nil)
	server := platform.NewServer(cfg, healthHandler, authHandler)
	handler := server.Handler()

	req := httptest.NewRequest(http.MethodGet, "/auth/me", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Errorf("expected status %d for /auth/me without session, got %d", http.StatusUnauthorized, rr.Code)
	}
}

func TestAuthLogoutRouteIsPublic(t *testing.T) {
	cfg := &platform.Config{
		ServerAddr: ":8080",
	}

	db := &mockDB{}
	healthHandler := health.NewHandler(db)
	authHandler := auth.NewHandler(nil, nil, nil)
	server := platform.NewServer(cfg, healthHandler, authHandler)
	handler := server.Handler()

	req := httptest.NewRequest(http.MethodGet, "/auth/logout", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status %d for public /auth/logout, got %d", http.StatusOK, rr.Code)
	}
}

func TestServerGracefulShutdown(t *testing.T) {
	cfg := &platform.Config{
		ServerAddr: ":0", // use random port
	}

	db := &mockDB{}
	healthHandler := health.NewHandler(db)
	server := platform.NewServer(cfg, healthHandler, nil)

	errCh := make(chan error, 1)
	go func() {
		errCh <- server.Start()
	}()

	// Wait briefly to ensure server has started
	// (a more robust test would wait for the port to open)

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // cancel immediately to test shutdown

	err := server.Shutdown(ctx)
	if err != nil && !errors.Is(err, context.Canceled) {
		t.Errorf("expected nil or context.Canceled during shutdown, got: %v", err)
	}

	// Wait for Start to return
	startErr := <-errCh
	if startErr != nil && !errors.Is(startErr, http.ErrServerClosed) {
		t.Errorf("expected http.ErrServerClosed, got: %v", startErr)
	}
}
