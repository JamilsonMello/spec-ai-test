package usecase

import (
	"errors"
	"time"

	"github.com/example/cadastro-de-usuarios/domain"
)

// Custom errors for post update validation
var (
	ErrUnauthorizedUpdate = errors.New("Acesso negado. Você não é o autor deste post.")
)

// UpdatePostRequest represents the input data for updating a post.
type UpdatePostRequest struct {
	ID       string `json:"-"` // From path parameter
	Content  string `json:"content"`
	AuthorID string `json:"-"` // From authentication context
}

// UpdatePostResponse represents the output data after updating a post.
type UpdatePostResponse struct {
	ID        string `json:"id"`
	Content   string `json:"content"`
	AuthorID  string `json:"authorId"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

// UpdatePostUseCase handles the business logic for post updates.
type UpdatePostUseCase struct {
	PostRepository domain.PostRepository
}

// NewUpdatePostUseCase creates a new UpdatePostUseCase.
func NewUpdatePostUseCase(repo domain.PostRepository) *UpdatePostUseCase {
	return &UpdatePostUseCase{
		PostRepository: repo,
	}
}

// Execute performs the post update process.
func (uc *UpdatePostUseCase) Execute(req UpdatePostRequest) (*UpdatePostResponse, error) {
	// 1. Fetch existing post
	post, err := uc.PostRepository.GetPostByID(req.ID)
	if err != nil {
		if errors.Is(err, domain.ErrPostNotFound) {
			return nil, domain.ErrPostNotFound
		}
		return nil, err
	}

	// 2. Verify authorization
	if post.AuthorID != req.AuthorID {
		return nil, ErrUnauthorizedUpdate
	}

	// 3. Validate new content
	post.Content = req.Content
	if !post.IsValidContent() {
		return nil, ErrInvalidContent
	}

	// 4. Update post fields
	post.UpdatedAt = time.Now()

	// 5. Persist changes
	err = uc.PostRepository.UpdatePost(post)
	if err != nil {
		return nil, err
	}

	// 6. Return response
	return &UpdatePostResponse{
		ID:        post.ID,
		Content:   post.Content,
		AuthorID:  post.AuthorID,
		CreatedAt: post.CreatedAt.Format(time.RFC3339),
		UpdatedAt: post.UpdatedAt.Format(time.RFC3339),
	}, nil
}
