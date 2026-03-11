package usecase

import (
	"github.com/example/cadastro-de-usuarios/domain"
)

type ValidateTokenInput struct {
	TokenString string
}

type ValidateTokenOutput struct {
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

func (uc *ValidateTokenUseCase) validateTokenString(token string) error {
	if token == "" {
		return domain.ErrInvalidToken
	}
	return nil
}

func (uc *ValidateTokenUseCase) Execute(req ValidateTokenInput) (*ValidateTokenOutput, error) {
	if err := uc.validateTokenString(req.TokenString); err != nil {
		return nil, err
	}

	payload, validationErr := uc.service.Validate(req.TokenString)
	if validationErr != nil {
		return nil, validationErr
	}

	if err := payload.Validate(); err != nil {
		return nil, err
	}

	return &ValidateTokenOutput{Subject: payload.Subject}, nil
}
