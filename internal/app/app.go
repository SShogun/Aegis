package app

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/SShogun/Aegis/internal/config"
)

type App struct {
	cfg    config.Config
	logger *slog.Logger
}

func New(cfg config.Config, logger *slog.Logger) *App {
	return &App{cfg: cfg, logger: logger}
}

func Run(ctx context.Context, cfg config.Config, logger *slog.Logger) error {
	a := New(cfg, logger)

	mux := http.NewServeMux()
	RegisterHealthHandlers(mux, a)

	srv := &http.Server{
		Addr:         cfg.HTTPAddr,
		Handler:      mux,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		IdleTimeout:  cfg.IdleTimeout,
	}

	serverErr := make(chan error, 1)
	go func() {
		a.logger.Info("http server listening", slog.String("addr", cfg.HTTPAddr))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			serverErr <- err
		} else {
			serverErr <- nil
		}
	}()

	select {
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
		defer cancel()
		a.logger.Info("shutting down http server")
		if err := srv.Shutdown(shutdownCtx); err != nil {
			a.logger.Error("server shutdown error", slog.String("err", err.Error()))
			return err
		}
		return ctx.Err()
	case err := <-serverErr:
		// server failed or crashed
		return err
	}
}

func (a *App) ReadinessCheck(ctx context.Context) error {
	_ = ctx
	return nil
}
