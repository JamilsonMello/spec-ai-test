package usecase

import (
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"strings"

	"github.com/google/uuid"

	"github.com/example/cadastro-de-usuarios/internal/domain"
)

var (
	ErrFileTooLarge          = errors.New("Arquivo muito grande")
	ErrUnsupportedFileFormat = errors.New("Formato de arquivo inválido")
	ErrUserNotFoundUpload    = errors.New("Usuário não encontrado")
	ErrSaveFileFailed        = errors.New("Falha ao salvar imagem")
)

const maxFileSize = 2 * 1024 * 1024

type UploadProfilePictureRequest struct {
	UserID string
	File   multipart.File
	Header *multipart.FileHeader
}

type UploadProfilePictureResponse struct {
	Message string `json:"message"`
	URL     string `json:"url"`
}

type UploadProfilePictureUseCase struct {
	UserRepository domain.UserRepository
	FileStorage    domain.FileStorage
}

func NewUploadProfilePictureUseCase(repo domain.UserRepository, fileStorage domain.FileStorage) *UploadProfilePictureUseCase {
	return &UploadProfilePictureUseCase{
		UserRepository: repo,
		FileStorage:    fileStorage,
	}
}

func (uc *UploadProfilePictureUseCase) Execute(req UploadProfilePictureRequest) (*UploadProfilePictureResponse, error) {
	if err := uc.validateUser(req.UserID); err != nil {
		return nil, err
	}

	if err := uc.validateFileSize(req.Header); err != nil {
		return nil, err
	}

	contentType, err := uc.detectContentType(req.File)
	if err != nil {
		return nil, err
	}

	ext := uc.extensionFromContentType(contentType)
	if ext == "" {
		return nil, ErrUnsupportedFileFormat
	}

	filename := fmt.Sprintf("%s%s", uuid.New().String(), ext)

	url, err := uc.FileStorage.Save(req.File, filename)
	if err != nil {
		return nil, ErrSaveFileFailed
	}

	if err := uc.UserRepository.UpdateProfilePictureURL(req.UserID, url); err != nil {
		return nil, ErrSaveFileFailed
	}

	return &UploadProfilePictureResponse{
		Message: "Foto atualizada com sucesso",
		URL:     url,
	}, nil
}

func (uc *UploadProfilePictureUseCase) validateUser(userID string) error {
	if userID == "" {
		return ErrUserNotFoundUpload
	}
	_, err := uc.UserRepository.FindUserByUuid(userID)
	if err != nil {
		return ErrUserNotFoundUpload
	}
	return nil
}

func (uc *UploadProfilePictureUseCase) validateFileSize(header *multipart.FileHeader) error {
	if header.Size > maxFileSize {
		return ErrFileTooLarge
	}
	return nil
}

func (uc *UploadProfilePictureUseCase) detectContentType(file multipart.File) (string, error) {
	buf := make([]byte, 512)
	n, err := file.Read(buf)
	if err != nil && err != io.EOF {
		return "", ErrUnsupportedFileFormat
	}

	if seeker, ok := file.(io.Seeker); ok {
		seeker.Seek(0, io.SeekStart)
	}

	contentType := http.DetectContentType(buf[:n])
	return contentType, nil
}

func (uc *UploadProfilePictureUseCase) extensionFromContentType(contentType string) string {
	ct := strings.ToLower(contentType)
	switch {
	case strings.HasPrefix(ct, "image/jpeg"):
		return ".jpg"
	case strings.HasPrefix(ct, "image/png"):
		return ".png"
	default:
		return ""
	}
}

