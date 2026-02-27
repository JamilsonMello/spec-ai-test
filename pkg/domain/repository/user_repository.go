package repository

import "github.com/example/cadastro-de-usuarios/pkg/domain/entity"

// UserFilter represents filter criteria for listing users.
type UserFilter struct {
	Name  string
	Email string
}

// UserRepository provides an interface for user persistence operations.
type UserRepository interface {
	SaveUser(user *entity.User) error
	GetUserByEmail(email string) (*entity.User, error)
	GetUserByID(id string) (*entity.User, error)
	DeleteUser(id string) error
	ListUsers(filter UserFilter, page int, limit int) ([]*entity.User, int, error)
	UpdateUser(user *entity.User) error
}
