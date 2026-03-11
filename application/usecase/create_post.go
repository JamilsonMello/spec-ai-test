package usecase

import (
	"errors"
	"time"

	"github.com/google/uuid"

	"github.com/example/cadastro-de-usuarios/domain"
)

var (
	ErrInvalidContent     = errors.New("conteúdo deve ter entre 1 e 600 caracteres")
	ErrUnauthorizedCreate = errors.New("usuário não autenticado")
)

type CreatePostInput struct {
	Content  string `json:"content"`
	AuthorID string `json:"-"`
}

type CreatePostOutput struct {
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

func (uc *CreatePostUseCase) validateAuthorID(authorID string) error {
	if authorID == "" {
		return ErrUnauthorizedCreate
	}
	return nil
}

func (uc *CreatePostUseCase) createPost(content, authorID string) *domain.Post {
	return &domain.Post{
		Content:  content,
		AuthorID: authorID,
	}
}

func (uc *CreatePostUseCase) validatePostContent(post *domain.Post) error {
	if !post.IsValidContent() {
		return ErrInvalidContent
	}
	return nil
}

func (uc *CreatePostUseCase) savePost(post *domain.Post) error {
	repositoryErr := uc.PostRepository.SavePost(post)
	if repositoryErr != nil {
		return repositoryErr
	}
	return nil
}

func (uc *CreatePostUseCase) Execute(req CreatePostInput) (*CreatePostOutput, error) {

	if authorValidationErr := uc.validateAuthorID(req.AuthorID); authorValidationErr != nil {
		return nil, authorValidationErr
	}

	post := uc.createPost(req.Content, req.AuthorID)

	if contentValidationErr := uc.validatePostContent(post); contentValidationErr != nil {
		return nil, contentValidationErr
	}

	post.ID = uuid.New().String()
	post.CreatedAt = time.Now()

	if saveErr := uc.savePost(post); saveErr != nil {
		return nil, saveErr
	}

	return &CreatePostOutput{
		ID:        post.ID,
		Content:   post.Content,
		AuthorID:  post.AuthorID,
		CreatedAt: post.CreatedAt.Format(time.RFC3339),
	}, nil
}
