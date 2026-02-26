package usecase

import (
	"errors"

	"github.com/example/cadastro-de-usuarios/domain"
	"golang.org/x/crypto/bcrypt"
)

// Custom errors for password reset
var (
	ErrInvalidToken        = errors.New("token inválido ou expirado")
	ErrPasswordMismatch    = errors.New("senha e confirmação não conferem")
	ErrPasswordTooShort    = errors.New("senha deve ter no mínimo 8 caracteres")
)

// ResetPasswordRequest is the DTO for password reset input.
type ResetPasswordRequest struct {
	Token           string `json:"token"`
	NewPassword     string `json:"newPassword"`
	ConfirmPassword string `json:"confirmPassword"`
}

// ResetPasswordResponse is the DTO for password reset output.
type ResetPasswordResponse struct {
	Message string `json:"message"`
}

// ResetPasswordUseCase handles the business logic for password reset.
type ResetPasswordUseCase struct {
	UserRepository            UserRepository
	PasswordRecoveryRepository PasswordRecoveryRepository
}

// NewResetPasswordUseCase creates a new ResetPasswordUseCase.
func NewResetPasswordUseCase(userRepo UserRepository, recoveryRepo PasswordRecoveryRepository) *ResetPasswordUseCase {
	return &ResetPasswordUseCase{
		UserRepository:            userRepo,
		PasswordRecoveryRepository: recoveryRepo,
	}
}

// Execute performs the password reset process.
func (uc *ResetPasswordUseCase) Execute(req ResetPasswordRequest) (*ResetPasswordResponse, error) {
	// 1. Validate passwords match
	if req.NewPassword != req.ConfirmPassword {
		return nil, ErrPasswordMismatch
	}

	// 2. Validate password strength
	if len(req.NewPassword) < 8 {
		return nil, ErrPasswordTooShort
	}

	// 3. Retrieve recovery token
	recovery, err := uc.PasswordRecoveryRepository.GetPasswordRecoveryByToken(req.Token)
	if err != nil {
		return nil, ErrInvalidToken
	}

	// 4. Validate token is still valid
	if !recovery.IsValid() {
		return nil, ErrInvalidToken
	}

	// 5. Retrieve user
	user, err := uc.UserRepository.GetUserByID(recovery.UserID)
	if err != nil {
		return nil, ErrUserNotFound
	}

	// 6. Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// 7. Update user password
	user.Password = string(hashedPassword)
	err = uc.UserRepository.UpdateUser(user)
	if err != nil {
		return nil, err
	}

	// 8. Mark token as used
	recovery.MarkAsUsed()
	err = uc.PasswordRecoveryRepository.UpdatePasswordRecovery(recovery)
	if err != nil {
		return nil, err
	}

	return &ResetPasswordResponse{
		Message: "Senha redefinida com sucesso",
	}, nil
}
