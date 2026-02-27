package domain

import "errors"

// Repository errors
var (
	ErrUserNotFound        = errors.New("user not found")
	ErrEmailAlreadyExists  = errors.New("email already exists")
	ErrPostNotFound        = errors.New("post not found")
	ErrRecoveryTokenNotFound = errors.New("recovery token not found")
)
