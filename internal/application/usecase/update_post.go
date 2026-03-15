package usecase

import (
	"errors"
	"time"

	"github.com/example/cadastro-de-usuarios/internal/domain"
)

var (
	ErrUnauthorizedUpdate = errors.New("Acesso negado. Você não é o autor deste post.")
)

type UpdatePostRequest struct {
	ID       string `json:"-"` 
	Content  string `json:"content"`
	AuthorID string `json:"-"` 
}

type UpdatePostResponse struct {
	ID        string `json:"id"`
	Content   string `json:"content"`
	AuthorID  string `json:"authorId"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

type UpdatePostUseCase struct {
	PostRepository domain.PostRepository
}

func NewUpdatePostUseCase(repo domain.PostRepository) *UpdatePostUseCase {
	return &UpdatePostUseCase{
		PostRepository: repo,
	}
}

func (uc *UpdatePostUseCase) Execute(req UpdatePostRequest) (*UpdatePostResponse, error) {
	post, err := uc.PostRepository.GetPostByID(req.ID)
	if err != nil {
		if errors.Is(err, domain.ErrPostNotFound) {
			return nil, domain.ErrPostNotFound
		}
		return nil, err
	}

	if post.AuthorID != req.AuthorID {
		return nil, ErrUnauthorizedUpdate
	}

	post.Content = req.Content
	if !post.IsValidContent() {
		return nil, ErrInvalidContent
	}

	post.UpdatedAt = time.Now()

	err = uc.PostRepository.UpdatePost(post)
	if err != nil {
		return nil, err
	}

	return &UpdatePostResponse{
		ID:        post.ID,
		Content:   post.Content,
		AuthorID:  post.AuthorID,
		CreatedAt: post.CreatedAt.Format(time.RFC3339),
		UpdatedAt: post.UpdatedAt.Format(time.RFC3339),
	}, nil
}
