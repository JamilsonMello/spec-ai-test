package usecase

import (
	"errors"
	"time"

	"github.com/example/cadastro-de-usuarios/domain"
)

var (
	ErrRecoveryTokenNotFound = errors.New("token de recuperação não encontrado")
	ErrTokenExpired          = errors.New("token expirado")
	ErrTokenAlreadyUsed      = errors.New("token já foi utilizado")
)

type RequestPasswordRecoveryInput struct {
	Email string `json:"email"`
}

type RequestPasswordRecoveryOutput struct {
	Token     string `json:"token"`
	Message   string `json:"message"`
	ExpiresAt string `json:"expiresAt"`
}

type RequestPasswordRecoveryUseCase struct {
	UserRepository             domain.UserRepository
	PasswordRecoveryRepository domain.PasswordRecoveryRepository
	EmailService               domain.EmailService
}

func NewRequestPasswordRecoveryUseCase(userRepo domain.UserRepository, recoveryRepo domain.PasswordRecoveryRepository, emailService domain.EmailService) *RequestPasswordRecoveryUseCase {
	return &RequestPasswordRecoveryUseCase{
		UserRepository:             userRepo,
		PasswordRecoveryRepository: recoveryRepo,
		EmailService:               emailService,
	}
}

func (uc *RequestPasswordRecoveryUseCase) getUserByEmail(email string) (*domain.User, error) {
	user, repositoryErr := uc.UserRepository.GetUserByEmail(email)
	if repositoryErr != nil {
		return nil, ErrUserNotFound
	}
	return user, nil
}

func (uc *RequestPasswordRecoveryUseCase) createPasswordRecovery(userID string) (*domain.PasswordRecovery, error) {
	recovery, creationErr := domain.NewPasswordRecovery(userID)
	if creationErr != nil {
		return nil, creationErr
	}
	return recovery, nil
}

func (uc *RequestPasswordRecoveryUseCase) saveRecovery(recovery *domain.PasswordRecovery) error {
	repositoryErr := uc.PasswordRecoveryRepository.SavePasswordRecovery(recovery)
	if repositoryErr != nil {
		return repositoryErr
	}
	return nil
}

func (uc *RequestPasswordRecoveryUseCase) sendRecoveryEmail(email string, token string) {
	if uc.EmailService != nil {
		_ = uc.EmailService.SendPasswordRecoveryEmail(email, token)
	}
}

func (uc *RequestPasswordRecoveryUseCase) Execute(req RequestPasswordRecoveryInput) (*RequestPasswordRecoveryOutput, error) {

	user, userErr := uc.getUserByEmail(req.Email)
	if userErr != nil {
		return nil, userErr
	}

	recovery, creationErr := uc.createPasswordRecovery(user.ID)
	if creationErr != nil {
		return nil, creationErr
	}

	if err := uc.saveRecovery(recovery); err != nil {
		return nil, err
	}

	uc.sendRecoveryEmail(user.Email, recovery.Token)

	return &RequestPasswordRecoveryOutput{
		Token:     recovery.Token,
		Message:   "Token de recuperação enviado com sucesso",
		ExpiresAt: recovery.ExpiresAt.Format(time.RFC3339),
	}, nil
}
