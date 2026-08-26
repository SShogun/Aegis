package platform

import (
	"context"
	"net/http"
	"time"

	"github.com/sshogun/Aegis/backend/internal/health"
)

// Server holds the configured HTTP server.
type Server struct {
	httpServer *http.Server
}

// NewServer creates and configures a new HTTP server.
func NewServer(cfg *Config, healthHandler *health.Handler) *Server {
	mux := http.NewServeMux()

	// Register routes
	mux.HandleFunc("/health/live", healthHandler.Liveness())
	mux.HandleFunc("/health/ready", healthHandler.Readiness())

	httpServer := &http.Server{
		Addr:         cfg.ServerAddr,
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	return &Server{
		httpServer: httpServer,
	}
}

// Start begins listening and serving HTTP requests.
func (s *Server) Start() error {
	return s.httpServer.ListenAndServe()
}

// Shutdown gracefully stops the HTTP server.
func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}

// Handler returns the underlying http.Handler for testing purposes.
func (s *Server) Handler() http.Handler {
	return s.httpServer.Handler
}
