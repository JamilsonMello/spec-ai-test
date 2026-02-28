package domain

import "errors"

// Common errors for token validation
var (
	ErrTokenExpired = errors.New("token expired")
	ErrInvalidToken = errors.New("invalid token")
)

// TokenPayload represents the payload data inside a JWT
type TokenPayload struct {
	Subject   string
	ExpiresAt int64
}

// Validate checks if the token payload has the minimum required fields
func (t *TokenPayload) Validate() error {
	if t.Subject == "" {
		return ErrInvalidToken
	}
	if t.ExpiresAt == 0 {
		return ErrInvalidToken
	}
	return nil
}