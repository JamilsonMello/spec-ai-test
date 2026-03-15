package usecase

import (
	"errors"
	"time"

	"github.com/example/cadastro-de-usuarios/internal/domain"
)

var (
	ErrRecoveryTokenNotFound = errors.New("token de recuperação não encontrado")
	ErrTokenExpired          = errors.New("token expirado")
	ErrTokenAlreadyUsed      = errors.New("token já foi utilizado")
)

type RequestPasswordRecoveryRequest struct {
	Email string `json:"email"`
}

type RequestPasswordRecoveryResponse struct {
	Token     string `json:"token"`
	Message   string `json:"message"`
	ExpiresAt string `json:"expiresAt"`
}

type RequestPasswordRecoveryUseCase struct {
	UserRepository             domain.UserRepository
	PasswordRecoveryRepository domain.PasswordRecoveryRepository
	EmailSender                domain.EmailSender
}

func NewRequestPasswordRecoveryUseCase(userRepo domain.UserRepository, recoveryRepo domain.PasswordRecoveryRepository, emailSender domain.EmailSender) *RequestPasswordRecoveryUseCase {
	return &RequestPasswordRecoveryUseCase{
		UserRepository:             userRepo,
		PasswordRecoveryRepository: recoveryRepo,
		EmailSender:                emailSender,
	}
}

func (uc *RequestPasswordRecoveryUseCase) Execute(req RequestPasswordRecoveryRequest) (*RequestPasswordRecoveryResponse, error) {
	user, err := uc.UserRepository.GetUserByEmail(req.Email)
	if err != nil {
		return nil, ErrUserNotFound
	}

	recovery, err := domain.NewPasswordRecovery(user.ID)
	if err != nil {
		return nil, err
	}

	err = uc.PasswordRecoveryRepository.SavePasswordRecovery(recovery)
	if err != nil {
		return nil, err
	}

	if uc.EmailSender != nil {
		_ = uc.EmailSender.SendPasswordRecoveryEmail(user.Email, recovery.Token)
	}

	return &RequestPasswordRecoveryResponse{
		Token:     recovery.Token,
		Message:   "Token de recuperação enviado com sucesso",
		ExpiresAt: recovery.ExpiresAt.Format(time.RFC3339),
	}, nil
}
