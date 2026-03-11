package usecase

import (
	"errors"

	"github.com/google/uuid"

	"github.com/example/cadastro-de-usuarios/domain"
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

	if _, uuidParseErr := uuid.Parse(userID); uuidParseErr != nil {
		return ErrInvalidUserID
	}

	user, repositoryErr := uc.UserRepository.GetUserByID(userID)
	if repositoryErr != nil {
		return ErrUserNotFound
	}
	if user == nil {
		return ErrUserNotFound
	}

	if deleteErr := uc.UserRepository.DeleteUser(userID); deleteErr != nil {
		return deleteErr
	}

	return nil
}
