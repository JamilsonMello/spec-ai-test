package usecase

import (
	"errors"

	"github.com/example/cadastro-de-usuarios/domain"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidToken     = errors.New("token inválido ou expirado")
	ErrPasswordMismatch = errors.New("senha e confirmação não conferem")
	ErrPasswordTooShort = errors.New("senha deve ter no mínimo 8 caracteres")
)

type ResetPasswordInput struct {
	Token           string `json:"token"`
	NewPassword     string `json:"newPassword"`
	ConfirmPassword string `json:"confirmPassword"`
}

type ResetPasswordOutput struct {
	Message string `json:"message"`
}

type ResetPasswordUseCase struct {
	UserRepository             domain.UserRepository
	PasswordRecoveryRepository domain.PasswordRecoveryRepository
}

func NewResetPasswordUseCase(userRepo domain.UserRepository, recoveryRepo domain.PasswordRecoveryRepository) *ResetPasswordUseCase {
	return &ResetPasswordUseCase{
		UserRepository:             userRepo,
		PasswordRecoveryRepository: recoveryRepo,
	}
}

func (uc *ResetPasswordUseCase) validatePasswords(newPassword, confirmPassword string) error {
	if newPassword != confirmPassword {
		return ErrPasswordMismatch
	}
	if len(newPassword) < 8 {
		return ErrPasswordTooShort
	}
	return nil
}

func (uc *ResetPasswordUseCase) validateRecoveryToken(token string) (*domain.PasswordRecovery, error) {
	recovery, repositoryErr := uc.PasswordRecoveryRepository.GetPasswordRecoveryByToken(token)
	if repositoryErr != nil {
		return nil, ErrInvalidToken
	}
	if !recovery.IsValid() {
		return nil, ErrInvalidToken
	}
	return recovery, nil
}

func (uc *ResetPasswordUseCase) getUserByRecovery(recovery *domain.PasswordRecovery) (*domain.User, error) {
	user, repositoryErr := uc.UserRepository.GetUserByID(recovery.UserID)
	if repositoryErr != nil {
		return nil, ErrUserNotFound
	}
	return user, nil
}

func (uc *ResetPasswordUseCase) hashPassword(password string) (string, error) {
	hashedPassword, hashErr := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if hashErr != nil {
		return "", hashErr
	}
	return string(hashedPassword), nil
}

func (uc *ResetPasswordUseCase) updateUserPassword(user *domain.User, hashedPassword string) error {
	user.Password = hashedPassword
	updateErr := uc.UserRepository.UpdateUser(user)
	if updateErr != nil {
		return updateErr
	}
	return nil
}

func (uc *ResetPasswordUseCase) markRecoveryUsed(recovery *domain.PasswordRecovery) error {
	recovery.MarkAsUsed()
	updateErr := uc.PasswordRecoveryRepository.UpdatePasswordRecovery(recovery)
	if updateErr != nil {
		return updateErr
	}
	return nil
}

func (uc *ResetPasswordUseCase) Execute(req ResetPasswordInput) (*ResetPasswordOutput, error) {

	if err := uc.validatePasswords(req.NewPassword, req.ConfirmPassword); err != nil {
		return nil, err
	}

	recovery, err := uc.validateRecoveryToken(req.Token)
	if err != nil {
		return nil, err
	}

	user, err := uc.getUserByRecovery(recovery)
	if err != nil {
		return nil, err
	}

	hashedPassword, err := uc.hashPassword(req.NewPassword)
	if err != nil {
		return nil, err
	}

	if err := uc.updateUserPassword(user, hashedPassword); err != nil {
		return nil, err
	}

	if err := uc.markRecoveryUsed(recovery); err != nil {
		return nil, err
	}

	return &ResetPasswordOutput{
		Message: "Senha redefinida com sucesso",
	}, nil
}
