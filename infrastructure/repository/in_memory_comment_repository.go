package repository

import (
	"sync"

	"github.com/example/cadastro-de-usuarios/domain"
)

// InMemoryCommentRepository implements CommentRepository interface for in-memory storage.
type InMemoryCommentRepository struct {
	mu       sync.RWMutex
	comments map[string]*domain.Comment // id -> Comment
}

// NewInMemoryCommentRepository creates a new InMemoryCommentRepository.
func NewInMemoryCommentRepository() *InMemoryCommentRepository {
	return &InMemoryCommentRepository{
		comments: make(map[string]*domain.Comment),
	}
}

// SaveComment saves a comment to the in-memory store.
func (r *InMemoryCommentRepository) SaveComment(comment *domain.Comment) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.comments[comment.ID] = comment
	return nil
}

// GetCommentByID retrieves a comment by its ID.
func (r *InMemoryCommentRepository) GetCommentByID(id string) (*domain.Comment, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	comment, ok := r.comments[id]
	if !ok {
		return nil, domain.ErrCommentNotFound
	}
	return comment, nil
}

// ListCommentsByPostID returns all comments associated with a given post ID.
func (r *InMemoryCommentRepository) ListCommentsByPostID(postID string) ([]*domain.Comment, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []*domain.Comment
	for _, c := range r.comments {
		if c.PostID == postID {
			result = append(result, c)
		}
	}
	return result, nil
}
