package repository

import (
	"database/sql"
	"fmt"

	"github.com/example/cadastro-de-usuarios/domain"
)

type PostgreSQLPostRepository struct {
	db *sql.DB
}

func NewPostgreSQLPostRepository(db *sql.DB) *PostgreSQLPostRepository {
	return &PostgreSQLPostRepository{db: db}
}

func (r *PostgreSQLPostRepository) SavePost(post *domain.Post) error {
	query := `INSERT INTO posts (id, content, author_id, created_at) VALUES ($1, $2, $3, $4)`
	_, err := r.db.Exec(query, post.ID, post.Content, post.AuthorID, post.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to save post: %w", err)
	}
	return nil
}
