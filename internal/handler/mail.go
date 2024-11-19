package handler

import (
	"io"
	"net/http"
	"strings"
	"zip/internal/service"
)

func MailHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		CustomResponse(w, "error", "unable to parse form")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		CustomResponse(w, "error", "could not retrieve file from form-data")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer file.Close()

	logFileInfo(fileHeader)

	if r.Header.Get("Content-Type") != "application/vnd.openxmlformats-officedocument.wordprocessingml.document" && r.Header.Get("Content-Type") != "application/pdf" {
		CustomResponse(w, "error", "invalid content type")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	emailData, _, err := r.FormFile("emails")
	if err != nil {
		CustomResponse(w, "error", "could not retrieve emails from form-data")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	emails, err := io.ReadAll(emailData)
	if err != nil {
		CustomResponse(w, "error", "failed while reading emails")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	emailList := strings.SplitAfter(string(emails), ",")

	data, err := io.ReadAll(emailData)
	if err != nil {
		CustomResponse(w, "error", "failed while reading file")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	service.SendMail(emailList, data)
}
