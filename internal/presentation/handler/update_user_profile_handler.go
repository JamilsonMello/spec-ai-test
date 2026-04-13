package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/example/cadastro-de-usuarios/internal/application/usecase"
)

type UpdateUserProfileHandler struct {
	UpdateUserProfileUseCase *usecase.UpdateUserProfileUseCase
}

func NewUpdateUserProfileHandler(updateProfileUC *usecase.UpdateUserProfileUseCase) *UpdateUserProfileHandler {
	return &UpdateUserProfileHandler{
		UpdateUserProfileUseCase: updateProfileUC,
	}
}

func (h *UpdateUserProfileHandler) UpdateUserProfile(c echo.Context) error {
	var req usecase.UpdateUserProfileRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request payload"})
	}

	req.UserID = c.Param("id")

	_, err := h.UpdateUserProfileUseCase.Execute(req)
	if err != nil {
		statusCode, message := MapErrorToHTTP(err)
		return c.JSON(statusCode, map[string]string{"error": message})
	}

	return c.NoContent(http.StatusNoContent)
}
