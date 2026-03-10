package repository

import (
	"github.com/example/cadastro-de-usuarios/domain"
)

type MockUserRepository struct {
	users               map[string]*domain.User
	ListUsersError      error
	LastListUsersFilter domain.UserFilter
	LastListUsersPage   int
	LastListUsersLimit  int
	ListUsersCallCount  int
}

func NewMockUserRepository() *MockUserRepository {
	return &MockUserRepository{
		users: make(map[string]*domain.User),
	}
}

func (r *MockUserRepository) SaveUser(user *domain.User) error {
	if _, exists := r.users[user.Email]; exists {
		return domain.ErrEmailAlreadyExists
	}
	r.users[user.Email] = user
	return nil
}

func (r *MockUserRepository) GetUserByEmail(email string) (*domain.User, error) {
	user, exists := r.users[email]
	if !exists {
		return nil, domain.ErrUserNotFound
	}
	return user, nil
}

func (r *MockUserRepository) GetUserByID(id string) (*domain.User, error) {
	for _, user := range r.users {
		if user.ID == id {
			return user, nil
		}
	}
	return nil, domain.ErrUserNotFound
}

func (r *MockUserRepository) UpdateUser(user *domain.User) error {
	if _, exists := r.users[user.Email]; !exists {
		return domain.ErrUserNotFound
	}
	r.users[user.Email] = user
	return nil
}

func (r *MockUserRepository) DeleteUser(id string) error {
	for email, user := range r.users {
		if user.ID == id {
			delete(r.users, email)
			return nil
		}
	}
	return domain.ErrUserNotFound
}

func (r *MockUserRepository) ListUsers(filter domain.UserFilter, page int, limit int) ([]*domain.User, int, error) {
	r.LastListUsersFilter = filter
	r.LastListUsersPage = page
	r.LastListUsersLimit = limit
	r.ListUsersCallCount++
	if r.ListUsersError != nil {
		return nil, 0, r.ListUsersError
	}
	// Simulate empty result for simplicity
	return []*domain.User{}, 0, nil
}
