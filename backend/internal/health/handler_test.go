package health_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/sshogun/Aegis/backend/internal/health"
)

// mockDB implements the health.Database interface for testing.
type mockDB struct {
	err error
}

func (m *mockDB) HealthCheck(ctx context.Context) error {
	return m.err
}

func TestLiveness(t *testing.T) {
	handler := health.NewHandler(&mockDB{})
	req := httptest.NewRequest(http.MethodGet, "/health/live", nil)
	rr := httptest.NewRecorder()

	handler.Liveness().ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status %v, got %v", http.StatusOK, rr.Code)
	}

	expected := `{"status":"UP"}` + "\n"
	if rr.Body.String() != expected {
		t.Errorf("expected body %q, got %q", expected, rr.Body.String())
	}
}

func TestReadiness_Success(t *testing.T) {
	handler := health.NewHandler(&mockDB{err: nil})
	req := httptest.NewRequest(http.MethodGet, "/health/ready", nil)
	rr := httptest.NewRecorder()

	handler.Readiness().ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status %v, got %v", http.StatusOK, rr.Code)
	}

	expected := `{"status":"UP"}` + "\n"
	if rr.Body.String() != expected {
		t.Errorf("expected body %q, got %q", expected, rr.Body.String())
	}
}

func TestReadiness_Failure(t *testing.T) {
	handler := health.NewHandler(&mockDB{err: errors.New("database connection refused")})
	req := httptest.NewRequest(http.MethodGet, "/health/ready", nil)
	rr := httptest.NewRecorder()

	handler.Readiness().ServeHTTP(rr, req)

	if rr.Code != http.StatusServiceUnavailable {
		t.Errorf("expected status %v, got %v", http.StatusServiceUnavailable, rr.Code)
	}

	expected := `{"status":"DOWN"}` + "\n"
	if rr.Body.String() != expected {
		t.Errorf("expected body %q, got %q", expected, rr.Body.String())
	}
}
