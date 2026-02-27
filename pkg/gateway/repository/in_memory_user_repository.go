package repository

import (
	"errors"
	"sort"
	"strings"
	"sync"

	"github.com/example/cadastro-de-usuarios/pkg/domain/entity"
	"github.com/example/cadastro-de-usuarios/pkg/domain/repository"
)

var ( // Define custom errors
	ErrUserNotFound       = errors.New("user not found")
	ErrEmailAlreadyExists = errors.New("email already exists")
)

// InMemoryUserRepository implements UserRepository interface for in-memory storage.
type InMemoryUserRepository struct {
	mu    sync.RWMutex
	users map[string]*entity.User // email -> User
	ids   map[string]string      // id -> email
}

// NewInMemoryUserRepository creates a new InMemoryUserRepository.
func NewInMemoryUserRepository() *InMemoryUserRepository {
	return &InMemoryUserRepository{
		users: make(map[string]*entity.User),
		ids:   make(map[string]string),
	}
}

// SaveUser saves a user to the in-memory store.
func (r *InMemoryUserRepository) SaveUser(user *entity.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.users[user.Email]; exists {
		return ErrEmailAlreadyExists
	}
	r.users[user.Email] = user
	r.ids[user.ID] = user.Email
	return nil
}

// GetUserByEmail retrieves a user by their email from the in-memory store.
func (r *InMemoryUserRepository) GetUserByEmail(email string) (*entity.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	user, exists := r.users[email]
	if !exists {
		return nil, ErrUserNotFound
	}
	return user, nil
}

// GetUserByID retrieves a user by their ID from the in-memory store.
func (r *InMemoryUserRepository) GetUserByID(id string) (*entity.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	email, exists := r.ids[id]
	if !exists {
		return nil, ErrUserNotFound
	}
	user, exists := r.users[email]
	if !exists {
		return nil, ErrUserNotFound
	}
	return user, nil
}

// DeleteUser removes a user from the in-memory store by ID.
func (r *InMemoryUserRepository) DeleteUser(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	email, ok := r.ids[id]
	if !ok {
		return ErrUserNotFound
	}

	delete(r.users, email)
	delete(r.ids, id)
	return nil
}


// ListUsers retrieves a paginated list of users with optional filters.
func (r *InMemoryUserRepository) ListUsers(filter repository.UserFilter, page int, limit int) ([]*entity.User, int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	// Collect all users into a slice
	var filteredUsers []*entity.User
	for _, user := range r.users {
		// Apply name filter (case-insensitive partial match)
		if filter.Name != "" && !strings.Contains(strings.ToLower(user.Name), strings.ToLower(filter.Name)) {
			continue
		}
		// Apply email filter (case-insensitive partial match)
		if filter.Email != "" && !strings.Contains(strings.ToLower(user.Email), strings.ToLower(filter.Email)) {
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
		return []*entity.User{}, totalCount, nil
	}

	end := offset + limit
	if end > totalCount {
		end = totalCount
	}

	return filteredUsers[offset:end], totalCount, nil
}

// UpdateUser updates an existing user in the in-memory store.
func (r *InMemoryUserRepository) UpdateUser(user *entity.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.ids[user.ID]; !exists {
		return ErrUserNotFound
	}

	r.users[user.Email] = user
	return nil
}
