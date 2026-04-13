package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/example/cadastro-de-usuarios/internal/application/usecase"
)

type UploadProfilePictureHandler struct {
	UploadProfilePictureUseCase *usecase.UploadProfilePictureUseCase
}

func NewUploadProfilePictureHandler(uploadPictureUC *usecase.UploadProfilePictureUseCase) *UploadProfilePictureHandler {
	return &UploadProfilePictureHandler{
		UploadProfilePictureUseCase: uploadPictureUC,
	}
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

func (h *UploadProfilePictureHandler) UploadProfilePicture(c echo.Context) error {
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

	resp, err := h.UploadProfilePictureUseCase.Execute(req)
	if err != nil {
		statusCode, message := MapErrorToHTTP(err)
		code := mapUploadErrorCode(err)
		return c.JSON(statusCode, map[string]string{"error": message, "code": code})
	}

	return c.JSON(http.StatusOK, resp)
}
