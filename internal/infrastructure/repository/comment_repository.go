package repository

import (
	"database/sql"
	"fmt"

	"github.com/example/cadastro-de-usuarios/internal/domain"
	"github.com/example/cadastro-de-usuarios/pkg/db/queries"
)

type CommentRepository struct {
	db *sql.DB
}

func NewCommentRepository(db *sql.DB) *CommentRepository {
	return &CommentRepository{db: db}
}

func (r *CommentRepository) SaveComment(comment *domain.Comment) error {
	_, err := r.db.Exec(queries.InsertCommentSQL, comment.ID, comment.PostID, comment.UserID, comment.Content, comment.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to save comment: %w", err)
	}
	return nil
}
