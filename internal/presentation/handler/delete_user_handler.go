package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/example/cadastro-de-usuarios/internal/application/usecase"
)

type DeleteUserHandler struct {
	DeleteUserUseCase *usecase.DeleteUserUseCase
}

func NewDeleteUserHandler(deleteUC *usecase.DeleteUserUseCase) *DeleteUserHandler {
	return &DeleteUserHandler{
		DeleteUserUseCase: deleteUC,
	}
}

func (h *DeleteUserHandler) DeleteUser(c echo.Context) error {
	userRole := c.Request().Header.Get("X-User-Role")
	if userRole == "" {
		userRole = c.QueryParam("role")
	}

	userID := c.Param("id")

	err := h.DeleteUserUseCase.Execute(userID, userRole)
	if err != nil {
		statusCode, message := MapErrorToHTTP(err)
		return c.JSON(statusCode, map[string]string{"error": message})
	}

	return c.NoContent(http.StatusNoContent)
}
