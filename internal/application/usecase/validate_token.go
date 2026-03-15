package usecase

import (
	"github.com/example/cadastro-de-usuarios/internal/domain"
)

type ValidateTokenRequest struct {
	TokenString string
}

type ValidateTokenResponse struct {
	Subject string
}

type TokenValidatorService interface {
	Validate(token string) (*domain.TokenPayload, error)
}

type ValidateTokenUseCase struct {
	service TokenValidatorService
}

func NewValidateTokenUseCase(s TokenValidatorService) *ValidateTokenUseCase {
	return &ValidateTokenUseCase{service: s}
}

func (uc *ValidateTokenUseCase) Execute(req ValidateTokenRequest) (*ValidateTokenResponse, error) {
	if req.TokenString == "" {
		return nil, domain.ErrInvalidToken
	}

	payload, err := uc.service.Validate(req.TokenString)
	if err != nil {
		return nil, err
	}

	if err := payload.Validate(); err != nil {
		return nil, err
	}

	return &ValidateTokenResponse{Subject: payload.Subject}, nil
}
