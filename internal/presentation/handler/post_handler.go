package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/example/cadastro-de-usuarios/internal/application/usecase"
)

type PostHandler struct {
	CreatePostUseCase *usecase.CreatePostUseCase
	UpdatePostUseCase *usecase.UpdatePostUseCase
}

func NewPostHandler(createPostUC *usecase.CreatePostUseCase, updatePostUC *usecase.UpdatePostUseCase) *PostHandler {
	return &PostHandler{
		CreatePostUseCase: createPostUC,
		UpdatePostUseCase: updatePostUC,
	}
}

func (handler *PostHandler) CreatePost(c echo.Context) error {
	userID := c.Request().Header.Get("X-User-ID")
	if userID == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
	}

	var req usecase.CreatePostRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request payload"})
	}

	req.AuthorID = userID

	resp, err := handler.CreatePostUseCase.Execute(req)
	if err != nil {
		statusCode, message := MapErrorToHTTP(err)
		return c.JSON(statusCode, map[string]string{"error": message})
	}

	return c.JSON(http.StatusCreated, resp)
}

func (handler *PostHandler) UpdatePost(c echo.Context) error {
	userID, ok := c.Get("user_id").(string)
	if !ok || userID == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Não autenticado", "code": "UNAUTHORIZED"})
	}

	postID := c.Param("id")
	if postID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "ID do post inválido", "code": "INVALID_ID"})
	}

	var req usecase.UpdatePostRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request payload", "code": "INVALID_REQUEST"})
	}

	req.ID = postID
	req.AuthorID = userID

	resp, err := handler.UpdatePostUseCase.Execute(req)
	if err != nil {
		statusCode, message := MapErrorToHTTP(err)
		return c.JSON(statusCode, map[string]string{"error": message})
	}

	return c.JSON(http.StatusOK, resp)
}
