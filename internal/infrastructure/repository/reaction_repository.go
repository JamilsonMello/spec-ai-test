package repository

import (
	"database/sql"
	"fmt"

	"github.com/example/cadastro-de-usuarios/internal/domain"
)

type ReactionRepository struct {
	db *sql.DB
}

func NewReactionRepository(db *sql.DB) *ReactionRepository {
	return &ReactionRepository{db: db}
}

func (r *ReactionRepository) GetReactionByCommentAndUser(commentID, userID, reactionType string) (*domain.Reaction, error) {
	query := `SELECT id, comment_id, user_id, reaction_type, created_at FROM comment_reactions WHERE comment_id = $1 AND user_id = $2 AND reaction_type = $3`
	reaction := &domain.Reaction{}
	err := r.db.QueryRow(query, commentID, userID, reactionType).Scan(&reaction.ID, &reaction.CommentID, &reaction.UserID, &reaction.Type, &reaction.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get reaction: %w", err)
	}
	return reaction, nil
}

func (r *ReactionRepository) SaveReaction(reaction *domain.Reaction) error {
	query := `INSERT INTO comment_reactions (id, comment_id, user_id, reaction_type, created_at) VALUES ($1, $2, $3, $4, $5)`
	_, err := r.db.Exec(query, reaction.ID, reaction.CommentID, reaction.UserID, reaction.Type, reaction.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to save reaction: %w", err)
	}
	return nil
}

func (r *ReactionRepository) DeleteReaction(id string) error {
	query := `DELETE FROM comment_reactions WHERE id = $1`
	_, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete reaction: %w", err)
	}
	return nil
}
