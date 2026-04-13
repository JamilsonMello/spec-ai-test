package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/example/cadastro-de-usuarios/internal/application/usecase"
)

type AuthHandler struct {
	LoginUseCase *usecase.LoginUseCase
}

func NewAuthHandler(loginUC *usecase.LoginUseCase) *AuthHandler {
	return &AuthHandler{
		LoginUseCase: loginUC,
	}
}

func (handler *AuthHandler) Login(c echo.Context) error {
	var req usecase.LoginRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request payload"})
	}

	resp, err := handler.LoginUseCase.Execute(req)
	if err != nil {
		if err == usecase.ErrInvalidCredentials {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": err.Error()})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
	}

	return c.JSON(http.StatusOK, resp)
}
