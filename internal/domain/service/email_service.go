package service

type EmailService interface {
	SendPasswordRecoveryEmail(email string, token string) error
}
