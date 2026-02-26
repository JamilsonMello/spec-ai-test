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

// UserRepository provides an interface for user persistence operations.
type UserRepository interface {
	SaveUser(user *domain.User) error
	GetUserByEmail(email string) (*domain.User, error)
	FindAll(params ListUsersParams) ([]*domain.User, error)
}

// ListUsersParams contains the parameters for listing users.
type ListUsersParams struct {
	Name    string
	Email   string
	Page    int
	Limit   int
}

// ListUsersResponse is the DTO for the user list output.
type ListUsersResponse struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Surname   string    `json:"surname"`
	Email     string    `json:"email"`
	BirthDate string    `json:"birthDate"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"createdAt"`
}
