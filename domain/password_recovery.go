package domain

import (
	"crypto/rand"
	"encoding/hex"
	"time"

	"github.com/google/uuid"
)

// PasswordRecovery represents a password recovery token in the system.
type PasswordRecovery struct {
	ID        string    `json:"id"`
	Token     string    `json:"-"` // Sensitive field, not exposed in JSON
	UserID    string    `json:"userId"`
	ExpiresAt time.Time `json:"expiresAt"`
	Used      bool      `json:"used"`
	CreatedAt time.Time `json:"createdAt"`
}

// NewPasswordRecovery creates a new password recovery token with 24-hour expiration.
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

// IsValid checks if the recovery token is valid (not expired and not used).
func (pr *PasswordRecovery) IsValid() bool {
	return !pr.Used && time.Now().Before(pr.ExpiresAt)
}

// MarkAsUsed marks the recovery token as used.
func (pr *PasswordRecovery) MarkAsUsed() {
	pr.Used = true
}

// generateSecureToken generates a cryptographically secure random token.
func generateSecureToken() (string, error) {
	bytes := make([]byte, 32) // 64 character hex string
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
