package usecase

import (
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"

	"github.com/example/cadastro-de-usuarios/internal/domain"
)

var (
	ErrFileTooLarge           = errors.New("Arquivo excede o limite de 5MB")
	ErrUnsupportedFileFormat  = errors.New("Formato de arquivo não suportado")
	ErrUserNotFoundUpload     = errors.New("Usuário não encontrado")
	ErrSaveFileFailed         = errors.New("Falha ao salvar imagem")
)

const maxFileSize = 5 * 1024 * 1024

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
	UploadPath     string
}

func NewUploadProfilePictureUseCase(repo domain.UserRepository, uploadPath string) *UploadProfilePictureUseCase {
	return &UploadProfilePictureUseCase{
		UserRepository: repo,
		UploadPath:     uploadPath,
	}
}

func (uc *UploadProfilePictureUseCase) Execute(req UploadProfilePictureRequest) (*UploadProfilePictureResponse, error) {
	if err := uc.validateUser(req.UserID); err != nil {
		return nil, err
	}

	if err := uc.validateFileSize(req.Header); err != nil {
		return nil, err
	}

	if err := uc.validateContentType(req.Header); err != nil {
		return nil, err
	}

	ext, err := uc.validateFileExtension(req.Header.Filename)
	if err != nil {
		return nil, err
	}

	relativePath, err := uc.saveFile(req.File, ext)
	if err != nil {
		return nil, err
	}

	if err := uc.UserRepository.UpdateProfilePictureURL(req.UserID, relativePath); err != nil {
		return nil, ErrSaveFileFailed
	}

	return &UploadProfilePictureResponse{
		Message: "Foto atualizada com sucesso",
		URL:     relativePath,
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

func (uc *UploadProfilePictureUseCase) validateContentType(header *multipart.FileHeader) error {
	contentType := header.Header.Get("Content-Type")
	if contentType != "image/jpeg" && contentType != "image/png" {
		return ErrUnsupportedFileFormat
	}
	return nil
}

func (uc *UploadProfilePictureUseCase) validateFileExtension(filename string) (string, error) {
	ext := strings.ToLower(filepath.Ext(filename))
	if ext != ".jpg" && ext != ".jpeg" && ext != ".png" {
		return "", ErrUnsupportedFileFormat
	}
	return ext, nil
}

func (uc *UploadProfilePictureUseCase) saveFile(file multipart.File, ext string) (string, error) {
	dir := filepath.Join(uc.UploadPath, "profile")
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", ErrSaveFileFailed
	}

	fileName := fmt.Sprintf("%s%s", uuid.New().String(), ext)
	fullPath := filepath.Join(dir, fileName)

	dst, err := os.Create(fullPath)
	if err != nil {
		return "", ErrSaveFileFailed
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		return "", ErrSaveFileFailed
	}

	relativePath := filepath.Join("/uploads/profile", fileName)
	return relativePath, nil
}
