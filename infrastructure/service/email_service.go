package service

import "log"

// EmailService implements email sending functionality.
type EmailService struct{}

// NewEmailService creates a new EmailService.
func NewEmailService() *EmailService {
	return &EmailService{}
}

// SendPasswordRecoveryEmail sends a password recovery email to the user.
// In production, this would integrate with a real email service (SendGrid, AWS SES, etc.)
func (s *EmailService) SendPasswordRecoveryEmail(email string, token string) error {
	// For demo purposes, just log the email
	log.Printf("[EMAIL] Sending password recovery email to %s with token: %s", email, token)
	return nil
}
