package service

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"log/slog"
	"net/smtp"
	"zip/internal/config"
)

const (
	smtpHost    = "smtp.gmail.com"
	smtpPort    = "587"
	boundary    = "boundary123"
	contentType = "application/vnd.openxmlformats-officedocument.wordprocessingml.document"
)

func SendMail(emails []string, message []byte, file io.Reader, fileName string) error {
	// Load email credentials from config
	senderEmail := config.Config.Email
	senderPassword := config.Config.Password

	// Set up authentication
	auth := smtp.PlainAuth("", senderEmail, senderPassword, smtpHost)

	// Read the file content
	fileContent, err := readFile(file)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	// Build the email headers and body
	emailHeaders := buildHeaders(senderEmail, "Auto message by gopher.")
	emailBody, err := buildBody("Hello World!", fileContent, fileName)
	if err != nil {
		return fmt.Errorf("failed to build email body: %w", err)
	}

	// Combine headers and body
	var msg bytes.Buffer
	for key, value := range emailHeaders {
		msg.WriteString(fmt.Sprintf("%s: %s\r\n", key, value))
	}
	msg.WriteString("\r\n")
	msg.Write(emailBody.Bytes())

	// Send the email
	if err := smtp.SendMail(fmt.Sprintf("%s:%s", smtpHost, smtpPort), auth, senderEmail, emails, msg.Bytes()); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	slog.Info("email sent successfully")
	return nil
}

func readFile(file io.Reader) ([]byte, error) {
	var buf bytes.Buffer
	_, err := buf.ReadFrom(file)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func buildHeaders(from, subject string) map[string]string {
	return map[string]string{
		"From":         from,
		"Subject":      subject,
		"MIME-Version": "1.0",
		"Content-Type": fmt.Sprintf("multipart/mixed; boundary=\"%s\"", boundary),
	}
}

func buildBody(bodyText string, fileContent []byte, fileName string) (*bytes.Buffer, error) {
	var body bytes.Buffer

	// Text part
	body.WriteString(fmt.Sprintf("--%s\r\n", boundary))
	body.WriteString("Content-Type: text/plain; charset=UTF-8\r\n\r\n")
	body.WriteString(bodyText + "\r\n")

	// File attachment part
	body.WriteString(fmt.Sprintf("--%s\r\n", boundary))
	body.WriteString(fmt.Sprintf("Content-Type: %s; name=\"%s\"\r\n", contentType, fileName))
	body.WriteString("Content-Transfer-Encoding: base64\r\n")
	body.WriteString(fmt.Sprintf("Content-Disposition: attachment; filename=\"%s\"\r\n\r\n", fileName))

	// Encode file content in Base64
	encodedContent := base64.StdEncoding.EncodeToString(fileContent)
	body.WriteString(encodedContent + "\r\n")

	// End boundary
	body.WriteString(fmt.Sprintf("--%s--", boundary))

	return &body, nil
}
