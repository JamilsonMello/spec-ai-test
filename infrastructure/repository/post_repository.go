package repository

import (
	"database/sql"
	"fmt"

	"github.com/example/cadastro-de-usuarios/domain"
)

type PostRepository struct {
	db *sql.DB
}

func NewPostRepository(db *sql.DB) *PostRepository {
	return &PostRepository{db: db}
}

func (postRepository *PostRepository) SavePost(post *domain.Post) error {
	query := `INSERT INTO posts (id, content, author_id, created_at) VALUES ($1, $2, $3, $4)`
	_, err := postRepository.db.Exec(query, post.ID, post.Content, post.AuthorID, post.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to save post: %w", err)
	}
	return nil
}
