package handler

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/example/cadastro-de-usuarios/application/usecase"
)

type PostHandler struct {
	CreatePostUseCase *usecase.CreatePostUseCase
}

func NewPostHandler(createPostUC *usecase.CreatePostUseCase) *PostHandler {
	return &PostHandler{
		CreatePostUseCase: createPostUC,
	}
}

func (postHandler *PostHandler) CreatePost(c echo.Context) error {

	userID := c.Request().Header.Get("X-User-ID")
	if userID == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
	}

	var req usecase.CreatePostInput
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request payload"})
	}

	req.AuthorID = userID

	resp, err := postHandler.CreatePostUseCase.Execute(req)
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
