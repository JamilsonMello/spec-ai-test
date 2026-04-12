package usecase

import (
	"errors"
	"time"

	"github.com/google/uuid"

	"github.com/example/cadastro-de-usuarios/internal/domain"
)

var (
	ErrInvalidCommentContent = errors.New("conteúdo deve ter entre 1 e 500 caracteres")
	ErrCommentPostNotFound   = errors.New("Post não encontrado")
)

type CreateCommentRequest struct {
	Content string `json:"content"`
	PostID  string `json:"-"`
	UserID  string `json:"-"`
}

type CreateCommentResponse struct {
	ID        string `json:"id"`
	PostID    string `json:"post_id"`
	UserID    string `json:"user_id"`
	Content   string `json:"content"`
	CreatedAt string `json:"created_at"`
}

type CreateCommentUseCase struct {
	CommentRepository domain.CommentRepository
	PostRepository    domain.PostRepository
}

func NewCreateCommentUseCase(commentRepo domain.CommentRepository, postRepo domain.PostRepository) *CreateCommentUseCase {
	return &CreateCommentUseCase{
		CommentRepository: commentRepo,
		PostRepository:    postRepo,
	}
}

func (uc *CreateCommentUseCase) Execute(req CreateCommentRequest) (*CreateCommentResponse, error) {
	if err := uc.validateRequest(req); err != nil {
		return nil, err
	}

	if err := uc.validatePostExists(req.PostID); err != nil {
		return nil, err
	}

	comment := uc.createComment(req)

	if err := uc.saveComment(comment); err != nil {
		return nil, err
	}

	return uc.buildResponse(comment), nil
}

func (uc *CreateCommentUseCase) validateRequest(req CreateCommentRequest) error {
	comment := domain.NewComment(req.PostID, req.UserID, req.Content)
	if !comment.IsValidContent() {
		return ErrInvalidCommentContent
	}
	return nil
}

func (uc *CreateCommentUseCase) validatePostExists(postID string) error {
	_, err := uc.PostRepository.GetPostByID(postID)
	if err != nil {
		if err == domain.ErrPostNotFound {
			return ErrCommentPostNotFound
		}
		return err
	}
	return nil
}

func (uc *CreateCommentUseCase) createComment(req CreateCommentRequest) *domain.Comment {
	comment := domain.NewComment(req.PostID, req.UserID, req.Content)
	comment.ID = uuid.New().String()
	comment.CreatedAt = time.Now()
	return comment
}

func (uc *CreateCommentUseCase) saveComment(comment *domain.Comment) error {
	return uc.CommentRepository.SaveComment(comment)
}

func (uc *CreateCommentUseCase) buildResponse(comment *domain.Comment) *CreateCommentResponse {
	return &CreateCommentResponse{
		ID:        comment.ID,
		PostID:    comment.PostID,
		UserID:    comment.UserID,
		Content:   comment.Content,
		CreatedAt: comment.CreatedAt.Format(time.RFC3339),
	}
}
