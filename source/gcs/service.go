package gcs

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"cloud.google.com/go/storage"
	"google.golang.org/api/iterator"
	cloudinterface "github.com/JonnyOrman/cloud-file-backup/source/interface"
)

// Service implements the CloudService interface for Google Cloud Storage
type Service struct {
	config cloudinterface.Config
	client *storage.Client
	bucket string
	prefix string
}

// NewService creates a new GCS source service
func NewService(config cloudinterface.Config) *Service {
	return &Service{
		config: config,
		bucket: config.GCSBucket,
		prefix: config.GCSPrefix,
	}
}

// GetServiceName returns the name of the service
func (s *Service) GetServiceName() string {
	return "gcs"
}

// Authenticate initializes the GCS client
func (s *Service) Authenticate(ctx context.Context) error {
	client, err := storage.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("failed to create GCS client: %w", err)
	}

	s.client = client
	return nil
}

// ListFiles lists all files in the specified path
func (s *Service) ListFiles(ctx context.Context, path string) ([]cloudinterface.CloudFile, error) {
	if s.client == nil {
		return nil, fmt.Errorf("GCS client not initialized")
	}

	// Build the prefix for listing
	listPrefix := s.prefix
	if path != "" && path != "/" {
		if !strings.HasSuffix(listPrefix, "/") {
			listPrefix += "/"
		}
		listPrefix += strings.TrimPrefix(path, "/")
	}

	var files []cloudinterface.CloudFile
	processedDirs := make(map[string]bool)

	// List objects in GCS
	bucket := s.client.Bucket(s.bucket)
	query := &storage.Query{
		Prefix: listPrefix,
	}

	it := bucket.Objects(ctx, query)
	for {
		attrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to iterate GCS objects: %w", err)
		}

		// Skip the prefix itself
		if attrs.Name == listPrefix {
			continue
		}

		// Calculate relative path
		relPath := strings.TrimPrefix(attrs.Name, s.prefix)
		relPath = strings.TrimPrefix(relPath, "/")

		// Handle directories (GCS doesn't have real directories, but we can infer them)
		if strings.HasSuffix(attrs.Name, "/") {
			dirPath := strings.TrimSuffix(relPath, "/")
			if !processedDirs[dirPath] {
				dirName := filepath.Base(dirPath)
				if dirName == "" {
					dirName = filepath.Base(strings.TrimSuffix(attrs.Name, "/"))
				}
				
				files = append(files, cloudinterface.CloudFile{
					ID:           dirPath,
					Name:         dirName,
					Size:         0,
					LastModified: attrs.Updated,
					ETag:         attrs.ETag,
					IsFolder:     true,
				})
				processedDirs[dirPath] = true
			}
			continue
		}

		// Handle files
		fileName := filepath.Base(attrs.Name)
		files = append(files, cloudinterface.CloudFile{
			ID:           relPath,
			Name:         fileName,
			Size:         attrs.Size,
			LastModified: attrs.Updated,
			ETag:         attrs.ETag,
			IsFolder:     false,
		})
	}

	return files, nil
}

// DownloadFile downloads a file from GCS to the specified local path
func (s *Service) DownloadFile(ctx context.Context, fileID, localPath string) error {
	if s.client == nil {
		return fmt.Errorf("GCS client not initialized")
	}

	// Build the GCS object name
	objectName := s.prefix
	if !strings.HasSuffix(objectName, "/") {
		objectName += "/"
	}
	objectName += fileID

	// Create destination directory if it doesn't exist
	destDir := filepath.Dir(localPath)
	if err := os.MkdirAll(destDir, 0755); err != nil {
		return fmt.Errorf("failed to create destination directory: %w", err)
	}

	// Get the object from GCS
	bucket := s.client.Bucket(s.bucket)
	obj := bucket.Object(objectName)

	// Create a reader
	reader, err := obj.NewReader(ctx)
	if err != nil {
		return fmt.Errorf("failed to create GCS reader: %w", err)
	}
	defer reader.Close()

	// Create the local file
	file, err := os.Create(localPath)
	if err != nil {
		return fmt.Errorf("failed to create local file: %w", err)
	}
	defer file.Close()

	// Copy the content
	_, err = io.Copy(file, reader)
	if err != nil {
		return fmt.Errorf("failed to copy file content: %w", err)
	}

	return nil
} 