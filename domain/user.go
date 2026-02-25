package domain

import (
	"regexp"
	"time"
)

// User represents a registered user in the system.
type User struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Surname     string    `json:"surname"`
	Email       string    `json:"email"`
	BirthDate   time.Time `json:"birthDate"`
}

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
var nameSurnameRegex = regexp.MustCompile(`^[a-zA-Z\s]+$`)

// IsValidName checks if the name is valid.
func (u *User) IsValidName() bool {
	return len(u.Name) >= 2 && len(u.Name) <= 50 && nameSurnameRegex.MatchString(u.Name)
}

// IsValidSurname checks if the surname is valid.
func (u *User) IsValidSurname() bool {
	return len(u.Surname) >= 2 && len(u.Surname) <= 50 && nameSurnameRegex.MatchString(u.Surname)
}

// IsValidEmailFormat checks if the email format is valid.
func (u *User) IsValidEmailFormat() bool {
	return emailRegex.MatchString(u.Email)
}

// IsAdult checks if the user is 18 years or older.
func (u *User) IsAdult() bool {
	eighteenYearsAgo := time.Now().AddDate(-18, 0, 0)
	return u.BirthDate.Before(eighteenYearsAgo) || u.BirthDate.Equal(eighteenYearsAgo)
}

// IsPastDate checks if the birth date is in the past.
func (u *User) IsPastDate() bool {
	return u.BirthDate.Before(time.Now())
}
