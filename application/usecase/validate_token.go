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
	if tokenValidationErr := uc.validateTokenString(req.TokenString); tokenValidationErr != nil {
		return nil, tokenValidationErr
	}

	payload, serviceValidationErr := uc.service.Validate(req.TokenString)
	if serviceValidationErr != nil {
		return nil, serviceValidationErr
	}

	if payloadValidationErr := payload.Validate(); payloadValidationErr != nil {
		return nil, payloadValidationErr
	}

	return &ValidateTokenOutput{Subject: payload.Subject}, nil
}
