package local

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	cloudinterface "github.com/JonnyOrman/cloud-file-backup/source/interface"
)

// Service implements the CloudService interface for local file system
type Service struct {
	config Config
	rootPath string
}

// NewService creates a new local source service
func NewService(config Config) *Service {
	return &Service{
		config: config,
		rootPath: config.LocalPath,
	}
}

// GetServiceName returns the name of the service
func (s *Service) GetServiceName() string {
	return "local"
}

// Authenticate performs authentication (no-op for local filesystem)
func (s *Service) Authenticate(ctx context.Context) error {
	// Verify the root path exists and is accessible
	if _, err := os.Stat(s.rootPath); os.IsNotExist(err) {
		return fmt.Errorf("local path does not exist: %s", s.rootPath)
	}
	return nil
}

// ListFiles lists all files in the specified path
func (s *Service) ListFiles(ctx context.Context, path string) ([]cloudinterface.CloudFile, error) {
	fullPath := filepath.Join(s.rootPath, path)
	
	// Check if path exists
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("path does not exist: %s", fullPath)
	}

	var files []cloudinterface.CloudFile
	
	err := filepath.Walk(fullPath, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip the root directory itself
		if filePath == fullPath {
			return nil
		}

		// Calculate relative path from root
		relPath, err := filepath.Rel(s.rootPath, filePath)
		if err != nil {
			return err
		}

		// Convert to forward slashes for consistency
		relPath = filepath.ToSlash(relPath)

		// Create CloudFile
		file := cloudinterface.CloudFile{
			ID:           relPath,
			Name:         info.Name(),
			Size:         info.Size(),
			LastModified: info.ModTime(),
			ETag:         fmt.Sprintf("%d-%d", info.Size(), info.ModTime().Unix()),
			IsFolder:     info.IsDir(),
		}

		files = append(files, file)

		// If it's a directory, don't walk into it (we'll handle it separately)
		if info.IsDir() {
			return filepath.SkipDir
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to list files: %w", err)
	}

	return files, nil
}

// DownloadFile copies a file from the local source to the specified destination
func (s *Service) DownloadFile(ctx context.Context, fileID, localPath string) error {
	sourcePath := filepath.Join(s.rootPath, fileID)
	
	// Check if source file exists
	if _, err := os.Stat(sourcePath); os.IsNotExist(err) {
		return fmt.Errorf("source file does not exist: %s", sourcePath)
	}

	// Create destination directory if it doesn't exist
	destDir := filepath.Dir(localPath)
	if err := os.MkdirAll(destDir, 0755); err != nil {
		return fmt.Errorf("failed to create destination directory: %w", err)
	}

	// Open source file
	sourceFile, err := os.Open(sourcePath)
	if err != nil {
		return fmt.Errorf("failed to open source file: %w", err)
	}
	defer sourceFile.Close()

	// Create destination file
	destFile, err := os.Create(localPath)
	if err != nil {
		return fmt.Errorf("failed to create destination file: %w", err)
	}
	defer destFile.Close()

	// Copy file contents
	_, err = io.Copy(destFile, sourceFile)
	if err != nil {
		return fmt.Errorf("failed to copy file: %w", err)
	}

	return nil
} 