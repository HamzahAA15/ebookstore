package notification

import (
	"ebookstore/utils/config"

	"github.com/gofiber/fiber/v2/log"
	"gopkg.in/gomail.v2"
)

type EmailPayload struct {
	To      string
	Subject string
	Body    string
}

type GmailNotification struct {
	gomailDialer *gomail.Dialer
}

func NewGmailNotification(gomailDialer *gomail.Dialer) INotificationService {
	return &GmailNotification{
		gomailDialer: gomailDialer,
	}
}

func (g *GmailNotification) SendNotification(payload EmailPayload) error {
	m := gomail.NewMessage()
	m.SetHeader("From", config.CONFIG_AUTH_EMAIL)
	m.SetHeader("To", payload.To)
	m.SetHeader("Subject", payload.Subject)
	m.SetBody("text/html", payload.Body)

	// Send the email to Bob, Cora and Dan.
	if err := g.gomailDialer.DialAndSend(m); err != nil {
		log.Info("Error sending email: ", err)
	}

	return nil
}
