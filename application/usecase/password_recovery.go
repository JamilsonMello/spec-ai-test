package usecase

import (
	"errors"
	"time"

	"github.com/example/cadastro-de-usuarios/domain"
)

// Custom errors for password recovery
var (
	ErrRecoveryTokenNotFound = errors.New("token de recuperação não encontrado")
	ErrTokenExpired          = errors.New("token expirado")
	ErrTokenAlreadyUsed      = errors.New("token já foi utilizado")
)

// RequestPasswordRecoveryRequest is the DTO for requesting password recovery.
type RequestPasswordRecoveryRequest struct {
	Email string `json:"email"`
}

// RequestPasswordRecoveryResponse is the DTO for password recovery output.
type RequestPasswordRecoveryResponse struct {
	Token     string `json:"token"` // In production, this would be sent via email
	Message   string `json:"message"`
	ExpiresAt string `json:"expiresAt"`
}

// RequestPasswordRecoveryUseCase handles the business logic for password recovery requests.
type RequestPasswordRecoveryUseCase struct {
	UserRepository            domain.UserRepository
	PasswordRecoveryRepository domain.PasswordRecoveryRepository
	EmailService              domain.EmailService
}

// NewRequestPasswordRecoveryUseCase creates a new RequestPasswordRecoveryUseCase.
func NewRequestPasswordRecoveryUseCase(userRepo domain.UserRepository, recoveryRepo domain.PasswordRecoveryRepository, emailService domain.EmailService) *RequestPasswordRecoveryUseCase {
	return &RequestPasswordRecoveryUseCase{
		UserRepository:            userRepo,
		PasswordRecoveryRepository: recoveryRepo,
		EmailService:              emailService,
	}
}

// Execute performs the password recovery request process.
func (uc *RequestPasswordRecoveryUseCase) Execute(req RequestPasswordRecoveryRequest) (*RequestPasswordRecoveryResponse, error) {
	// 1. Find user by email
	user, err := uc.UserRepository.GetUserByEmail(req.Email)
	if err != nil {
		return nil, ErrUserNotFound
	}

	// 2. Create password recovery token
	recovery, err := domain.NewPasswordRecovery(user.ID)
	if err != nil {
		return nil, err
	}

	// 3. Save recovery token
	err = uc.PasswordRecoveryRepository.SavePasswordRecovery(recovery)
	if err != nil {
		return nil, err
	}

	if uc.EmailService != nil {
		_ = uc.EmailService.SendTemplatedEmail(user.Email, "PasswordRecovery", map[string]interface{}{
			"token": recovery.Token,
			"email": user.Email,
		})
	}

	return &RequestPasswordRecoveryResponse{
		Token:     recovery.Token,
		Message:   "Token de recuperação enviado com sucesso",
		ExpiresAt: recovery.ExpiresAt.Format(time.RFC3339),
	}, nil
}
