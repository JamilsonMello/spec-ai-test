package service

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"os"
	"strings"
	"time"

	"github.com/example/cadastro-de-usuarios/domain"
)

// JWTValidatorService validates a JWT token using HS256
type JWTValidatorService struct {
	secret string
}

// NewJWTValidatorService creates a new JWTValidatorService with the secret from environment
func NewJWTValidatorService() *JWTValidatorService {
	secret := os.Getenv("JWT_SECRET")
	return &JWTValidatorService{secret: secret}
}

// Validate decodes the JWT token, checks the signature, and verifies the expiration time
func (s *JWTValidatorService) Validate(token string) (*domain.TokenPayload, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return nil, domain.ErrInvalidToken
	}

	// Calculate and verify signature
	message := parts[0] + "." + parts[1]
	mac := hmac.New(sha256.New, []byte(s.secret))
	mac.Write([]byte(message))
	expectedSignature := base64.RawURLEncoding.EncodeToString(mac.Sum(nil))

	if parts[2] != expectedSignature {
		return nil, domain.ErrInvalidToken
	}

	// Decode payload (allow for both RawURLEncoding and URLEncoding with padding by standardizing it)
	payloadStr := parts[1]
	if pad := len(payloadStr) % 4; pad != 0 {
		payloadStr += strings.Repeat("=", 4-pad)
	}

	payloadBytes, err := base64.URLEncoding.DecodeString(payloadStr)
	if err != nil {
		return nil, domain.ErrInvalidToken
	}

	var claims struct {
		Sub string `json:"sub"`
		Exp int64  `json:"exp"`
	}

	if err := json.Unmarshal(payloadBytes, &claims); err != nil {
		return nil, domain.ErrInvalidToken
	}

	// Check expiration
	if time.Now().Unix() >= claims.Exp {
		return nil, domain.ErrTokenExpired
	}

	return &domain.TokenPayload{
		Subject:   claims.Sub,
		ExpiresAt: claims.Exp,
	}, nil
}