package domain

import "errors"

var (
	ErrUserNotFound          = errors.New("user not found")
	ErrEmailAlreadyExists    = errors.New("email already exists")
	ErrPostNotFound          = errors.New("Post não encontrado")
	ErrRecoveryTokenNotFound = errors.New("recovery token not found")
)
