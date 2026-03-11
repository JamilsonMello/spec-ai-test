package handler

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/example/cadastro-de-usuarios/application/usecase"
)

type PasswordRecoveryHandler struct {
	RequestPasswordRecoveryUseCase *usecase.RequestPasswordRecoveryUseCase
	ResetPasswordUseCase           *usecase.ResetPasswordUseCase
}

func NewPasswordRecoveryHandler(recoveryUC *usecase.RequestPasswordRecoveryUseCase, resetUC *usecase.ResetPasswordUseCase) *PasswordRecoveryHandler {
	return &PasswordRecoveryHandler{
		RequestPasswordRecoveryUseCase: recoveryUC,
		ResetPasswordUseCase:           resetUC,
	}
}

func (passwordRecoveryHandler *PasswordRecoveryHandler) RequestPasswordRecovery(c echo.Context) error {
	var req usecase.RequestPasswordRecoveryInput
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request payload"})
	}

	resp, err := passwordRecoveryHandler.RequestPasswordRecoveryUseCase.Execute(req)
	if err != nil {
		if err == usecase.ErrUserNotFound {

			return c.JSON(http.StatusOK, map[string]string{"message": "Se o email existir em nossa base, você receberá instruções de recuperação"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
	}

	return c.JSON(http.StatusOK, resp)
}

func (passwordRecoveryHandler *PasswordRecoveryHandler) ResetPassword(c echo.Context) error {
	var req usecase.ResetPasswordInput
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request payload"})
	}

	resp, err := passwordRecoveryHandler.ResetPasswordUseCase.Execute(req)
	if err != nil {
		if errors.Is(err, usecase.ErrInvalidToken) {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		} else if errors.Is(err, usecase.ErrPasswordMismatch) ||
			errors.Is(err, usecase.ErrPasswordTooShort) {
			return c.JSON(http.StatusUnprocessableEntity, map[string]string{"error": err.Error()})
		} else if errors.Is(err, usecase.ErrUserNotFound) {
			return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
	}

	return c.JSON(http.StatusOK, resp)
}
