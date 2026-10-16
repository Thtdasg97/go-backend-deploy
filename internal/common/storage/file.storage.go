package storage

import "go-backend/internal/dto"

// implement storage để lưu file
type FileStorage interface {
	Upload(input dto.UploadInput) (*dto.UploadReturn, error)
	Delete(fileName string) error
}
