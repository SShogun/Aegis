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

	"github.com/sshogun/Aegis/backend/auth"
	"github.com/sshogun/Aegis/backend/internal/health"
	"github.com/sshogun/Aegis/backend/internal/platform"
	"github.com/sshogun/Aegis/backend/users"
)

func main() {
	if err := run(); err != nil {
		log.Fatalf("application startup failed: %v", err)
	}
}

func run() error {
	log.Println("Starting Aegis backend...")

	// 1. Load configuration.
	cfg, err := platform.LoadConfig()
	if err != nil {
		return err
	}

	log.Printf("Configuration loaded: %s", cfg.String())

	ctx := context.Background()

	// 2. Create the PostgreSQL connection pool.
	db, err := platform.NewDatabase(ctx, cfg)
	if err != nil {
		return err
	}

	// 3. Ensure database resources are closed during shutdown.
	defer func() {
		log.Println("Closing database resources...")
		db.Close()
	}()

	// 4. Create health handlers.
	healthHandler := health.NewHandler(db)

	// 5. Create the OIDC client.
	oidcClient, err := auth.NewOIDC(ctx, auth.Config{
		ProviderName: os.Getenv("AEGIS_OIDC_PROVIDER"),
		IssuerURL:    os.Getenv("AEGIS_OIDC_ISSUER_URL"),
		ClientID:     os.Getenv("AEGIS_OIDC_CLIENT_ID"),
		ClientSecret: os.Getenv("AEGIS_OIDC_CLIENT_SECRET"),
		RedirectURL:  os.Getenv("AEGIS_OIDC_REDIRECT_URL"),
	})
	if err != nil {
		return err
	}

	userRepository := users.NewRepository(db.Pool)
	userService := users.NewService(userRepository)

	sessionRepository := users.NewSessionRepository(db.Pool)
	sessionService := users.NewSessionService(sessionRepository)

	authHandler := auth.NewHandler(
		oidcClient,
		userService,
		sessionService,
	)

	// 9. Create the HTTP server.
	srv := platform.NewServer(cfg, healthHandler, authHandler)

	serverErrors := make(chan error, 1)

	// 10. Start the server.
	go func() {
		log.Printf("HTTP server listening on %s", cfg.ServerAddr)
		serverErrors <- srv.Start()
	}()

	// 11. Handle operating-system shutdown signals.
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-serverErrors:
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			return err
		}

	case sig := <-shutdown:
		log.Printf("Received signal %v, starting graceful shutdown", sig)

		shutdownCtx, cancel := context.WithTimeout(
			context.Background(),
			10*time.Second,
		)
		defer cancel()

		if err := srv.Shutdown(shutdownCtx); err != nil {
			log.Printf("Graceful shutdown failed: %v", err)
			return err
		}
	}

	log.Println("Aegis backend stopped gracefully.")

	return nil
}
