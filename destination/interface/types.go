package storageinterface

import (
	"context"
	"io"
)

// StorageService defines the interface that all storage services must implement
type StorageService interface {
	// GetStorageName returns the name of the storage service
	GetStorageName() string
	
	// Initialize initializes the storage service
	Initialize(ctx context.Context) error
	
	// WriteFile writes a file to the storage service
	WriteFile(ctx context.Context, relativePath string, content io.Reader, size int64) error
	
	// FileExists checks if a file exists in the storage service
	FileExists(ctx context.Context, relativePath string) (bool, error)
	
	// CreateDirectory creates a directory in the storage service
	CreateDirectory(ctx context.Context, relativePath string) error
}
