package usecase

import "github.com/google/uuid"
import "github.com/example/cadastro-de-usuarios/domain"

// RegisterUserRequest is the DTO for user registration input.
type RegisterUserRequest struct {
	Nome           string `json:"nome"`
	Sobrenome      string `json:"sobrenome"`
	Email          string `json:"email"`
	DataNascimento string `json:"dataNascimento"` // YYYY-MM-DD
	Senha          string `json:"senha"`
}

// RegisterUserResponse is the DTO for user registration output.
type RegisterUserResponse struct {
	ID             uuid.UUID `json:"id"`
	Nome           string    `json:"nome"`
	Sobrenome      string    `json:"sobrenome"`
	Email          string    `json:"email"`
	DataNascimento string    `json:"dataNascimento"`
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
