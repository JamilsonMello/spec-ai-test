package usecase

import (
	"github.com/example/cadastro-de-usuarios/internal/domain"
)

type RequestPasswordRecoveryRequest struct {
	Email string `json:"email"`
}

type RequestPasswordRecoveryResponse struct {
	Message string `json:"message"`
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
	genericResponse := &RequestPasswordRecoveryResponse{
		Message: "Se um usuário com esse e-mail existir, você receberá instruções em breve.",
	}

	user, err := uc.UserRepository.GetUserByEmail(req.Email)
	if err != nil || user == nil {
		return genericResponse, nil
	}

	recovery, err := domain.NewPasswordRecovery(user.ID)
	if err != nil {
		return genericResponse, nil
	}

	err = uc.PasswordRecoveryRepository.SavePasswordRecovery(recovery)
	if err != nil {
		return genericResponse, nil
	}

	if uc.EmailSender != nil {
		_ = uc.EmailSender.SendPasswordRecoveryEmail(user.Email, recovery.Token)
	}

	return genericResponse, nil
}
