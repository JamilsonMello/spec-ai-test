package storage

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type LocalStorage struct {
	BasePath string
}

func NewLocalStorage(basePath string) *LocalStorage {
	return &LocalStorage{BasePath: basePath}
}

func (s *LocalStorage) Save(file io.Reader, filename string) (string, error) {
	dir := filepath.Join(s.BasePath, "profile")
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", fmt.Errorf("falha ao criar diretório: %w", err)
	}

	fullPath := filepath.Join(dir, filename)

	dst, err := os.Create(fullPath)
	if err != nil {
		return "", fmt.Errorf("falha ao criar arquivo: %w", err)
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		return "", fmt.Errorf("falha ao salvar arquivo: %w", err)
	}

	relativePath := filepath.Join("/uploads/profile", filename)
	return relativePath, nil
}
