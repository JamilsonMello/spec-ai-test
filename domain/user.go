package domain

import (
	"regexp"
	"time"
)

// User represents a registered user in the system.
type User struct {
	ID            string    `json:"id"`
	Name          string    `json:"name"`
	Email         string    `json:"email"`
	Password      string    `json:"-"` // Omitted from JSON response
	RecoveryToken string    `json:"-"` // Omitted from JSON response
	Role          string    `json:"role"`
	CreatedAt     time.Time `json:"createdAt"`
}

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
var nameRegex = regexp.MustCompile(`^[a-zA-Z\s]+$`)

// IsValidName checks if the name is valid.
func (u *User) IsValidName() bool {
	return len(u.Name) >= 2 && len(u.Name) <= 50 && nameRegex.MatchString(u.Name)
}

// IsValidEmailFormat checks if the email format is valid.
func (u *User) IsValidEmailFormat() bool {
	return emailRegex.MatchString(u.Email)
}
