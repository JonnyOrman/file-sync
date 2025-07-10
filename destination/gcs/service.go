package gcs

import (
	"context"
	"io"
	"path"
	"strings"

	"cloud.google.com/go/storage"
	"google.golang.org/api/option"
	storageinterface "github.com/JonnyOrman/cloud-file-backup/destination/interface"
)

// Service implements Google Cloud Storage
type Service struct {
	config    Config
	client    *storage.Client
	bucket    *storage.BucketHandle
}

// NewService creates a new GCS storage service
func NewService(config Config) storageinterface.StorageService {
	return &Service{
		config: config,
	}
}

// GetStorageName returns the storage service name
func (s *Service) GetStorageName() string {
	return "Google Cloud Storage"
}

// Initialize sets up the GCS client and validates the bucket
func (s *Service) Initialize(ctx context.Context) error {
	var err error
	
	// Create client with optional service account key
	if s.config.GCSKeyFile != "" {
		s.client, err = storage.NewClient(ctx, option.WithCredentialsFile(s.config.GCSKeyFile))
	} else {
		// Use default credentials (environment variables, metadata server, etc.)
		s.client, err = storage.NewClient(ctx)
	}
	if err != nil {
		return err
	}
	
	// Get bucket handle
	s.bucket = s.client.Bucket(s.config.GCSBucket)
	
	// Verify bucket exists and is accessible
	_, err = s.bucket.Attrs(ctx)
	return err
}

// WriteFile uploads content to GCS
func (s *Service) WriteFile(ctx context.Context, relativePath string, content io.Reader, size int64) error {
	objectName := s.buildObjectName(relativePath)
	
	// Create object writer
	obj := s.bucket.Object(objectName)
	writer := obj.NewWriter(ctx)
	
	// Copy content to GCS
	_, err := io.Copy(writer, content)
	if err != nil {
		writer.Close()
		return err
	}
	
	// Close writer to finalize upload
	return writer.Close()
}

// FileExists checks if a file exists in GCS
func (s *Service) FileExists(ctx context.Context, relativePath string) (bool, error) {
	objectName := s.buildObjectName(relativePath)
	
	obj := s.bucket.Object(objectName)
	_, err := obj.Attrs(ctx)
	
	if err != nil {
		if err == storage.ErrObjectNotExist {
			return false, nil
		}
		return false, err
	}
	
	return true, nil
}

// CreateDirectory is a no-op for GCS (directories don't exist in GCS)
func (s *Service) CreateDirectory(ctx context.Context, relativePath string) error {
	// GCS doesn't have directories, so this is a no-op
	return nil
}

// buildObjectName constructs the full GCS object name from the relative path
func (s *Service) buildObjectName(relativePath string) string {
	// Normalize path separators to forward slashes for GCS
	relativePath = strings.ReplaceAll(relativePath, "\\", "/")
	
	if s.config.GCSPrefix == "" {
		return relativePath
	}
	
	// Ensure prefix ends with a slash
	prefix := strings.TrimSuffix(s.config.GCSPrefix, "/") + "/"
	return path.Join(prefix, relativePath)
} 