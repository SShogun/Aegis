package health

import (
	"context"
	"encoding/json"
	"net/http"
)

// Database defines the required operations for health checks.
// Using an interface allows for easy mocking in tests.
type Database interface {
	HealthCheck(ctx context.Context) error
}

// Handler manages health check HTTP endpoints.
type Handler struct {
	db Database
}

// NewHandler creates a new Handler.
func NewHandler(db Database) *Handler {
	return &Handler{
		db: db,
	}
}

// Liveness confirms that the Go process is running.
// It does not check dependencies like PostgreSQL.
func (h *Handler) Liveness() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]string{"status": "UP"})
	}
}

// Readiness confirms that the application is ready to accept traffic.
// It verifies that PostgreSQL is reachable.
func (h *Handler) Readiness() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := h.db.HealthCheck(r.Context()); err != nil {
			// In a real application, the error should be logged internally here.
			// We return a generic 503 to avoid exposing database credentials or internal states.
			writeJSON(w, http.StatusServiceUnavailable, map[string]string{"status": "DOWN"})
			return
		}
		writeJSON(w, http.StatusOK, map[string]string{"status": "UP"})
	}
}

// writeJSON is a small helper for standardizing JSON responses.
func writeJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(data)
}
