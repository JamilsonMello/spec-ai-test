package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	"github.com/example/cadastro-de-usuarios/application/usecase"
	"github.com/example/cadastro-de-usuarios/domain"
)

type UserHandler struct {
	RegisterUserUseCase      *usecase.RegisterUserUseCase
	ListUsersUseCase         *usecase.ListUsersUseCase
	UpdateUserProfileUseCase *usecase.UpdateUserProfileUseCase
	DeleteUserUseCase        *usecase.DeleteUserUseCase
}

func NewUserHandler(registerUC *usecase.RegisterUserUseCase, listUC *usecase.ListUsersUseCase, updateProfileUC *usecase.UpdateUserProfileUseCase, deleteUC *usecase.DeleteUserUseCase) *UserHandler {
	return &UserHandler{
		RegisterUserUseCase:      registerUC,
		ListUsersUseCase:         listUC,
		UpdateUserProfileUseCase: updateProfileUC,
		DeleteUserUseCase:        deleteUC,
	}
}

func (userHandler *UserHandler) getRole(c echo.Context) string {
	role := c.Request().Header.Get("X-User-Role")
	if role == "" {
		role = c.QueryParam("role")
	}
	return role
}

func parsePage(c echo.Context) int {
	page := 1
	if pageStr := c.QueryParam("page"); pageStr != "" {
		if parsedPage, err := strconv.Atoi(pageStr); err == nil && parsedPage > 0 {
			page = parsedPage
		}
	}
	return page
}

func parseLimit(c echo.Context) int {
	limit := 30
	if limitStr := c.QueryParam("limit"); limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 {
			limit = parsedLimit
			if limit > 30 {
				limit = 30
			}
		}
	}
	return limit
}

func (userHandler *UserHandler) RegisterUser(c echo.Context) error {
	var req usecase.RegisterUserInput
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request payload"})
	}

	resp, err := userHandler.RegisterUserUseCase.Execute(req)
	if err != nil {
		if errors.Is(err, usecase.ErrInvalidName) ||
			errors.Is(err, usecase.ErrInvalidSurname) ||
			errors.Is(err, usecase.ErrInvalidEmail) ||
			errors.Is(err, usecase.ErrInvalidBirthDate) ||
			errors.Is(err, usecase.ErrUserTooYoung) ||
			errors.Is(err, usecase.ErrFutureBirthDate) {
			return c.JSON(http.StatusUnprocessableEntity, map[string]string{"error": err.Error()})
		} else if errors.Is(err, domain.ErrEmailAlreadyExists) || errors.Is(err, usecase.ErrEmailInUse) {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
	}

	return c.JSON(http.StatusOK, resp)
}

func (userHandler *UserHandler) ListUsers(c echo.Context) error {
	if userHandler.getRole(c) != "admin" {
		return c.JSON(http.StatusForbidden, map[string]string{"error": "Access denied. Admin role required."})
	}

	req := usecase.ListUsersInput{
		Name:  c.QueryParam("name"),
		Email: c.QueryParam("email"),
		Page:  parsePage(c),
		Limit: parseLimit(c),
	}

	resp, err := userHandler.ListUsersUseCase.Execute(req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
	}

	return c.JSON(http.StatusOK, resp)
}

func (userHandler *UserHandler) DeleteUser(c echo.Context) error {

	userRole := c.Request().Header.Get("X-User-Role")
	if userRole == "" {
		userRole = c.QueryParam("role")
	}

	userID := c.Param("id")

	err := userHandler.DeleteUserUseCase.Execute(userID, userRole)
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

func (userHandler *UserHandler) UpdateUserProfile(c echo.Context) error {
	var req usecase.UpdateUserProfileInput
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request payload"})
	}

	req.UserID = c.Param("id")

	_, err := userHandler.UpdateUserProfileUseCase.Execute(req)
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
