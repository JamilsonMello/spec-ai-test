package handler

import (
	"net/http"

	"github.com/example/cadastro-de-usuarios/internal/application/usecase"
	"github.com/example/cadastro-de-usuarios/internal/domain"
)

func MapErrorToHTTP(err error) (int, string) {
	switch err {
	case usecase.ErrFileTooLarge, usecase.ErrUnsupportedFileFormat:
		return http.StatusBadRequest, err.Error()
	case usecase.ErrSaveFileFailed:
		return http.StatusInternalServerError, err.Error()
	case usecase.ErrInvalidName, usecase.ErrInvalidSurname, usecase.ErrInvalidEmail,
		usecase.ErrInvalidBirthDate, usecase.ErrUserTooYoung, usecase.ErrFutureBirthDate,
		usecase.ErrInvalidNameUpdate, usecase.ErrInvalidBirthDateUpdate, usecase.ErrFutureBirthDateUpdate,
		usecase.ErrInvalidUserID, usecase.ErrInvalidContent, usecase.ErrInvalidToken,
		usecase.ErrPasswordMismatch, usecase.ErrPasswordTooShort:
		return http.StatusUnprocessableEntity, err.Error()
	case domain.ErrEmailAlreadyExists, usecase.ErrEmailInUse:
		return http.StatusBadRequest, err.Error()
	case domain.ErrUserNotFound, usecase.ErrUserNotFoundUpdate, usecase.ErrUserNotFound, usecase.ErrUserNotFoundUpload,
		domain.ErrPostNotFound, domain.ErrRecoveryTokenNotFound, usecase.ErrRecoveryTokenNotFound:
		return http.StatusNotFound, err.Error()
	case usecase.ErrUnauthorizedRole, usecase.ErrUnauthorizedCreate, usecase.ErrUnauthorizedUpdate:
		return http.StatusForbidden, err.Error()
	default:
		return http.StatusInternalServerError, "Internal server error"
	}
}
