package usecase

import (
	"errors"

	"github.com/google/uuid"

	"github.com/example/cadastro-de-usuarios/internal/domain"
)

var (
	ErrInvalidUserID    = errors.New("invalid user ID format")
	ErrUserNotFound     = errors.New("user not found")
	ErrUnauthorizedRole = errors.New("unauthorized: admin role required")
)

type DeleteUserUseCase struct {
	UserRepository domain.UserRepository
}

func NewDeleteUserUseCase(repo domain.UserRepository) *DeleteUserUseCase {
	return &DeleteUserUseCase{
		UserRepository: repo,
	}
}

func (uc *DeleteUserUseCase) Execute(userID string, userRole string) error {
	if userRole != "admin" {
		return ErrUnauthorizedRole
	}

	if _, err := uuid.Parse(userID); err != nil {
		return ErrInvalidUserID
	}

	user, err := uc.UserRepository.FindUserByUuid(userID)
	if err != nil {
		return ErrUserNotFound
	}
	if user == nil {
		return ErrUserNotFound
	}

	if err := uc.UserRepository.DeleteUser(userID); err != nil {
		return err
	}

	return nil
}
