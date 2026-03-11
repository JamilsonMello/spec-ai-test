package service

import "log"

type EmailService struct{}

func NewEmailService() *EmailService {
	return &EmailService{}
}

func (emailService *EmailService) SendPasswordRecoveryEmail(email string, token string) error {

	log.Printf("[EMAIL] Sending password recovery email to %s with token: %s", email, token)
	return nil
}
