package repository

import (
	"sync"

	"github.com/example/cadastro-de-usuarios/domain"
)

// InMemoryPostRepository implements PostRepository interface for in-memory storage.
type InMemoryPostRepository struct {
	mu    sync.RWMutex
	posts map[string]*domain.Post // id -> Post
}

// NewInMemoryPostRepository creates a new InMemoryPostRepository.
func NewInMemoryPostRepository() *InMemoryPostRepository {
	return &InMemoryPostRepository{
		posts: make(map[string]*domain.Post),
	}
}

// SavePost saves a post to the in-memory store.
func (r *InMemoryPostRepository) SavePost(post *domain.Post) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.posts[post.ID] = post
	return nil
}
