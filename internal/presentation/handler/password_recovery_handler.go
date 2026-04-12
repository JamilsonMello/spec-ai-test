package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/example/cadastro-de-usuarios/internal/application/usecase"
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

func (handler *PasswordRecoveryHandler) RequestPasswordRecovery(c echo.Context) error {
	var req usecase.RequestPasswordRecoveryRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request payload"})
	}

	resp, err := handler.RequestPasswordRecoveryUseCase.Execute(req)
	if err != nil {
		return c.JSON(http.StatusOK, map[string]string{"message": "Se um usuário com esse e-mail existir, você receberá instruções em breve."})
	}

	return c.JSON(http.StatusOK, resp)
}

func (handler *PasswordRecoveryHandler) ResetPassword(c echo.Context) error {
	var req usecase.ResetPasswordRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Token inválido, expirado ou já utilizado."})
	}

	resp, err := handler.ResetPasswordUseCase.Execute(req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Token inválido, expirado ou já utilizado."})
	}

	return c.JSON(http.StatusOK, resp)
}
