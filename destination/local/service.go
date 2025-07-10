package local

import (
	"context"
	"io"
	"os"
	"path/filepath"

	storageinterface "github.com/JonnyOrman/cloud-file-backup/destination/interface"
)

// Service implements local file system storage
type Service struct {
	config storageinterface.LocalConfig
}

// NewService creates a new local storage service
func NewService(config storageinterface.LocalConfig) storageinterface.StorageService {
	return &Service{
		config: config,
	}
}

// GetStorageName returns the storage service name
func (s *Service) GetStorageName() string {
	return "Local File System"
}

// Initialize creates the base directory if it doesn't exist
func (s *Service) Initialize(ctx context.Context) error {
	return os.MkdirAll(s.config.LocalPath, 0755)
}

// WriteFile writes content to a local file
func (s *Service) WriteFile(ctx context.Context, relativePath string, content io.Reader, size int64) error {
	fullPath := filepath.Join(s.config.LocalPath, relativePath)
	
	// Ensure directory exists
	dir := filepath.Dir(fullPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	
	// Create and write file
	file, err := os.Create(fullPath)
	if err != nil {
		return err
	}
	defer file.Close()
	
	_, err = io.Copy(file, content)
	return err
}

// FileExists checks if a file exists locally
func (s *Service) FileExists(ctx context.Context, relativePath string) (bool, error) {
	fullPath := filepath.Join(s.config.LocalPath, relativePath)
	_, err := os.Stat(fullPath)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// CreateDirectory creates a directory locally
func (s *Service) CreateDirectory(ctx context.Context, relativePath string) error {
	fullPath := filepath.Join(s.config.LocalPath, relativePath)
	return os.MkdirAll(fullPath, 0755)
} 