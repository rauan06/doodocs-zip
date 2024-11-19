package handler

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"mime/multipart"
	"net/http"
	"zip/internal/service"
)

func InformationHandler(w http.ResponseWriter, r *http.Request) {
	// Parse the form with a 10 MB limit
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		CustomResponse(w, "error", "unable to parse form")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Retrieve the file
	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		CustomResponse(w, "error", "could not retrieve file from form-data")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer file.Close()

	logFileInfo(fileHeader)

	// Unzip the file
	zipParser := &service.Zip{}
	if err := zipParser.Unzip(file); err != nil {
		CustomResponse(w, "error", fmt.Sprintf("failed to unzip file: %v", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Create and send the response
	if err := sendZipResponse(w, zipParser, fileHeader.Filename); err != nil {
		CustomResponse(w, "error", fmt.Sprintf("failed to create response: %v", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func logFileInfo(fileHeader *multipart.FileHeader) {
	slog.Info(
		"opened file",
		slog.String("file_name", fileHeader.Filename),
		slog.Int("file_size", int(fileHeader.Size)),
	)
}

func sendZipResponse(w http.ResponseWriter, zipParser service.ZipParser, filename string) error {
	response, err := zipParser.CreateResponse()
	if err != nil {
		return err
	}

	response.Filename = filename

	// Send JSON response
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(response)
}

func CustomResponse(w http.ResponseWriter, key string, value interface{}) {
	response := map[string]interface{}{key: value}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
