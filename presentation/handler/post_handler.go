package handler

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/example/cadastro-de-usuarios/application/usecase"
)

// PostHandler handles HTTP requests related to posts.
type PostHandler struct {
	CreatePostUseCase *usecase.CreatePostUseCase
}

// NewPostHandler creates a new PostHandler.
func NewPostHandler(createPostUC *usecase.CreatePostUseCase) *PostHandler {
	return &PostHandler{
		CreatePostUseCase: createPostUC,
	}
}

// CreatePost handles the POST /posts request.
func (h *PostHandler) CreatePost(c echo.Context) error {
	// Extract authenticated user ID from header
	userID := c.Request().Header.Get("X-User-ID")
	if userID == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
	}

	var req usecase.CreatePostRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request payload"})
	}

	// Set author ID from authentication context
	req.AuthorID = userID

	resp, err := h.CreatePostUseCase.Execute(req)
	if err != nil {
		if errors.Is(err, usecase.ErrInvalidContent) {
			return c.JSON(http.StatusUnprocessableEntity, map[string]string{"error": err.Error()})
		}
		if errors.Is(err, usecase.ErrUnauthorizedCreate) {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": err.Error()})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
	}

	return c.JSON(http.StatusCreated, resp)
}
