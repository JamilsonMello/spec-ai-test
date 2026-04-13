package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/example/cadastro-de-usuarios/internal/application/usecase"
)

type RegisterUserHandler struct {
	RegisterUserUseCase *usecase.RegisterUserUseCase
}

func NewRegisterUserHandler(registerUC *usecase.RegisterUserUseCase) *RegisterUserHandler {
	return &RegisterUserHandler{
		RegisterUserUseCase: registerUC,
	}
}

func (h *RegisterUserHandler) RegisterUser(c echo.Context) error {
	var req usecase.RegisterUserRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request payload"})
	}

	resp, err := h.RegisterUserUseCase.Execute(req)
	if err != nil {
		statusCode, message := MapErrorToHTTP(err)
		return c.JSON(statusCode, map[string]string{"error": message})
	}

	return c.JSON(http.StatusCreated, resp)
}
