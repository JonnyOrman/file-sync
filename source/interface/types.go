package cloudinterface

import (
	"context"
	"time"
)

// CloudFile represents a file in cloud storage
type CloudFile struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Size         int64     `json:"size"`
	LastModified time.Time `json:"last_modified"`
	ETag         string    `json:"etag"`
	IsFolder     bool      `json:"is_folder"`
	DownloadURL  string    `json:"download_url,omitempty"`
}

// CloudService defines the interface that all cloud services must implement
type CloudService interface {
	// GetServiceName returns the name of the cloud service
	GetServiceName() string
	
	// Authenticate performs authentication with the cloud service
	Authenticate(ctx context.Context) error
	
	// ListFiles lists all files in the specified path
	ListFiles(ctx context.Context, path string) ([]CloudFile, error)
	
	// DownloadFile downloads a file to the specified local path
	DownloadFile(ctx context.Context, fileID, localPath string) error
}
