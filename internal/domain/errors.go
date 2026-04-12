package domain

import "errors"

var (
	ErrUserNotFound              = errors.New("user not found")
	ErrEmailAlreadyExists        = errors.New("email already exists")
	ErrPostNotFound              = errors.New("Post não encontrado")
	ErrRecoveryTokenNotFound     = errors.New("recovery token not found")
	ErrCommentNotFound           = errors.New("Comment not found")
	ErrCommunityAlreadyExists    = errors.New("community already exists")
	ErrCommunityNotFound         = errors.New("community_not_found")
	ErrInvalidProductName        = errors.New("invalid product name")
	ErrInvalidProductDescription = errors.New("invalid product description")
	ErrInvalidProductPrice       = errors.New("invalid product price")
	ErrInvalidProductStock       = errors.New("invalid product stock")
)
