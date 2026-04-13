package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"time"
)

type JWTProvider struct {
	secret string
}

func NewJWTProvider(secret string) *JWTProvider {
	return &JWTProvider{secret: secret}
}

func (p *JWTProvider) GenerateToken(userID string) (string, error) {
	header := base64.RawURLEncoding.EncodeToString([]byte(`{"alg":"HS256","typ":"JWT"}`))

	claims := map[string]interface{}{
		"sub": userID,
		"exp": time.Now().Add(12 * time.Hour).Unix(),
	}

	claimsJSON, err := json.Marshal(claims)
	if err != nil {
		return "", err
	}

	payload := base64.RawURLEncoding.EncodeToString(claimsJSON)

	message := header + "." + payload
	mac := hmac.New(sha256.New, []byte(p.secret))
	mac.Write([]byte(message))
	signature := base64.RawURLEncoding.EncodeToString(mac.Sum(nil))

	return message + "." + signature, nil
}
