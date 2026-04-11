package usecase

import (
	"errors"
	"log"
	"time"

	"github.com/google/uuid"

	"github.com/example/cadastro-de-usuarios/internal/domain"
)

var (
	ErrInvalidContent     = errors.New("conteúdo deve ter entre 1 e 5000 caracteres")
	ErrUnauthorizedCreate = errors.New("usuário não autenticado")
	ErrInvalidRequest     = errors.New("dados de entrada inválidos")
	ErrRepositoryFailure  = errors.New("falha ao persistir dados")
)

type CreatePostRequest struct {
	Title    string `json:"title"`
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
	postRepository domain.PostRepository
}

func NewCreatePostUseCase(repo domain.PostRepository) *CreatePostUseCase {
	return &CreatePostUseCase{
		postRepository: repo,
	}
}

func (uc *CreatePostUseCase) Execute(req CreatePostRequest) (*CreatePostResponse, error) {
	if err := uc.validateRequest(req); err != nil {
		log.Printf("[WARN] falha na validação do post: %v", err)
		return nil, ErrInvalidRequest
	}

	post := uc.createPost(req)

	if err := uc.postRepository.SavePost(post); err != nil {
		log.Printf("[ERROR] falha ao salvar post: %v", err)
		return nil, ErrRepositoryFailure
	}

	return uc.buildResponse(post), nil
}

func (uc *CreatePostUseCase) validateRequest(req CreatePostRequest) error {
	if len(req.Title) < 5 || len(req.Title) > 100 {
		return errors.New("título deve ter entre 5 e 100 caracteres")
	}

	if req.Content == "" {
		return errors.New("conteúdo não pode ser vazio")
	}

	return nil
}

func (uc *CreatePostUseCase) createPost(req CreatePostRequest) *domain.Post {
	return &domain.Post{
		ID:        uuid.New().String(),
		Content:   req.Content,
		AuthorID:  req.AuthorID,
		CreatedAt: time.Now(),
	}
}

func (uc *CreatePostUseCase) buildResponse(post *domain.Post) *CreatePostResponse {
	return &CreatePostResponse{
		ID:        post.ID,
		Content:   post.Content,
		AuthorID:  post.AuthorID,
		CreatedAt: post.CreatedAt.Format(time.RFC3339),
	}
}
