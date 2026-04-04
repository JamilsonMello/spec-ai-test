package service

import (
	"context"
	"mime/multipart"
)

type FileStorageService interface {
	Upload(ctx context.Context, file multipart.File, fileName string) (string, error)
}
