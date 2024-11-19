package api

import (
	"log/slog"
	"net/http"
	"zip/internal/handler"
)

func SetupRouter() *http.ServeMux {
	mux := http.NewServeMux()

	// Program Endpoints
	mux.HandleFunc("POST /api/archive/information", handler.InformationHandler)
	mux.HandleFunc("POST /api/archive/files", handler.CompressHandler)
	mux.HandleFunc("POST /api/mail/file", handler.MailHandler)

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		handler.CustomResponse(w, "404", "Not Found")
		slog.Info("404 not found", slog.String("adress", r.URL.Path))
	})

	return mux
}
