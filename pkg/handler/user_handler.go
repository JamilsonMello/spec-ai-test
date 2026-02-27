package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	user_usecase "github.com/example/cadastro-de-usuarios/pkg/usecase/user"
)

// UserHandler handles HTTP requests related to users.
type UserHandler struct {
	RegisterUserUseCase      *user_usecase.RegisterUserUseCase
	ListUsersUseCase         *user_usecase.ListUsersUseCase
	UpdateUserProfileUseCase *user_usecase.UpdateUserProfileUseCase
	DeleteUserUseCase        *user_usecase.DeleteUserUseCase
}

// NewUserHandler creates a new UserHandler.
func NewUserHandler(registerUC *user_usecase.RegisterUserUseCase, listUC *user_usecase.ListUsersUseCase, updateProfileUC *user_usecase.UpdateUserProfileUseCase, deleteUC *user_usecase.DeleteUserUseCase) *UserHandler {
	return &UserHandler{
		RegisterUserUseCase:      registerUC,
		ListUsersUseCase:         listUC,
		UpdateUserProfileUseCase: updateProfileUC,
		DeleteUserUseCase:        deleteUC,
	}
}

// RegisterUser handles the POST /usuarios request.
func (h *UserHandler) RegisterUser(c echo.Context) error {
	var req user_usecase.RegisterUserRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request payload"})
	}

	resp, err := h.RegisterUserUseCase.Execute(req)
	if err != nil {
		if errors.Is(err, user_usecase.ErrInvalidName) ||
			errors.Is(err, user_usecase.ErrInvalidSurname) ||
			errors.Is(err, user_usecase.ErrInvalidEmail) ||
			errors.Is(err, user_usecase.ErrInvalidBirthDate) ||
			errors.Is(err, user_usecase.ErrUserTooYoung) ||
			errors.Is(err, user_usecase.ErrFutureBirthDate) {
			return c.JSON(http.StatusUnprocessableEntity, map[string]string{"error": err.Error()})
		} else if errors.Is(err, user_usecase.ErrEmailInUse) {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
	}

	return c.JSON(http.StatusOK, resp)
}

// ListUsers handles the GET /usuarios/listar request.
func (h *UserHandler) ListUsers(c echo.Context) error {
	userRole := c.Request().Header.Get("X-User-Role")
	if userRole == "" {
		userRole = c.QueryParam("role")
	}

	if userRole != "admin" {
		return c.JSON(http.StatusForbidden, map[string]string{"error": "Access denied. Admin role required."})
	}

	req := user_usecase.ListUsersRequest{
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
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
	}

	return c.JSON(http.StatusOK, resp)
}

// DeleteUser handles the DELETE /usuarios/:id request.
func (h *UserHandler) DeleteUser(c echo.Context) error {
	userRole := c.Request().Header.Get("X-User-Role")
	if userRole == "" {
		userRole = c.QueryParam("role")
	}

	userID := c.Param("id")

	err := h.DeleteUserUseCase.Execute(userID, userRole)
	if err != nil {
		if errors.Is(err, user_usecase.ErrInvalidUserID) {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		} else if errors.Is(err, user_usecase.ErrUserNotFound) {
			return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
		} else if errors.Is(err, user_usecase.ErrUnauthorizedRole) {
			return c.JSON(http.StatusForbidden, map[string]string{"error": err.Error()})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
	}

	return c.NoContent(http.StatusNoContent)
}

// UpdateUserProfile handles the PUT /usuarios/:id request.
func (h *UserHandler) UpdateUserProfile(c echo.Context) error {
	var req user_usecase.UpdateUserProfileRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request payload"})
	}

	req.UserID = c.Param("id")

	_, err := h.UpdateUserProfileUseCase.Execute(req)
	if err != nil {
		if errors.Is(err, user_usecase.ErrInvalidNameUpdate) ||
			errors.Is(err, user_usecase.ErrInvalidBirthDateUpdate) ||
			errors.Is(err, user_usecase.ErrFutureBirthDateUpdate) {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		}
		if errors.Is(err, user_usecase.ErrUserNotFoundUpdate) {
			return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
	}

	return c.NoContent(http.StatusNoContent)
}
