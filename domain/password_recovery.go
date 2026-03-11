package domain

import (
	"crypto/rand"
	"encoding/hex"
	"time"

	"github.com/google/uuid"
)

type PasswordRecovery struct {
	ID        string    `json:"id"`
	Token     string    `json:"-"`
	UserID    string    `json:"userId"`
	ExpiresAt time.Time `json:"expiresAt"`
	Used      bool      `json:"used"`
	CreatedAt time.Time `json:"createdAt"`
}

func NewPasswordRecovery(userID string) (*PasswordRecovery, error) {
	token, err := generateSecureToken()
	if err != nil {
		return nil, err
	}

	now := time.Now()
	return &PasswordRecovery{
		ID:        uuid.New().String(),
		Token:     token,
		UserID:    userID,
		ExpiresAt: now.Add(24 * time.Hour),
		Used:      false,
		CreatedAt: now,
	}, nil
}

func (passwordRecovery *PasswordRecovery) IsValid() bool {
	return !passwordRecovery.Used && time.Now().Before(passwordRecovery.ExpiresAt)
}

func (passwordRecovery *PasswordRecovery) MarkAsUsed() {
	passwordRecovery.Used = true
}

func generateSecureToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
