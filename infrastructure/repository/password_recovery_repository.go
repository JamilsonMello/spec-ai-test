package repository

import (
	"errors"
	"sync"

	"github.com/example/cadastro-de-usuarios/domain"
)

var ErrRecoveryTokenNotFound = errors.New("recovery token not found")

// InMemoryPasswordRecoveryRepository implements PasswordRecoveryRepository interface for in-memory storage.
type InMemoryPasswordRecoveryRepository struct {
	mu       sync.RWMutex
	recovery map[string]*domain.PasswordRecovery // token -> PasswordRecovery
}

// NewInMemoryPasswordRecoveryRepository creates a new InMemoryPasswordRecoveryRepository.
func NewInMemoryPasswordRecoveryRepository() *InMemoryPasswordRecoveryRepository {
	return &InMemoryPasswordRecoveryRepository{
		recovery: make(map[string]*domain.PasswordRecovery),
	}
}

// SavePasswordRecovery saves a password recovery token to the in-memory store.
func (r *InMemoryPasswordRecoveryRepository) SavePasswordRecovery(recovery *domain.PasswordRecovery) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.recovery[recovery.Token] = recovery
	return nil
}

// GetPasswordRecoveryByToken retrieves a password recovery token by its token string.
func (r *InMemoryPasswordRecoveryRepository) GetPasswordRecoveryByToken(token string) (*domain.PasswordRecovery, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	recovery, exists := r.recovery[token]
	if !exists {
		return nil, ErrRecoveryTokenNotFound
	}
	return recovery, nil
}

// UpdatePasswordRecovery updates an existing password recovery token in the in-memory store.
func (r *InMemoryPasswordRecoveryRepository) UpdatePasswordRecovery(recovery *domain.PasswordRecovery) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.recovery[recovery.Token] = recovery
	return nil
}
