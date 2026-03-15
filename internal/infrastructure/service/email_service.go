package service

import "log"

type EmailSender struct{}

func NewEmailSender() *EmailSender {
	return &EmailSender{}
}

func (s *EmailSender) SendPasswordRecoveryEmail(email string, token string) error {
	log.Printf("[EMAIL] Sending password recovery email to %s with token: %s", email, token)
	return nil
}
