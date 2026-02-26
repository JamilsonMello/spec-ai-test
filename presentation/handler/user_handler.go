package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	"github.com/example/cadastro-de-usuarios/infrastructure/repository"
	"github.com/example/cadastro-de-usuarios/application/usecase"
)

// UserHandler handles HTTP requests related to users.
type UserHandler struct {
	RegisterUserUseCase      *usecase.RegisterUserUseCase
	ListUsersUseCase         *usecase.ListUsersUseCase
	UpdateUserProfileUseCase *usecase.UpdateUserProfileUseCase
	DeleteUserUseCase        *usecase.DeleteUserUseCase
}

// NewUserHandler creates a new UserHandler.
func NewUserHandler(registerUC *usecase.RegisterUserUseCase, listUC *usecase.ListUsersUseCase, updateProfileUC *usecase.UpdateUserProfileUseCase, deleteUC *usecase.DeleteUserUseCase) *UserHandler {
	return &UserHandler{
		RegisterUserUseCase:      registerUC,
		ListUsersUseCase:         listUC,
		UpdateUserProfileUseCase: updateProfileUC,
		DeleteUserUseCase:        deleteUC,
	}
}

// RegisterUser handles the POST /usuarios request.
func (h *UserHandler) RegisterUser(c echo.Context) error {
	var req usecase.RegisterUserRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request payload"})
	}

	resp, err := h.RegisterUserUseCase.Execute(req)
	if err != nil {
		if errors.Is(err, usecase.ErrInvalidName) ||
			errors.Is(err, usecase.ErrInvalidSurname) ||
			errors.Is(err, usecase.ErrInvalidEmail) ||
			errors.Is(err, usecase.ErrInvalidBirthDate) ||
			errors.Is(err, usecase.ErrUserTooYoung) ||
			errors.Is(err, usecase.ErrFutureBirthDate) {
			return c.JSON(http.StatusUnprocessableEntity, map[string]string{"error": err.Error()})
		} else if errors.Is(err, repository.ErrEmailAlreadyExists) || errors.Is(err, usecase.ErrEmailInUse) {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
	}

	return c.JSON(http.StatusOK, resp)
}

// ListUsers handles the GET /usuarios/listar request.
// This endpoint requires admin role for access.
func (h *UserHandler) ListUsers(c echo.Context) error {
	// Check for admin role (basic auth check via header or query param for demo)
	// In production, this should use proper JWT/session validation
	userRole := c.Request().Header.Get("X-User-Role")
	if userRole == "" {
		userRole = c.QueryParam("role")
	}

	// Only allow admin users to access this endpoint
	if userRole != "admin" {
		return c.JSON(http.StatusForbidden, map[string]string{"error": "Access denied. Admin role required."})
	}

	// Parse query parameters
	req := usecase.ListUsersRequest{
		Name:  c.QueryParam("name"),
		Email: c.QueryParam("email"),
	}

	// Parse page (default to 1)
	page := 1
	if pageStr := c.QueryParam("page"); pageStr != "" {
		if parsedPage, err := strconv.Atoi(pageStr); err == nil && parsedPage > 0 {
			page = parsedPage
		}
	}
	req.Page = page

	// Parse limit (default to 30, max 30)
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

	// Execute use case
	resp, err := h.ListUsersUseCase.Execute(req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
	}

	return c.JSON(http.StatusOK, resp)
}

// DeleteUser handles the DELETE /usuarios/:id request.
// This endpoint requires admin role for access.
func (h *UserHandler) DeleteUser(c echo.Context) error {
	// Check for admin role (basic auth check via header or query param for demo)
	// In production, this should use proper JWT/session validation
	userRole := c.Request().Header.Get("X-User-Role")
	if userRole == "" {
		userRole = c.QueryParam("role")
	}

	// Get user ID from path parameter
	userID := c.Param("id")

	// Execute use case
	err := h.DeleteUserUseCase.Execute(userID, userRole)
	if err != nil {
		if errors.Is(err, usecase.ErrInvalidUserID) {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		} else if errors.Is(err, usecase.ErrUserNotFound) {
			return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
		} else if errors.Is(err, usecase.ErrUnauthorizedRole) {
			return c.JSON(http.StatusForbidden, map[string]string{"error": err.Error()})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
	}

	return c.NoContent(http.StatusNoContent)
}

// UpdateUserProfile handles the PUT /usuarios/:id request.
func (h *UserHandler) UpdateUserProfile(c echo.Context) error {
	var req usecase.UpdateUserProfileRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request payload"})
	}

	// Get user ID from path parameter
	req.UserID = c.Param("id")

	resp, err := h.UpdateUserProfileUseCase.Execute(req)
	if err != nil {
		if errors.Is(err, usecase.ErrInvalidNameUpdate) ||
			errors.Is(err, usecase.ErrInvalidBirthDateUpdate) ||
			errors.Is(err, usecase.ErrFutureBirthDateUpdate) {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		}
		if errors.Is(err, usecase.ErrUserNotFoundUpdate) {
			return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
	}

	return c.NoContent(http.StatusNoContent)
}
