package handler

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/example/cadastro-de-usuarios/application/usecase"
)

// PasswordRecoveryHandler handles HTTP requests related to password recovery.
type PasswordRecoveryHandler struct {
	RequestPasswordRecoveryUseCase    *usecase.RequestPasswordRecoveryUseCase
	ResetPasswordUseCase              *usecase.ResetPasswordUseCase
	ResendPasswordRecoveryTokenUseCase *usecase.ResendPasswordRecoveryTokenUseCase
}

// NewPasswordRecoveryHandler creates a new PasswordRecoveryHandler.
func NewPasswordRecoveryHandler(recoveryUC *usecase.RequestPasswordRecoveryUseCase, resetUC *usecase.ResetPasswordUseCase, resendUC *usecase.ResendPasswordRecoveryTokenUseCase) *PasswordRecoveryHandler {
	return &PasswordRecoveryHandler{
		RequestPasswordRecoveryUseCase:    recoveryUC,
		ResetPasswordUseCase:              resetUC,
		ResendPasswordRecoveryTokenUseCase: resendUC,
	}
}

// RequestPasswordRecovery handles the POST /password-recovery request.
func (h *PasswordRecoveryHandler) RequestPasswordRecovery(c echo.Context) error {
	var req usecase.RequestPasswordRecoveryRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request payload"})
	}

	resp, err := h.RequestPasswordRecoveryUseCase.Execute(req)
	if err != nil {
		if err == usecase.ErrUserNotFound {
			// Return generic message to avoid email enumeration
			return c.JSON(http.StatusOK, map[string]string{"message": "Se o email existir em nossa base, você receberá instruções de recuperação"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
	}

	return c.JSON(http.StatusOK, resp)
}

// ResendToken handles the POST /api/password/resend-token request.
func (h *PasswordRecoveryHandler) ResendToken(c echo.Context) error {
	var req usecase.ResendTokenRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request payload"})
	}

	if err := h.ResendPasswordRecoveryTokenUseCase.Execute(req); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Se o email existir em nossa base, você receberá um novo token de recuperação"})
}

// ResetPassword handles the POST /password-recovery/reset request.
func (h *PasswordRecoveryHandler) ResetPassword(c echo.Context) error {
	var req usecase.ResetPasswordRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request payload"})
	}

	resp, err := h.ResetPasswordUseCase.Execute(req)
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
