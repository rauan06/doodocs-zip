package handler

import (
	"io"
	"log/slog"
	"net/http"
	"strings"
	"zip/internal/service"
)

func MailHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		CustomResponse(w, "error", "unable to parse form: "+err.Error())
		return
	}

	// Retrieve the file (PDF)
	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		CustomResponse(w, "error", "could not retrieve file from form-data: "+err.Error())
		return
	}
	defer file.Close()

	// Log file info
	logFileInfo(fileHeader)

	// Validate content type
	contentType := fileHeader.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "application/pdf") {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		CustomResponse(w, "error", "invalid content type: "+contentType)
		return
	}

	// Retrieve emails
	emailData, _, err := r.FormFile("emails")
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		CustomResponse(w, "error", "could not retrieve emails from form-data: "+err.Error())
		return
	}

	emails, err := io.ReadAll(emailData)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		CustomResponse(w, "error", "failed while reading emails: "+err.Error())
		return
	}

	// Split and clean up email addresses
	emailList := strings.FieldsFunc(string(emails), func(r rune) bool {
		return r == ',' || r == '\n' || r == '\r'
	})

	// Send email with attachment
	if err := service.SendMail(emailList, []byte("Here is your document."), file, fileHeader.Filename); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		CustomResponse(w, "error", err.Error())
		return
	}

	// Log success
	slog.Info("mails were sent successfully")
	CustomResponse(w, "success", "mails were sent successfully")
}
