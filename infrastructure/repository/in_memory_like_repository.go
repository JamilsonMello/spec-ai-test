package repository

import (
	"sync"

	"github.com/example/cadastro-de-usuarios/domain"
)

// likeKey uniquely identifies a like by user and target.
type likeKey struct {
	userID   string
	targetID string
}

// InMemoryLikeRepository implements LikeRepository interface for in-memory storage.
type InMemoryLikeRepository struct {
	mu    sync.RWMutex
	likes map[likeKey]*domain.Like
}

// NewInMemoryLikeRepository creates a new InMemoryLikeRepository.
func NewInMemoryLikeRepository() *InMemoryLikeRepository {
	return &InMemoryLikeRepository{
		likes: make(map[likeKey]*domain.Like),
	}
}

// SaveLike saves a like to the in-memory store.
func (r *InMemoryLikeRepository) SaveLike(like *domain.Like) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.likes[likeKey{like.UserID, like.TargetID}] = like
	return nil
}

// DeleteLike removes a like from the in-memory store.
func (r *InMemoryLikeRepository) DeleteLike(userID, targetID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	delete(r.likes, likeKey{userID, targetID})
	return nil
}

// GetLike retrieves a like by user ID and target ID.
func (r *InMemoryLikeRepository) GetLike(userID, targetID string) (*domain.Like, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	like, ok := r.likes[likeKey{userID, targetID}]
	if !ok {
		return nil, nil
	}
	return like, nil
}

// CountLikesByTarget returns the number of likes for a given target ID.
func (r *InMemoryLikeRepository) CountLikesByTarget(targetID string) (int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	count := 0
	for key := range r.likes {
		if key.targetID == targetID {
			count++
		}
	}
	return count, nil
}
