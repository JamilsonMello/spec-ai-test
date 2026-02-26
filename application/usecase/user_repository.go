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

// UserFilter represents filter criteria for listing users.
type UserFilter struct {
	Name  string
	Email string
}

// UserRepository provides an interface for user persistence operations.
type UserRepository interface {
	SaveUser(user *domain.User) error
	GetUserByEmail(email string) (*domain.User, error)
	GetUserByID(id string) (*domain.User, error)
	UpdateUser(user *domain.User) error
	ListUsers(filter UserFilter, page int, limit int) ([]*domain.User, int, error)
}
