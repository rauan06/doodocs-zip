package api

import (
	"net/http"
	"zip/internal/handler"
)

func SetupRouter() *http.ServeMux {
	mux := http.NewServeMux()

	// Program Endpoints
	mux.HandleFunc("POST /api/archive/information", handler.InformationHandler)
	// mux.HandleFunc("POST /api/archive/fiels", )
	// mux.HandleFunc("POST /api/mail/file", )

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.NotFound(w, r) // Respond with a 404 Not Found for undefined routes
	})

	return mux
}
