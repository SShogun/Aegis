package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sshogun/Aegis/backend/internal/health"
	"github.com/sshogun/Aegis/backend/internal/platform"
)

func main() {
	if err := run(); err != nil {
		log.Fatalf("application startup failed: %v", err)
	}
}

func run() error {
	log.Println("Starting Aegis backend...")

	// 1. Load configuration.
	// 2. Fail clearly if configuration is invalid.
	cfg, err := platform.LoadConfig()
	if err != nil {
		return err
	}
	log.Printf("Configuration loaded: %s", cfg.String())

	// Context for database initialization
	ctx := context.Background()

	// 3. Create the PostgreSQL connection pool.
	db, err := platform.NewDatabase(ctx, cfg)
	if err != nil {
		return err
	}
	
	// 4. Ensure database resources are closed during shutdown.
	defer func() {
		log.Println("Closing database resources...")
		db.Close()
	}()

	// 5. Create health handlers.
	healthHandler := health.NewHandler(db)

	// 6. Create the HTTP server.
	srv := platform.NewServer(cfg, healthHandler)

	// Channel to listen for errors from the HTTP server.
	serverErrors := make(chan error, 1)

	// 7. Start the server.
	go func() {
		log.Printf("HTTP server listening on %s", cfg.ServerAddr)
		serverErrors <- srv.Start()
	}()

	// 8. Handle operating-system shutdown signals.
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-serverErrors:
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			return err
		}
	case sig := <-shutdown:
		log.Printf("Received signal %v, starting graceful shutdown", sig)

		// 9. Shut down the HTTP server gracefully with a bounded timeout.
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := srv.Shutdown(shutdownCtx); err != nil {
			log.Printf("Graceful shutdown failed: %v", err)
			return err
		}
	}
	
	// 10. Database resources will be cleanly closed by the defer block above.

	log.Println("Aegis backend stopped gracefully.")
	return nil
}
