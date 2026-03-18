package domain

import (
	"regexp"
	"time"
)

type User struct {
	ID                string    `json:"id"`
	Name              string    `json:"name"`
	Surname           string    `json:"surname"`
	Email             string    `json:"email"`
	BirthDate         time.Time `json:"birthDate"`
	Password          string    `json:"-"`
	RecoveryToken     string    `json:"-"`
	Role              string    `json:"role"`
	ProfilePictureURL string    `json:"profile_picture_url"`
	CreatedAt         time.Time `json:"createdAt"`
}

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
var nameSurnameRegex = regexp.MustCompile(`^[a-zA-Z\s]+$`)

func (u *User) IsValidName() bool {
	return len(u.Name) >= 2 && len(u.Name) <= 50 && nameSurnameRegex.MatchString(u.Name)
}

func (u *User) IsValidSurname() bool {
	return len(u.Surname) >= 2 && len(u.Surname) <= 50 && nameSurnameRegex.MatchString(u.Surname)
}

func (u *User) IsValidEmailFormat() bool {
	return emailRegex.MatchString(u.Email)
}

func (u *User) IsAdult() bool {
	eighteenYearsAgo := time.Now().AddDate(-18, 0, 0)
	return u.BirthDate.Before(eighteenYearsAgo) || u.BirthDate.Equal(eighteenYearsAgo)
}

func (u *User) IsPastDate() bool {
	return u.BirthDate.Before(time.Now())
}

func (u *User) IsValidPassword(password string) bool {
	return len(password) >= 8
}
