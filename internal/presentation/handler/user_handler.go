package handler

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	"github.com/example/cadastro-de-usuarios/internal/application/usecase"
)

type UserHandler struct {
	RegisterUserUseCase            *usecase.RegisterUserUseCase
	ListUsersUseCase               *usecase.ListUsersUseCase
	UpdateUserProfileUseCase       *usecase.UpdateUserProfileUseCase
	DeleteUserUseCase              *usecase.DeleteUserUseCase
	UploadProfilePictureUseCase    *usecase.UploadProfilePictureUseCase
}

func NewUserHandler(registerUC *usecase.RegisterUserUseCase, listUC *usecase.ListUsersUseCase, updateProfileUC *usecase.UpdateUserProfileUseCase, deleteUC *usecase.DeleteUserUseCase, uploadPictureUC *usecase.UploadProfilePictureUseCase) *UserHandler {
	return &UserHandler{
		RegisterUserUseCase:         registerUC,
		ListUsersUseCase:            listUC,
		UpdateUserProfileUseCase:    updateProfileUC,
		DeleteUserUseCase:           deleteUC,
		UploadProfilePictureUseCase: uploadPictureUC,
	}
}

func (handler *UserHandler) RegisterUser(c echo.Context) error {
	var req usecase.RegisterUserRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request payload"})
	}

	resp, err := handler.RegisterUserUseCase.Execute(req)
	if err != nil {
		statusCode, message := MapErrorToHTTP(err)
		return c.JSON(statusCode, map[string]string{"error": message})
	}

	return c.JSON(http.StatusOK, resp)
}

func (handler *UserHandler) ListUsers(c echo.Context) error {
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

	resp, err := handler.ListUsersUseCase.Execute(req)
	if err != nil {
		statusCode, message := MapErrorToHTTP(err)
		return c.JSON(statusCode, map[string]string{"error": message})
	}

	return c.JSON(http.StatusOK, resp)
}

func (handler *UserHandler) DeleteUser(c echo.Context) error {
	userRole := c.Request().Header.Get("X-User-Role")
	if userRole == "" {
		userRole = c.QueryParam("role")
	}

	userID := c.Param("id")

	err := handler.DeleteUserUseCase.Execute(userID, userRole)
	if err != nil {
		statusCode, message := MapErrorToHTTP(err)
		return c.JSON(statusCode, map[string]string{"error": message})
	}

	return c.NoContent(http.StatusNoContent)
}

func (handler *UserHandler) UpdateUserProfile(c echo.Context) error {
	var req usecase.UpdateUserProfileRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request payload"})
	}

	req.UserID = c.Param("id")

	_, err := handler.UpdateUserProfileUseCase.Execute(req)
	if err != nil {
		statusCode, message := MapErrorToHTTP(err)
		return c.JSON(statusCode, map[string]string{"error": message})
	}

	return c.NoContent(http.StatusNoContent)
}

func mapUploadErrorCode(err error) string {
	switch err {
	case usecase.ErrFileTooLarge, usecase.ErrUnsupportedFileFormat:
		return "INVALID_FILE"
	case usecase.ErrUserNotFoundUpload:
		return "USER_NOT_FOUND"
	default:
		return "INTERNAL_ERROR"
	}
}

func (handler *UserHandler) UploadProfilePicture(c echo.Context) error {
	userID := c.Param("id")

	file, header, err := c.Request().FormFile("foto")
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Arquivo inválido ou muito grande", "code": "INVALID_FILE"})
	}
	defer file.Close()

	req := usecase.UploadProfilePictureRequest{
		UserID: userID,
		File:   file,
		Header: header,
	}

	resp, err := handler.UploadProfilePictureUseCase.Execute(req)
	if err != nil {
		statusCode, message := MapErrorToHTTP(err)
		code := mapUploadErrorCode(err)
		return c.JSON(statusCode, map[string]string{"error": message, "code": code})
	}

	return c.JSON(http.StatusOK, resp)
}
