package usecase

import (
	"errors"
	"time"

	"github.com/google/uuid"

	"github.com/example/cadastro-de-usuarios/internal/domain"
)

var (
	ErrInvalidCommentContent = errors.New("Invalid content")
	ErrUnauthorizedComment   = errors.New("Unauthorized")
)

type CreateCommentRequest struct {
	Content string `json:"content"`
	PostID  string `json:"-"`
	UserID  string `json:"-"`
}

type CreateCommentResponse struct {
	ID      string `json:"id"`
	Content string `json:"content"`
}

type CreateCommentUseCase struct {
	CommentRepository domain.CommentRepository
}

func NewCreateCommentUseCase(repo domain.CommentRepository) *CreateCommentUseCase {
	return &CreateCommentUseCase{
		CommentRepository: repo,
	}
}

func (uc *CreateCommentUseCase) Execute(req CreateCommentRequest) (*CreateCommentResponse, error) {
	if req.UserID == "" {
		return nil, ErrUnauthorizedComment
	}

	comment := &domain.Comment{
		ID:        uuid.New().String(),
		PostID:    req.PostID,
		UserID:    req.UserID,
		Content:   req.Content,
		CreatedAt: time.Now(),
	}

	if !comment.IsValidContent() {
		return nil, ErrInvalidCommentContent
	}

	if err := uc.CommentRepository.SaveComment(comment); err != nil {
		return nil, err
	}

	return &CreateCommentResponse{
		ID:      comment.ID,
		Content: comment.Content,
	}, nil
}
