package usecase

import (
	"errors"

	"github.com/example/cadastro-de-usuarios/internal/domain"
)

var (
	ErrInvalidToken     = errors.New("Token inválido, expirado ou já utilizado.")
	ErrPasswordMismatch = errors.New("senha e confirmação não conferem")
	ErrPasswordTooShort = errors.New("senha deve ter no mínimo 8 caracteres")
)

type ResetPasswordRequest struct {
	Token       string `json:"token"`
	NewPassword string `json:"new_password"`
}

type ResetPasswordResponse struct {
	Message string `json:"message"`
}

type ResetPasswordUseCase struct {
	UserRepository             domain.UserRepository
	UserRepositoryTx           domain.UserRepositoryTx
	PasswordRecoveryRepository domain.PasswordRecoveryRepository
	PasswordRecoveryRepoTx     domain.PasswordRecoveryRepositoryTx
	PasswordHasher             domain.PasswordHasher
	TransactionManager         domain.TransactionManager
}

func NewResetPasswordUseCase(
	userRepo domain.UserRepository,
	userRepoTx domain.UserRepositoryTx,
	recoveryRepo domain.PasswordRecoveryRepository,
	recoveryRepoTx domain.PasswordRecoveryRepositoryTx,
	hasher domain.PasswordHasher,
	txManager domain.TransactionManager,
) *ResetPasswordUseCase {
	return &ResetPasswordUseCase{
		UserRepository:             userRepo,
		UserRepositoryTx:           userRepoTx,
		PasswordRecoveryRepository: recoveryRepo,
		PasswordRecoveryRepoTx:     recoveryRepoTx,
		PasswordHasher:             hasher,
		TransactionManager:         txManager,
	}
}

func (uc *ResetPasswordUseCase) Execute(req ResetPasswordRequest) (*ResetPasswordResponse, error) {
	if len(req.NewPassword) < 8 {
		return nil, ErrPasswordTooShort
	}

	recovery, err := uc.PasswordRecoveryRepository.GetPasswordRecoveryByToken(req.Token)
	if err != nil {
		return nil, ErrInvalidToken
	}

	if !recovery.IsValid() {
		return nil, ErrInvalidToken
	}

	user, err := uc.UserRepository.FindUserByUuid(recovery.UserID)
	if err != nil {
		return nil, ErrInvalidToken
	}

	hashedPassword, err := uc.PasswordHasher.Hash(req.NewPassword)
	if err != nil {
		return nil, err
	}

	user.Password = hashedPassword
	recovery.MarkAsUsed()

	err = uc.TransactionManager.ExecuteInTransaction(func(tx interface{}) error {
		if err := uc.UserRepositoryTx.UpdateUserTx(tx, user); err != nil {
			return err
		}
		if err := uc.PasswordRecoveryRepoTx.UpdatePasswordRecoveryTx(tx, recovery); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return &ResetPasswordResponse{
		Message: "Senha alterada com sucesso.",
	}, nil
}
