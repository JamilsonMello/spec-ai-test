package domain

import (
	"regexp"
	"time"
)

type User struct {
	ID            string    `json:"id"`
	Name          string    `json:"name"`
	Surname       string    `json:"surname"`
	Email         string    `json:"email"`
	BirthDate     time.Time `json:"birthDate"`
	Password      string    `json:"-"`
	RecoveryToken string    `json:"-"`
	Role          string    `json:"role"`
	CreatedAt     time.Time `json:"createdAt"`
}

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
var nameSurnameRegex = regexp.MustCompile(`^[a-zA-Z\s]+$`)

func (user *User) IsValidName() bool {
	return len(user.Name) >= 2 && len(user.Name) <= 50 && nameSurnameRegex.MatchString(user.Name)
}

func (user *User) IsValidSurname() bool {
	return len(user.Surname) >= 2 && len(user.Surname) <= 50 && nameSurnameRegex.MatchString(user.Surname)
}

func (user *User) IsValidEmailFormat() bool {
	return emailRegex.MatchString(user.Email)
}

func (user *User) IsAdult() bool {
	eighteenYearsAgo := time.Now().AddDate(-18, 0, 0)
	return user.BirthDate.Before(eighteenYearsAgo) || user.BirthDate.Equal(eighteenYearsAgo)
}

func (user *User) IsPastDate() bool {
	return user.BirthDate.Before(time.Now())
}

func (user *User) IsValidPassword(password string) bool {
	return len(password) >= 8
}
