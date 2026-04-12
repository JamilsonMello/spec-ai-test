package usecase

import (
	"errors"
	"time"

	"github.com/google/uuid"

	"github.com/example/cadastro-de-usuarios/internal/domain"
)

var (
	ErrInvalidCommunityName        = errors.New("nome da comunidade deve ter entre 1 e 100 caracteres")
	ErrInvalidCommunityDescription = errors.New("descrição da comunidade deve ter entre 1 e 500 caracteres")
	ErrCommunityAlreadyExists      = errors.New("community already exists")
	ErrUnauthorizedCommunity       = errors.New("usuário não autenticado")
)

type CreateCommunityRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	OwnerID     string `json:"-"`
}

type CreateCommunityResponse struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	OwnerID string `json:"owner_id"`
}

type CreateCommunityUseCase struct {
	CommunityRepository domain.CommunityRepository
}

func NewCreateCommunityUseCase(repo domain.CommunityRepository) *CreateCommunityUseCase {
	return &CreateCommunityUseCase{
		CommunityRepository: repo,
	}
}

func (uc *CreateCommunityUseCase) Execute(req CreateCommunityRequest) (*CreateCommunityResponse, error) {
	if err := uc.validateRequest(req); err != nil {
		return nil, err
	}

	community, err := uc.createCommunity(req)
	if err != nil {
		return nil, err
	}

	exists, err := uc.CommunityRepository.ExistsByName(community.Name)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrCommunityAlreadyExists
	}

	if err := uc.CommunityRepository.SaveCommunity(community); err != nil {
		return nil, err
	}

	return uc.buildResponse(community), nil
}

func (uc *CreateCommunityUseCase) validateRequest(req CreateCommunityRequest) error {
	if req.OwnerID == "" {
		return ErrUnauthorizedCommunity
	}
	return nil
}

func (uc *CreateCommunityUseCase) createCommunity(req CreateCommunityRequest) (*domain.Community, error) {
	now := time.Now()
	community := &domain.Community{
		ID:          uuid.New().String(),
		Name:        req.Name,
		Description: req.Description,
		OwnerID:     req.OwnerID,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if !community.IsValidName() {
		return nil, ErrInvalidCommunityName
	}

	if !community.IsValidDescription() {
		return nil, ErrInvalidCommunityDescription
	}

	return community, nil
}

func (uc *CreateCommunityUseCase) buildResponse(community *domain.Community) *CreateCommunityResponse {
	return &CreateCommunityResponse{
		ID:      community.ID,
		Name:    community.Name,
		OwnerID: community.OwnerID,
	}
}
