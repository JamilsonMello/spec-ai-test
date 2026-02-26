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

// DeleteUserUseCase handles the business logic for deleting a user by ID.
type DeleteUserUseCase struct {
	UserRepository domain.UserRepository
}

// NewDeleteUserUseCase creates a new DeleteUserUseCase.
func NewDeleteUserUseCase(repo domain.UserRepository) *DeleteUserUseCase {
	return &DeleteUserUseCase{
		UserRepository: repo,
	}
}

// Execute performs the user deletion process.
func (uc *DeleteUserUseCase) Execute(userID string, userRole string) error {
	// Validate admin role
	if userRole != "admin" {
		return ErrUnauthorizedRole
	}

	// Validate UUID format
	if _, err := uuid.Parse(userID); err != nil {
		return ErrInvalidUserID
	}

	// Check if user exists
	user, err := uc.UserRepository.GetUserByID(userID)
	if err != nil {
		return ErrUserNotFound
	}
	if user == nil {
		return ErrUserNotFound
	}

	// Delete user from repository
	if err := uc.UserRepository.DeleteUser(userID); err != nil {
		return err
	}

	return nil
}
