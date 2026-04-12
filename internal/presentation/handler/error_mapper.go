package handler

import (
	"net/http"

	"github.com/example/cadastro-de-usuarios/internal/application/usecase"
	"github.com/example/cadastro-de-usuarios/internal/domain"
)

func MapErrorToHTTP(err error) (int, string) {
	switch err {
	case usecase.ErrFileTooLarge:
		return http.StatusRequestEntityTooLarge, err.Error()
	case usecase.ErrUnsupportedFileFormat:
		return http.StatusBadRequest, err.Error()
	case usecase.ErrSaveFileFailed:
		return http.StatusInternalServerError, err.Error()
	case usecase.ErrInvalidToken, usecase.ErrPasswordTooShort:
		return http.StatusBadRequest, err.Error()
	case usecase.ErrInvalidName, usecase.ErrInvalidSurname, usecase.ErrInvalidEmail,
		usecase.ErrInvalidBirthDate, usecase.ErrUserTooYoung, usecase.ErrFutureBirthDate,
		usecase.ErrInvalidNameUpdate, usecase.ErrInvalidBirthDateUpdate, usecase.ErrFutureBirthDateUpdate,
		usecase.ErrInvalidUserID, usecase.ErrInvalidContent,
		usecase.ErrPasswordMismatch:
		return http.StatusUnprocessableEntity, err.Error()
	case domain.ErrEmailAlreadyExists, usecase.ErrEmailInUse:
		return http.StatusBadRequest, err.Error()
	case usecase.ErrInvalidCommentContent, usecase.ErrInvalidReactionType:
		return http.StatusBadRequest, err.Error()
	case usecase.ErrCommentNotFoundToggle, domain.ErrCommentNotFound:
		return http.StatusNotFound, err.Error()
	case usecase.ErrUnauthorizedComment, usecase.ErrUnauthorizedReaction:
		return http.StatusUnauthorized, err.Error()
	case usecase.ErrInvalidUUIDFormat:
		return http.StatusBadRequest, err.Error()
	case domain.ErrCommunityNotFound:
		return http.StatusNotFound, err.Error()
	case domain.ErrUserNotFound, usecase.ErrUserNotFoundUpdate, usecase.ErrUserNotFound, usecase.ErrUserNotFoundUpload,
		domain.ErrPostNotFound, domain.ErrRecoveryTokenNotFound:
		return http.StatusNotFound, err.Error()
	case usecase.ErrUnauthorizedRole, usecase.ErrUnauthorizedCreate, usecase.ErrUnauthorizedUpdate:
		return http.StatusForbidden, err.Error()
	case usecase.ErrUnauthorizedCreateProduct:
		return http.StatusUnauthorized, err.Error()
	case usecase.ErrUnauthorizedUpdateProduct:
		return http.StatusForbidden, err.Error()
	case domain.ErrProductNotFound:
		return http.StatusNotFound, err.Error()
	case domain.ErrInvalidProductName, domain.ErrInvalidProductDescription,
		domain.ErrInvalidProductPrice, domain.ErrInvalidProductStock,
		usecase.ErrInvalidProductInput:
		return http.StatusBadRequest, err.Error()
	default:
		return http.StatusInternalServerError, "Internal server error"
	}
}
