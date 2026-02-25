package repository

import (
	"errors"
	"sync"

	"github.com/example/cadastro-de-usuarios/domain"
	"github.com/example/cadastro-de-usuarios/application/usecase"
)

var ( // Define custom errors
	ErrUserNotFound    = errors.New("user not found")
	ErrEmailAlreadyExists = errors.New("email already exists")
)

// InMemoryUserRepository implements UserRepository interface for in-memory storage.
type InMemoryUserRepository struct {
	mu    sync.RWMutex
	users map[string]*domain.User // email -> User
}

// NewInMemoryUserRepository creates a new InMemoryUserRepository.
func NewInMemoryUserRepository() *InMemoryUserRepository {
	return &InMemoryUserRepository{
		users: make(map[string]*domain.User),
	}
}

// SaveUser saves a user to the in-memory store.
func (r *InMemoryUserRepository) SaveUser(user *domain.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.users[user.Email]; exists {
		return ErrEmailAlreadyExists
	}

	r.users[user.Email] = user
	return nil
}

// GetUserByEmail retrieves a user by their email from the in-memory store.
func (r *InMemoryUserRepository) GetUserByEmail(email string) (*domain.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	user, exists := r.users[email]
	if !exists {
		return nil, ErrUserNotFound
	}
	return user, nil
}
