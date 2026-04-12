package repository

import (
	"database/sql"
	"fmt"

	"github.com/example/cadastro-de-usuarios/internal/domain"
)

type CommentRepository struct {
	db *sql.DB
}

func NewCommentRepository(db *sql.DB) *CommentRepository {
	return &CommentRepository{db: db}
}

func (r *CommentRepository) SaveComment(comment *domain.Comment) error {
	query := `INSERT INTO comments (id, post_id, user_id, content, created_at) VALUES ($1, $2, $3, $4, $5)`
	_, err := r.db.Exec(query, comment.ID, comment.PostID, comment.UserID, comment.Content, comment.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to save comment: %w", err)
	}
	return nil
}

func (r *CommentRepository) GetCommentsByPostID(postID string) ([]*domain.Comment, error) {
	query := `SELECT id, post_id, user_id, content, created_at FROM comments WHERE post_id = $1 ORDER BY created_at ASC`
	rows, err := r.db.Query(query, postID)
	if err != nil {
		return nil, fmt.Errorf("failed to get comments: %w", err)
	}
	defer rows.Close()

	var comments []*domain.Comment
	for rows.Next() {
		comment := &domain.Comment{}
		err := rows.Scan(&comment.ID, &comment.PostID, &comment.UserID, &comment.Content, &comment.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan comment: %w", err)
		}
		comments = append(comments, comment)
	}
	return comments, nil
}

func (r *CommentRepository) GetCommentByID(id string) (*domain.Comment, error) {
	query := `SELECT id, post_id, user_id, content, created_at FROM comments WHERE id = $1`
	comment := &domain.Comment{}
	err := r.db.QueryRow(query, id).Scan(&comment.ID, &comment.PostID, &comment.UserID, &comment.Content, &comment.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, domain.ErrCommentNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get comment by id: %w", err)
	}
	return comment, nil
}
