package main

import (
	"context"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/SShogun/Aegis/internal/app"
	"github.com/SShogun/Aegis/internal/config"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	logger := newLogger(string(cfg.AppEnv))
	slog.SetDefault(logger)

	logger.Info("starting Aegis",
		slog.String("env", string(cfg.AppEnv)),
		slog.String("service", cfg.ServiceName),
	)

	if err := app.Run(ctx, cfg, logger); err != nil {
		logger.Error("Aegis stopped with error", slog.String("error", err.Error()))
		os.Exit(1)
	}
	logger.Info("Aegis stopped cleanly")
}

func newLogger(appEnv string) *slog.Logger {
	opts := &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}

	if appEnv == "production" {
		return slog.New(slog.NewJSONHandler(os.Stdout, opts))
	}

	return slog.New(slog.NewTextHandler(os.Stdout, opts))
}
