package config

import (
	"log/slog"
	"os"
	"zip/models"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

// SetLogger updates the Logger value in the configuration based on the specified Env.
func SetLogger(conf *models.Config) {
	switch conf.Env {
	case envLocal:
		conf.Logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envProd:
		conf.Logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	case envDev:
		conf.Logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}))
	}
}
