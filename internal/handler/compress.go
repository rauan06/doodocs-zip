package handler

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"zip/internal/service"
	"zip/models"
)

func CompressHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		handleError(w, "unable to parse form", http.StatusBadRequest)
		return
	}

	files, err := processFiles(r)
	if err != nil {
		handleError(w, err.Error(), http.StatusBadRequest)
		return
	}

	zipService := &service.Zip{}
	zipData, err := zipService.CreateZip(files)
	if err != nil {
		handleError(w, fmt.Sprintf("failed while creating zip file: %v", err), http.StatusInternalServerError)
		return
	}

	// Set headers and return the zip file
	w.Header().Set("Content-Type", "application/zip")
	w.Header().Set("Content-Disposition", "attachment; filename=files.zip")
	w.Write(zipData.Bytes())

	slog.Info("zip file created successfully")
}

func processFiles(r *http.Request) ([]models.File, error) {
	var files []models.File

	for _, fileHeader := range r.MultipartForm.File["files[]"] {
		file, err := fileHeader.Open()
		if err != nil {
			return nil, fmt.Errorf("failed while opening file %s", fileHeader.Filename)
		}
		defer file.Close()

		content, err := io.ReadAll(file)
		if err != nil {
			return nil, fmt.Errorf("failed while reading file %s", fileHeader.Filename)
		}

		files = append(files, models.File{
			FilePath: fileHeader.Filename,
			Content:  content,
			Size:     float64(len(content)),
			MimeType: http.DetectContentType(content),
		})
	}

	return files, nil
}

func handleError(w http.ResponseWriter, message string, statusCode int) {
	CustomResponse(w, "error", message)
	w.WriteHeader(statusCode)
}
