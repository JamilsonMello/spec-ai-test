package repository

import (
	"errors"
	"sort"
	"strings"
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

// ListUsers retrieves a paginated list of users with optional filters.
func (r *InMemoryUserRepository) ListUsers(filter usecase.UserFilter, page int, limit int) ([]*domain.User, int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	// Collect all users into a slice
	var filteredUsers []*domain.User
	for _, user := range r.users {
		// Apply name filter (case-insensitive partial match)
		if filter.Name != "" && !strings.Contains(strings.ToLower(user.Name), filter.Name) {
			continue
		}
		// Apply email filter (case-insensitive partial match)
		if filter.Email != "" && !strings.Contains(strings.ToLower(user.Email), filter.Email) {
			continue
		}
		filteredUsers = append(filteredUsers, user)
	}

	// Sort by CreatedAt descending (most recent first)
	sort.Slice(filteredUsers, func(i, j int) bool {
		return filteredUsers[i].CreatedAt.After(filteredUsers[j].CreatedAt)
	})

	// Calculate total count before pagination
	totalCount := len(filteredUsers)

	// Apply pagination
	offset := (page - 1) * limit
	if offset >= totalCount {
		return []*domain.User{}, totalCount, nil
	}

	end := offset + limit
	if end > totalCount {
		end = totalCount
	}

	return filteredUsers[offset:end], totalCount, nil
}
