package usecase

import (
	"errors"
	"time"

	"github.com/google/uuid"

	"github.com/example/cadastro-de-usuarios/domain"
)

// Custom errors for post creation validation
var (
	ErrInvalidContent     = errors.New("conteúdo deve ter entre 1 e 600 caracteres")
	ErrUnauthorizedCreate = errors.New("usuário não autenticado")
)

// CreatePostRequest represents the input data for creating a post.
type CreatePostRequest struct {
	Content  string `json:"content"`
	AuthorID string `json:"-"` // Set from authentication context
}

// CreatePostResponse represents the output data after creating a post.
type CreatePostResponse struct {
	ID        string `json:"id"`
	Content   string `json:"content"`
	AuthorID  string `json:"authorId"`
	CreatedAt string `json:"createdAt"`
}

// CreatePostUseCase handles the business logic for post creation.
type CreatePostUseCase struct {
	PostRepository PostRepository
}

// NewCreatePostUseCase creates a new CreatePostUseCase.
func NewCreatePostUseCase(repo PostRepository) *CreatePostUseCase {
	return &CreatePostUseCase{
		PostRepository: repo,
	}
}

// Execute performs the post creation process.
func (uc *CreatePostUseCase) Execute(req CreatePostRequest) (*CreatePostResponse, error) {
	// 1. Check authentication
	if req.AuthorID == "" {
		return nil, ErrUnauthorizedCreate
	}

	// 2. Create post entity
	post := &domain.Post{
		Content:  req.Content,
		AuthorID: req.AuthorID,
	}

	// 3. Validate content
	if !post.IsValidContent() {
		return nil, ErrInvalidContent
	}

	// 4. Generate ID and timestamp
	post.ID = uuid.New().String()
	post.CreatedAt = time.Now()

	// 5. Save post
	err := uc.PostRepository.SavePost(post)
	if err != nil {
		return nil, err
	}

	// 6. Return response
	return &CreatePostResponse{
		ID:        post.ID,
		Content:   post.Content,
		AuthorID:  post.AuthorID,
		CreatedAt: post.CreatedAt.Format(time.RFC3339),
	}, nil
}
