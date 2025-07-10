package s3

import (
	"context"
	"errors"
	"io"
	"path"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	storageinterface "github.com/JonnyOrman/cloud-file-backup/destination/interface"
)

// Service implements AWS S3 storage
type Service struct {
	config   Config
	s3Client *s3.Client
}

// NewService creates a new S3 storage service
func NewService(config Config) storageinterface.StorageService {
	return &Service{
		config: config,
	}
}

// GetStorageName returns the storage service name
func (s *Service) GetStorageName() string {
	return "AWS S3"
}

// Initialize sets up the S3 client and validates the bucket
func (s *Service) Initialize(ctx context.Context) error {
	// Load AWS configuration
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(s.config.S3Region))
	if err != nil {
		return err
	}
	
	// Create S3 client
	s.s3Client = s3.NewFromConfig(cfg)
	
	// Verify bucket exists and is accessible
	_, err = s.s3Client.HeadBucket(ctx, &s3.HeadBucketInput{
		Bucket: aws.String(s.config.S3Bucket),
	})
	
	return err
}

// WriteFile uploads content to S3
func (s *Service) WriteFile(ctx context.Context, relativePath string, content io.Reader, size int64) error {
	key := s.buildKey(relativePath)
	
	_, err := s.s3Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:        aws.String(s.config.S3Bucket),
		Key:           aws.String(key),
		Body:          content,
		ContentLength: aws.Int64(size),
	})
	
	return err
}

// FileExists checks if a file exists in S3
func (s *Service) FileExists(ctx context.Context, relativePath string) (bool, error) {
	key := s.buildKey(relativePath)
	
	_, err := s.s3Client.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(s.config.S3Bucket),
		Key:    aws.String(key),
	})
	
	if err != nil {
		// Check if it's a "not found" type error
		var notFound *types.NotFound
		var noSuchKey *types.NoSuchKey
		if errors.As(err, &notFound) || errors.As(err, &noSuchKey) {
			return false, nil
		}
		return false, err
	}
	
	return true, nil
}

// CreateDirectory is a no-op for S3 (directories don't exist in S3)
func (s *Service) CreateDirectory(ctx context.Context, relativePath string) error {
	// S3 doesn't have directories, so this is a no-op
	return nil
}

// buildKey constructs the full S3 key from the relative path
func (s *Service) buildKey(relativePath string) string {
	// Normalize path separators to forward slashes for S3
	relativePath = strings.ReplaceAll(relativePath, "\\", "/")
	
	if s.config.S3Prefix == "" {
		return relativePath
	}
	
	// Ensure prefix ends with a slash
	prefix := strings.TrimSuffix(s.config.S3Prefix, "/") + "/"
	return path.Join(prefix, relativePath)
} 