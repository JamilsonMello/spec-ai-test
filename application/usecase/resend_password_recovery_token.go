package usecase

import (
	"time"

	"github.com/example/cadastro-de-usuarios/domain"
)

// ResendTokenRequest is the DTO for resending a password recovery token.
type ResendTokenRequest struct {
	Email string `json:"email"`
}

// ResendPasswordRecoveryTokenUseCase handles the business logic for resending recovery tokens.
type ResendPasswordRecoveryTokenUseCase struct {
	UserRepository             domain.UserRepository
	PasswordRecoveryRepository domain.PasswordRecoveryRepository
	EmailService               domain.EmailService
}

// NewResendPasswordRecoveryTokenUseCase creates a new ResendPasswordRecoveryTokenUseCase.
func NewResendPasswordRecoveryTokenUseCase(userRepo domain.UserRepository, recoveryRepo domain.PasswordRecoveryRepository, emailService domain.EmailService) *ResendPasswordRecoveryTokenUseCase {
	return &ResendPasswordRecoveryTokenUseCase{
		UserRepository:             userRepo,
		PasswordRecoveryRepository: recoveryRepo,
		EmailService:               emailService,
	}
}

// Execute performs the token resend process.
func (uc *ResendPasswordRecoveryTokenUseCase) Execute(req ResendTokenRequest) error {
	user, err := uc.UserRepository.GetUserByEmail(req.Email)
	if err != nil {
		return nil
	}

	latest, err := uc.PasswordRecoveryRepository.GetLatestPasswordRecoveryByUserID(user.ID)
	if err == nil && time.Since(latest.CreatedAt) < 60*time.Second {
		return nil
	}

	if err := uc.PasswordRecoveryRepository.InvalidateAllUserTokens(user.ID); err != nil {
		return err
	}

	recovery, err := domain.NewPasswordRecovery(user.ID)
	if err != nil {
		return err
	}

	if err := uc.PasswordRecoveryRepository.SavePasswordRecovery(recovery); err != nil {
		return err
	}

	if uc.EmailService != nil {
		_ = uc.EmailService.SendPasswordRecoveryEmail(user.Email, recovery.Token)
	}

	return nil
}
