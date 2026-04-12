package mocks

import (
	"time"

	"github.com/example/cadastro-de-usuarios/internal/domain"
)

type JWTValidatorMock struct {
	Subject string
}

func NewJWTValidatorMock(subject string) *JWTValidatorMock {
	return &JWTValidatorMock{Subject: subject}
}

func (m *JWTValidatorMock) Validate(token string) (*domain.TokenPayload, error) {
	if token == "" || token == "invalid" {
		return nil, domain.ErrInvalidToken
	}

	return &domain.TokenPayload{
		Subject:   m.Subject,
		ExpiresAt: time.Now().Add(1 * time.Hour).Unix(),
	}, nil
}
