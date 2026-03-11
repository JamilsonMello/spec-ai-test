package domain

import "errors"

var (
	ErrTokenExpired = errors.New("token expired")
	ErrInvalidToken = errors.New("invalid token")
)

type TokenPayload struct {
	Subject   string
	ExpiresAt int64
}

func (payload *TokenPayload) Validate() error {
	if payload.Subject == "" {
		return ErrInvalidToken
	}
	if payload.ExpiresAt == 0 {
		return ErrInvalidToken
	}
	return nil
}
