package usecase

import (
	"errors"

	"github.com/example/cadastro-de-usuarios/internal/domain"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidToken     = errors.New("token inválido ou expirado")
	ErrPasswordMismatch = errors.New("senha e confirmação não conferem")
	ErrPasswordTooShort = errors.New("senha deve ter no mínimo 8 caracteres")
)

type ResetPasswordRequest struct {
	Token           string `json:"token"`
	NewPassword     string `json:"newPassword"`
	ConfirmPassword string `json:"confirmPassword"`
}

type ResetPasswordResponse struct {
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

func (uc *ResetPasswordUseCase) Execute(req ResetPasswordRequest) (*ResetPasswordResponse, error) {
	if err := uc.validatePassword(req); err != nil {
		return nil, err
	}

	recovery, err := uc.getValidRecovery(req.Token)
	if err != nil {
		return nil, err
	}

	user, err := uc.getUser(recovery.UserID)
	if err != nil {
		return nil, err
	}

	if err := uc.updateUserPassword(user, req.NewPassword); err != nil {
		return nil, err
	}

	if err := uc.markRecoveryAsUsed(recovery); err != nil {
		return nil, err
	}

	return uc.buildSuccessResponse(), nil
}

func (uc *ResetPasswordUseCase) validatePassword(req ResetPasswordRequest) error {
	if req.NewPassword != req.ConfirmPassword {
		return ErrPasswordMismatch
	}

	if len(req.NewPassword) < 8 {
		return ErrPasswordTooShort
	}

	return nil
}

func (uc *ResetPasswordUseCase) getValidRecovery(token string) (*domain.PasswordRecovery, error) {
	recovery, err := uc.PasswordRecoveryRepository.GetPasswordRecoveryByToken(token)
	if err != nil {
		return nil, ErrInvalidToken
	}

	if !recovery.IsValid() {
		return nil, ErrInvalidToken
	}

	return recovery, nil
}

func (uc *ResetPasswordUseCase) getUser(userID string) (*domain.User, error) {
	user, err := uc.UserRepository.FindUserByUuid(userID)
	if err != nil {
		return nil, ErrUserNotFound
	}
	return user, nil
}

func (uc *ResetPasswordUseCase) updateUserPassword(user *domain.User, newPassword string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.Password = string(hashedPassword)
	err = uc.UserRepository.UpdateUser(user)
	if err != nil {
		return err
	}

	return nil
}

func (uc *ResetPasswordUseCase) markRecoveryAsUsed(recovery *domain.PasswordRecovery) error {
	recovery.MarkAsUsed()
	err := uc.PasswordRecoveryRepository.UpdatePasswordRecovery(recovery)
	if err != nil {
		return err
	}
	return nil
}

func (uc *ResetPasswordUseCase) buildSuccessResponse() *ResetPasswordResponse {
	return &ResetPasswordResponse{
		Message: "Senha redefinida com sucesso",
	}
}
