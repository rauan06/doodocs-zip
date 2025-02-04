package main

import (
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"zip/internal/api"
	"zip/internal/config"
)

func main() {
	conf := config.SetupConfig()
	slog.SetDefault(conf.Logger)

	addr := fmt.Sprintf("0.0.0.0:%d", conf.Port)

	srv := &http.Server{
		Addr:    addr,
		Handler: api.SetupRouter(),
	}

	slog.Info("starting http server", slog.String("Env", conf.Env), slog.String("addr", addr))

	// Start server and handle potential errors
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
