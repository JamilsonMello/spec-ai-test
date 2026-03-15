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

func (t *TokenPayload) Validate() error {
	if t.Subject == "" {
		return ErrInvalidToken
	}
	if t.ExpiresAt == 0 {
		return ErrInvalidToken
	}
	return nil
}
