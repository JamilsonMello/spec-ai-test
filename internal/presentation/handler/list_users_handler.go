package handler

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	"github.com/example/cadastro-de-usuarios/internal/application/usecase"
)

type ListUsersHandler struct {
	ListUsersUseCase *usecase.ListUsersUseCase
}

func NewListUsersHandler(listUC *usecase.ListUsersUseCase) *ListUsersHandler {
	return &ListUsersHandler{
		ListUsersUseCase: listUC,
	}
}

func (h *ListUsersHandler) ListUsers(c echo.Context) error {
	userRole := c.Request().Header.Get("X-User-Role")
	if userRole == "" {
		userRole = c.QueryParam("role")
	}

	if userRole != "admin" {
		return c.JSON(http.StatusForbidden, map[string]string{"error": "Access denied. Admin role required."})
	}

	req := usecase.ListUsersRequest{
		Name:  c.QueryParam("name"),
		Email: c.QueryParam("email"),
	}

	page := 1
	if pageStr := c.QueryParam("page"); pageStr != "" {
		if parsedPage, err := strconv.Atoi(pageStr); err == nil && parsedPage > 0 {
			page = parsedPage
		}
	}
	req.Page = page

	limit := 30
	if limitStr := c.QueryParam("limit"); limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 {
			limit = parsedLimit
			if limit > 30 {
				limit = 30
			}
		}
	}
	req.Limit = limit

	resp, err := h.ListUsersUseCase.Execute(req)
	if err != nil {
		statusCode, message := MapErrorToHTTP(err)
		return c.JSON(statusCode, map[string]string{"error": message})
	}

	return c.JSON(http.StatusOK, resp)
}
