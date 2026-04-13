package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGenerateToken_Success(t *testing.T) {
	provider := NewJWTProvider("test-secret")

	token, err := provider.GenerateToken("user-123")

	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	parts := strings.Split(token, ".")
	assert.Equal(t, 3, len(parts))
}

func TestGenerateToken_ValidSignature(t *testing.T) {
	secret := "test-secret"
	provider := NewJWTProvider(secret)

	token, err := provider.GenerateToken("user-123")
	assert.NoError(t, err)

	parts := strings.Split(token, ".")
	message := parts[0] + "." + parts[1]
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(message))
	expectedSig := base64.RawURLEncoding.EncodeToString(mac.Sum(nil))

	assert.Equal(t, expectedSig, parts[2])
}

func TestGenerateToken_ContainsSubjectAndExpiration(t *testing.T) {
	provider := NewJWTProvider("test-secret")

	token, err := provider.GenerateToken("user-456")
	assert.NoError(t, err)

	parts := strings.Split(token, ".")
	payloadBytes, err := base64.RawURLEncoding.DecodeString(parts[1])
	assert.NoError(t, err)

	var claims struct {
		Sub string `json:"sub"`
		Exp int64  `json:"exp"`
	}
	err = json.Unmarshal(payloadBytes, &claims)
	assert.NoError(t, err)

	assert.Equal(t, "user-456", claims.Sub)

	expectedExp := time.Now().Add(12 * time.Hour).Unix()
	assert.InDelta(t, expectedExp, claims.Exp, 5)
}
