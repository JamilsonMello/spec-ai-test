package usecase

import (
	"errors"
	"time"

	"github.com/google/uuid"

	"github.com/example/cadastro-de-usuarios/internal/domain"
)

var (
	ErrInvalidContent     = errors.New("conteúdo deve ter entre 1 e 5000 caracteres")
	ErrUnauthorizedCreate = errors.New("usuário não autenticado")
)

type CreatePostRequest struct {
	Content  string `json:"content"`
	AuthorID string `json:"-"`
}

type CreatePostResponse struct {
	ID        string `json:"id"`
	Content   string `json:"content"`
	AuthorID  string `json:"authorId"`
	CreatedAt string `json:"createdAt"`
}

type CreatePostUseCase struct {
	PostRepository domain.PostRepository
}

func NewCreatePostUseCase(repo domain.PostRepository) *CreatePostUseCase {
	return &CreatePostUseCase{
		PostRepository: repo,
	}
}

func (uc *CreatePostUseCase) Execute(req CreatePostRequest) (*CreatePostResponse, error) {
	if err := uc.validateRequest(req); err != nil {
		return nil, err
	}

	post, err := uc.createPost(req)
	if err != nil {
		return nil, err
	}

	if err := uc.savePost(post); err != nil {
		return nil, err
	}

	return uc.buildResponse(post), nil
}

func (uc *CreatePostUseCase) validateRequest(req CreatePostRequest) error {
	if req.AuthorID == "" {
		return ErrUnauthorizedCreate
	}
	return nil
}

func (uc *CreatePostUseCase) createPost(req CreatePostRequest) (*domain.Post, error) {
	post := &domain.Post{
		Content:   req.Content,
		AuthorID:  req.AuthorID,
		ID:        uuid.New().String(),
		CreatedAt: time.Now(),
	}

	if !post.IsValidContent() {
		return nil, ErrInvalidContent
	}

	return post, nil
}

func (uc *CreatePostUseCase) savePost(post *domain.Post) error {
	err := uc.PostRepository.SavePost(post)
	if err != nil {
		return err
	}
	return nil
}

func (uc *CreatePostUseCase) buildResponse(post *domain.Post) *CreatePostResponse {
	return &CreatePostResponse{
		ID:        post.ID,
		Content:   post.Content,
		AuthorID:  post.AuthorID,
		CreatedAt: post.CreatedAt.Format(time.RFC3339),
	}
}
