package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/example/cadastro-de-usuarios/internal/application/usecase"
)

type ReactionHandler struct {
	ToggleReactionUseCase *usecase.ToggleCommentReactionUseCase
}

func NewReactionHandler(toggleReactionUC *usecase.ToggleCommentReactionUseCase) *ReactionHandler {
	return &ReactionHandler{
		ToggleReactionUseCase: toggleReactionUC,
	}
}

func (h *ReactionHandler) ToggleReaction(c echo.Context) error {
	userID, ok := c.Get("user_id").(string)
	if !ok || userID == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
	}

	commentID := c.Param("id")
	if commentID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid comment ID"})
	}

	var req usecase.ToggleReactionRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request payload"})
	}

	req.UserID = userID
	req.CommentID = commentID

	resp, err := h.ToggleReactionUseCase.Execute(req)
	if err != nil {
		statusCode, message := MapErrorToHTTP(err)
		return c.JSON(statusCode, map[string]string{"error": message})
	}

	if resp == nil {
		return c.NoContent(http.StatusNoContent)
	}

	return c.JSON(http.StatusCreated, resp)
}
