package usecase

import "github.com/google/uuid"
import "github.com/example/cadastro-de-usuarios/domain"

// RegisterUserRequest is the DTO for user registration input.
type RegisterUserRequest struct {
	Name        string `json:"name"`
	Surname     string `json:"surname"`
	Email       string `json:"email"`
	BirthDate   string `json:"birthDate"` // YYYY-MM-DD
}

// RegisterUserResponse is the DTO for user registration output.
type RegisterUserResponse struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Surname     string    `json:"surname"`
	Email       string    `json:"email"`
	BirthDate   string    `json:"birthDate"`
}

// ListUsersRequest is the DTO for listing users input.
type ListUsersRequest struct {
	Name   string
	Email  string
	Limit  int
	Offset int
}

// ListUsersResponse is the DTO for listing users output.
type ListUsersResponse struct {
	Users      []*domain.User `json:"users"`
	TotalCount int            `json:"totalCount"`
}

// UserRepository provides an interface for user persistence operations.
type UserRepository interface {
	SaveUser(user *domain.User) error
	GetUserByEmail(email string) (*domain.User, error)
	ListUsers(request ListUsersRequest) (*ListUsersResponse, error)
}
