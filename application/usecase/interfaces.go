package usecase

import "github.com/example/cadastro-de-usuarios/domain"

// PasswordRecoveryRepository provides an interface for password recovery token persistence.
type PasswordRecoveryRepository interface {
	SavePasswordRecovery(recovery *domain.PasswordRecovery) error
	GetPasswordRecoveryByToken(token string) (*domain.PasswordRecovery, error)
}

// EmailService provides an interface for sending emails.
type EmailService interface {
	SendPasswordRecoveryEmail(email string, token string) error
}
