package usecase

import (
	"github.com/example/cadastro-de-usuarios/domain"
)

// ValidateTokenRequest represents the request payload to validate a token
type ValidateTokenRequest struct {
	TokenString string
}

// ValidateTokenResponse represents the response payload with the token subject
type ValidateTokenResponse struct {
	Subject string
}

// TokenValidatorService defines the interface for token validation
type TokenValidatorService interface {
	Validate(token string) (*domain.TokenPayload, error)
}

// ValidateTokenUseCase validates a given token and returns the subject
type ValidateTokenUseCase struct {
	service TokenValidatorService
}

// NewValidateTokenUseCase creates a new ValidateTokenUseCase
func NewValidateTokenUseCase(s TokenValidatorService) *ValidateTokenUseCase {
	return &ValidateTokenUseCase{service: s}
}

// Execute performs the token validation logic
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