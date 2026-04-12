package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/example/cadastro-de-usuarios/internal/application/usecase"
)

type CommentHandler struct {
	CreateCommentUseCase *usecase.CreateCommentUseCase
}

func NewCommentHandler(createCommentUC *usecase.CreateCommentUseCase) *CommentHandler {
	return &CommentHandler{
		CreateCommentUseCase: createCommentUC,
	}
}

func (h *CommentHandler) CreateComment(c echo.Context) error {
	userID, ok := c.Get("user_id").(string)
	if !ok || userID == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized", "code": "UNAUTHORIZED"})
	}

	postID := c.Param("id")
	if postID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "ID do post inválido", "code": "INVALID_INPUT"})
	}

	var req usecase.CreateCommentRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request payload", "code": "INVALID_INPUT"})
	}

	req.PostID = postID
	req.UserID = userID

	resp, err := h.CreateCommentUseCase.Execute(req)
	if err != nil {
		switch err {
		case usecase.ErrInvalidCommentContent:
			return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error(), "code": "INVALID_INPUT"})
		case usecase.ErrCommentPostNotFound:
			return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error(), "code": "POST_NOT_FOUND"})
		default:
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal server error", "code": "INTERNAL_ERROR"})
		}
	}

	return c.JSON(http.StatusCreated, resp)
}
