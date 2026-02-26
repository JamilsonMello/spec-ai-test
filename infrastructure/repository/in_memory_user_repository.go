package repository

import (
	"errors"
	"sync"
	"sort"
	"strings"

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

// FindAll retrieves all users, with optional filtering and pagination.
func (r *InMemoryUserRepository) FindAll(params usecase.ListUsersParams) ([]*domain.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var allUsers []*domain.User
	for _, user := range r.users {
		allUsers = append(allUsers, user)
	}

	// Sort by CreatedAt descending
	sort.Slice(allUsers, func(i, j int) bool {
		return allUsers[i].CreatedAt.After(allUsers[j].CreatedAt)
	})

	var filteredUsers []*domain.User
	for _, user := range allUsers {
		if (params.Name == "" || strings.Contains(strings.ToLower(user.Name), strings.ToLower(params.Name))) &&
			(params.Email == "" || strings.Contains(strings.ToLower(user.Email), strings.ToLower(params.Email))) {
			filteredUsers = append(filteredUsers, user)
		}
	}

	// Paginate
	start := (params.Page - 1) * params.Limit
	end := start + params.Limit

	if start > len(filteredUsers) {
		return []*domain.User{}, nil
	}

	if end > len(filteredUsers) {
		end = len(filteredUsers)
	}

	return filteredUsers[start:end], nil
}
