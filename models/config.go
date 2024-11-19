package models

import "log/slog"

type Config struct {
	Help     bool
	Port     int    // Default if empty: 8000
	Env      string // Default if empty: "local"
	Logger   *slog.Logger
	Email    string
	Password string
}
