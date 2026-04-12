package usecase

import (
	"errors"
	"time"

	"github.com/google/uuid"

	"github.com/example/cadastro-de-usuarios/internal/domain"
)

var (
	ErrInvalidReactionType   = errors.New("Invalid reaction type")
	ErrUnauthorizedReaction  = errors.New("Unauthorized")
	ErrCommentNotFoundToggle = errors.New("Comment not found")
)

type ToggleReactionRequest struct {
	Type      string `json:"type"`
	CommentID string `json:"-"`
	UserID    string `json:"-"`
}

type ToggleReactionResponse struct {
	ID   string `json:"id"`
	Type string `json:"type"`
}

type ToggleCommentReactionUseCase struct {
	ReactionRepository domain.ReactionRepository
	CommentRepository  domain.CommentRepository
}

func NewToggleCommentReactionUseCase(reactionRepo domain.ReactionRepository, commentRepo domain.CommentRepository) *ToggleCommentReactionUseCase {
	return &ToggleCommentReactionUseCase{
		ReactionRepository: reactionRepo,
		CommentRepository:  commentRepo,
	}
}

func (uc *ToggleCommentReactionUseCase) Execute(req ToggleReactionRequest) (*ToggleReactionResponse, error) {
	if req.UserID == "" {
		return nil, ErrUnauthorizedReaction
	}

	if req.Type != "LIKE" && req.Type != "DISLIKE" {
		return nil, ErrInvalidReactionType
	}

	_, err := uc.CommentRepository.GetCommentByID(req.CommentID)
	if err != nil {
		if err == domain.ErrCommentNotFound {
			return nil, ErrCommentNotFoundToggle
		}
		return nil, err
	}

	existing, err := uc.ReactionRepository.GetReactionByCommentAndUser(req.CommentID, req.UserID, req.Type)
	if err != nil {
		return nil, err
	}

	if existing != nil {
		if err := uc.ReactionRepository.DeleteReaction(existing.ID); err != nil {
			return nil, err
		}
		return nil, nil
	}

	reaction := &domain.Reaction{
		ID:        uuid.New().String(),
		CommentID: req.CommentID,
		UserID:    req.UserID,
		Type:      req.Type,
		CreatedAt: time.Now(),
	}

	if err := uc.ReactionRepository.SaveReaction(reaction); err != nil {
		return nil, err
	}

	return &ToggleReactionResponse{
		ID:   reaction.ID,
		Type: reaction.Type,
	}, nil
}
