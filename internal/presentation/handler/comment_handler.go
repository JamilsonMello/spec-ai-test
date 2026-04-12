package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/example/cadastro-de-usuarios/internal/application/usecase"
)

type CommentHandler struct {
	CreateCommentUseCase *usecase.CreateCommentUseCase
	ListCommentsUseCase  *usecase.ListCommentsUseCase
}

func NewCommentHandler(createCommentUC *usecase.CreateCommentUseCase, listCommentsUC *usecase.ListCommentsUseCase) *CommentHandler {
	return &CommentHandler{
		CreateCommentUseCase: createCommentUC,
		ListCommentsUseCase:  listCommentsUC,
	}
}

func (h *CommentHandler) CreateComment(c echo.Context) error {
	userID, ok := c.Get("user_id").(string)
	if !ok || userID == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
	}

	postID := c.Param("id")
	if postID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid post ID"})
	}

	var req usecase.CreateCommentRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request payload"})
	}

	req.UserID = userID
	req.PostID = postID

	resp, err := h.CreateCommentUseCase.Execute(req)
	if err != nil {
		statusCode, message := MapErrorToHTTP(err)
		return c.JSON(statusCode, map[string]string{"error": message})
	}

	return c.JSON(http.StatusCreated, resp)
}

func (h *CommentHandler) ListComments(c echo.Context) error {
	postID := c.Param("id")
	if postID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid post ID"})
	}

	resp, err := h.ListCommentsUseCase.Execute(usecase.ListCommentsRequest{PostID: postID})
	if err != nil {
		statusCode, message := MapErrorToHTTP(err)
		return c.JSON(statusCode, map[string]string{"error": message})
	}

	return c.JSON(http.StatusOK, resp)
}
