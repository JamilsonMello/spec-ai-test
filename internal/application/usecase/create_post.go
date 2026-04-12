package usecase

import (
	"errors"
	"time"

	"github.com/google/uuid"

	"github.com/example/cadastro-de-usuarios/internal/domain"
)

var ErrInvalidUUIDFormat = errors.New("invalid_uuid_format")

var (
	ErrInvalidContent     = errors.New("conteúdo deve ter entre 1 e 5000 caracteres")
	ErrUnauthorizedCreate = errors.New("usuário não autenticado")
)

type CreatePostRequest struct {
	Content     string  `json:"content"`
	AuthorID    string  `json:"-"`
	CommunityID *string `json:"community_id,omitempty"`
}

type CreatePostResponse struct {
	ID          string `json:"id"`
	Content     string `json:"content"`
	AuthorID    string `json:"authorId"`
	CommunityID string `json:"community_id,omitempty"`
	CreatedAt   string `json:"createdAt"`
}

type CreatePostUseCase struct {
	PostRepository      domain.PostRepository
	CommunityRepository domain.CommunityRepository
}

func NewCreatePostUseCase(repo domain.PostRepository, communityRepo domain.CommunityRepository) *CreatePostUseCase {
	return &CreatePostUseCase{
		PostRepository:      repo,
		CommunityRepository: communityRepo,
	}
}

func (uc *CreatePostUseCase) Execute(req CreatePostRequest) (*CreatePostResponse, error) {
	if err := uc.validateRequest(req); err != nil {
		return nil, err
	}

	if req.CommunityID != nil && *req.CommunityID != "" {
		if _, err := uuid.Parse(*req.CommunityID); err != nil {
			return nil, ErrInvalidUUIDFormat
		}

		exists, err := uc.CommunityRepository.ExistsByID(*req.CommunityID)
		if err != nil {
			return nil, err
		}
		if !exists {
			return nil, domain.ErrCommunityNotFound
		}
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

	if req.CommunityID != nil && *req.CommunityID != "" {
		post.CommunityID = *req.CommunityID
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
		ID:          post.ID,
		Content:     post.Content,
		AuthorID:    post.AuthorID,
		CommunityID: post.CommunityID,
		CreatedAt:   post.CreatedAt.Format(time.RFC3339),
	}
}
