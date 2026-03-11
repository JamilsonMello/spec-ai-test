package handler

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/example/cadastro-de-usuarios/application/usecase"
	"github.com/example/cadastro-de-usuarios/domain"
)

// PostHandler handles HTTP requests related to posts.
type PostHandler struct {
	CreatePostUseCase *usecase.CreatePostUseCase
	UpdatePostUseCase *usecase.UpdatePostUseCase
}

// NewPostHandler creates a new PostHandler.
func NewPostHandler(createPostUC *usecase.CreatePostUseCase, updatePostUC *usecase.UpdatePostUseCase) *PostHandler {
	return &PostHandler{
		CreatePostUseCase: createPostUC,
		UpdatePostUseCase: updatePostUC,
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

// UpdatePost handles the PUT /posts/{id} request.
func (h *PostHandler) UpdatePost(c echo.Context) error {
	// Extract authenticated user ID from context (set by AuthMiddleware)
	userID, ok := c.Get("user_id").(string)
	if !ok || userID == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Não autenticado", "code": "UNAUTHORIZED"})
	}

	// Extract post ID from path parameter
	postID := c.Param("id")
	if postID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "ID do post inválido", "code": "INVALID_ID"})
	}

	var req usecase.UpdatePostRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request payload", "code": "INVALID_REQUEST"})
	}

	// Set IDs
	req.ID = postID
	req.AuthorID = userID

	resp, err := h.UpdatePostUseCase.Execute(req)
	if err != nil {
		if errors.Is(err, domain.ErrPostNotFound) {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Post não encontrado", "code": "POST_NOT_FOUND"})
		}
		if errors.Is(err, usecase.ErrUnauthorizedUpdate) {
			return c.JSON(http.StatusForbidden, map[string]string{"error": err.Error(), "code": "FORBIDDEN_ACCESS"})
		}
		if errors.Is(err, usecase.ErrInvalidContent) {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Content inválido", "code": "INVALID_CONTENT"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Erro interno do servidor", "code": "INTERNAL_SERVER_ERROR"})
	}

	return c.JSON(http.StatusOK, resp)
}
