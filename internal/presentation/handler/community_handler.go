package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/example/cadastro-de-usuarios/internal/application/usecase"
)

type CommunityHandler struct {
	CreateCommunityUseCase *usecase.CreateCommunityUseCase
}

func NewCommunityHandler(createCommunityUC *usecase.CreateCommunityUseCase) *CommunityHandler {
	return &CommunityHandler{
		CreateCommunityUseCase: createCommunityUC,
	}
}

func (h *CommunityHandler) CreateCommunity(c echo.Context) error {
	userID, ok := c.Get("user_id").(string)
	if !ok || userID == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
	}

	var req usecase.CreateCommunityRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request payload", "code": "INVALID_INPUT"})
	}

	req.OwnerID = userID

	resp, err := h.CreateCommunityUseCase.Execute(req)
	if err != nil {
		if err == usecase.ErrCommunityAlreadyExists {
			return c.JSON(http.StatusConflict, map[string]string{"error": "community already exists"})
		}
		if err == usecase.ErrInvalidCommunityName || err == usecase.ErrInvalidCommunityDescription {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error(), "code": "INVALID_INPUT"})
		}
		statusCode, message := MapErrorToHTTP(err)
		return c.JSON(statusCode, map[string]string{"error": message})
	}

	return c.JSON(http.StatusCreated, resp)
}
