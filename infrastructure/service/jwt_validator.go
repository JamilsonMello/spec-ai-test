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

type JWTValidatorService struct {
	secret string
}

func NewJWTValidatorService() *JWTValidatorService {
	secret := os.Getenv("JWT_SECRET")
	return &JWTValidatorService{secret: secret}
}

func (jWTValidatorService *JWTValidatorService) Validate(token string) (*domain.TokenPayload, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return nil, domain.ErrInvalidToken
	}

	message := parts[0] + "." + parts[1]
	mac := hmac.New(sha256.New, []byte(jWTValidatorService.secret))
	mac.Write([]byte(message))
	expectedSignature := base64.RawURLEncoding.EncodeToString(mac.Sum(nil))

	if parts[2] != expectedSignature {
		return nil, domain.ErrInvalidToken
	}

	payloadStr := parts[1]
	if pad := len(payloadStr) % 4; pad != 0 {
		payloadStr += strings.Repeat("=", 4-pad)
	}

	payloadBytes, decodeErr := base64.URLEncoding.DecodeString(payloadStr)
	if decodeErr != nil {
		return nil, domain.ErrInvalidToken
	}

	var claims struct {
		Sub string `json:"sub"`
		Exp int64  `json:"exp"`
	}

	if unmarshalErr := json.Unmarshal(payloadBytes, &claims); unmarshalErr != nil {
		return nil, domain.ErrInvalidToken
	}

	if time.Now().Unix() >= claims.Exp {
		return nil, domain.ErrTokenExpired
	}

	return &domain.TokenPayload{
		Subject:   claims.Sub,
		ExpiresAt: claims.Exp,
	}, nil
}
