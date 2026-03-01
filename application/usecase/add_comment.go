package usecase

import (
	"errors"
	"time"

	"github.com/google/uuid"

	"github.com/example/cadastro-de-usuarios/domain"
)

var (
	ErrInvalidCommentContent = errors.New("o conteúdo do comentário deve ter entre 1 e 400 caracteres")
	ErrCommentUnauthorized   = errors.New("usuário não autenticado")
)

// AddCommentRequest represents the input data for adding a comment.
type AddCommentRequest struct {
	PostID   string `json:"postId"`
	AuthorID string `json:"-"`
	Content  string `json:"content"`
}

// AddCommentResponse represents the output data after adding a comment.
type AddCommentResponse struct {
	ID        string `json:"id"`
	PostID    string `json:"postId"`
	AuthorID  string `json:"authorId"`
	Content   string `json:"content"`
	CreatedAt string `json:"createdAt"`
}

// AddCommentUseCase handles the business logic for adding a comment to a post.
type AddCommentUseCase struct {
	PostRepository    domain.PostRepository
	CommentRepository domain.CommentRepository
}

// NewAddCommentUseCase creates a new AddCommentUseCase.
func NewAddCommentUseCase(postRepo domain.PostRepository, commentRepo domain.CommentRepository) *AddCommentUseCase {
	return &AddCommentUseCase{
		PostRepository:    postRepo,
		CommentRepository: commentRepo,
	}
}

// Execute performs the comment creation process.
func (uc *AddCommentUseCase) Execute(req AddCommentRequest) (*AddCommentResponse, error) {
	if req.AuthorID == "" {
		return nil, ErrCommentUnauthorized
	}

	_, err := uc.PostRepository.GetPostByID(req.PostID)
	if err != nil {
		return nil, domain.ErrPostNotFound
	}

	comment := &domain.Comment{
		PostID:   req.PostID,
		AuthorID: req.AuthorID,
		Content:  req.Content,
	}

	if !comment.IsValidContent() {
		return nil, ErrInvalidCommentContent
	}

	comment.ID = uuid.New().String()
	comment.CreatedAt = time.Now()

	if err := uc.CommentRepository.SaveComment(comment); err != nil {
		return nil, err
	}

	return &AddCommentResponse{
		ID:        comment.ID,
		PostID:    comment.PostID,
		AuthorID:  comment.AuthorID,
		Content:   comment.Content,
		CreatedAt: comment.CreatedAt.Format(time.RFC3339),
	}, nil
}
