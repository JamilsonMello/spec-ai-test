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
	ErrVideoFileTooLarge       = errors.New("Arquivo excede o limite de 50MB")
	ErrUnsupportedVideoFormat  = errors.New("Formato de vídeo não suportado")
	ErrPostNotFoundUpload      = errors.New("Post não encontrado")
	ErrForbiddenUpload         = errors.New("Acesso negado. Você não é o autor deste post.")
	ErrSaveVideoFailed         = errors.New("Falha ao salvar vídeo")
)

const maxVideoFileSize int64 = 50 * 1024 * 1024

type UploadPostVideoRequest struct {
	PostID string
	UserID string
	File   multipart.File
	Header *multipart.FileHeader
}

type UploadPostVideoResponse struct {
	Message  string `json:"message"`
	VideoURL string `json:"video_url"`
}

type UploadPostVideoUseCase struct {
	PostRepository domain.PostRepository
	UploadPath     string
}

func NewUploadPostVideoUseCase(repo domain.PostRepository, uploadPath string) *UploadPostVideoUseCase {
	return &UploadPostVideoUseCase{
		PostRepository: repo,
		UploadPath:     uploadPath,
	}
}

func (uc *UploadPostVideoUseCase) Execute(req UploadPostVideoRequest) (*UploadPostVideoResponse, error) {
	post, err := uc.PostRepository.GetPostByID(req.PostID)
	if err != nil {
		if errors.Is(err, domain.ErrPostNotFound) {
			return nil, ErrPostNotFoundUpload
		}
		return nil, err
	}

	if post.AuthorID != req.UserID {
		return nil, ErrForbiddenUpload
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

	if err := uc.PostRepository.UpdatePostVideo(req.PostID, relativePath); err != nil {
		return nil, ErrSaveVideoFailed
	}

	return &UploadPostVideoResponse{
		Message:  "Video enviado com sucesso",
		VideoURL: relativePath,
	}, nil
}

func (uc *UploadPostVideoUseCase) validateFileSize(header *multipart.FileHeader) error {
	if header.Size > maxVideoFileSize {
		return ErrVideoFileTooLarge
	}
	return nil
}

func (uc *UploadPostVideoUseCase) validateContentType(header *multipart.FileHeader) error {
	contentType := header.Header.Get("Content-Type")
	if contentType != "video/mp4" && contentType != "video/quicktime" {
		return ErrUnsupportedVideoFormat
	}
	return nil
}

func (uc *UploadPostVideoUseCase) validateFileExtension(filename string) (string, error) {
	ext := strings.ToLower(filepath.Ext(filename))
	if ext != ".mp4" && ext != ".mov" {
		return "", ErrUnsupportedVideoFormat
	}
	return ext, nil
}

func (uc *UploadPostVideoUseCase) saveFile(file multipart.File, ext string) (string, error) {
	dir := filepath.Join(uc.UploadPath, "videos")
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", ErrSaveVideoFailed
	}

	fileName := fmt.Sprintf("%s%s", uuid.New().String(), ext)
	fullPath := filepath.Join(dir, fileName)

	dst, err := os.Create(fullPath)
	if err != nil {
		return "", ErrSaveVideoFailed
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		return "", ErrSaveVideoFailed
	}

	relativePath := filepath.Join("/uploads/videos", fileName)
	return relativePath, nil
}
