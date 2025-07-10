package s3

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	cloudinterface "github.com/JonnyOrman/cloud-file-backup/source/interface"
)

// Service implements the CloudService interface for AWS S3
type Service struct {
	config cloudinterface.Config
	client *s3.Client
	bucket string
	prefix string
}

// NewService creates a new S3 source service
func NewService(config cloudinterface.Config) *Service {
	return &Service{
		config: config,
		bucket: config.S3Bucket,
		prefix: config.S3Prefix,
	}
}

// GetServiceName returns the name of the service
func (s *Service) GetServiceName() string {
	return "s3"
}

// Authenticate initializes the S3 client
func (s *Service) Authenticate(ctx context.Context) error {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return fmt.Errorf("failed to load AWS config: %w", err)
	}

	s.client = s3.NewFromConfig(cfg)
	return nil
}

// ListFiles lists all files in the specified path
func (s *Service) ListFiles(ctx context.Context, path string) ([]cloudinterface.CloudFile, error) {
	if s.client == nil {
		return nil, fmt.Errorf("S3 client not initialized")
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

	// List objects in S3
	paginator := s3.NewListObjectsV2Paginator(s.client, &s3.ListObjectsV2Input{
		Bucket: aws.String(s.bucket),
		Prefix: aws.String(listPrefix),
	})

	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to list S3 objects: %w", err)
		}

		for _, obj := range page.Contents {
			// Skip the prefix itself
			if *obj.Key == listPrefix {
				continue
			}

			// Calculate relative path
			relPath := strings.TrimPrefix(*obj.Key, s.prefix)
			relPath = strings.TrimPrefix(relPath, "/")

			// Handle directories (S3 doesn't have real directories, but we can infer them)
			if strings.HasSuffix(*obj.Key, "/") {
				dirPath := strings.TrimSuffix(relPath, "/")
				if !processedDirs[dirPath] {
					dirName := filepath.Base(dirPath)
					if dirName == "" {
						dirName = filepath.Base(strings.TrimSuffix(*obj.Key, "/"))
					}
					
					files = append(files, cloudinterface.CloudFile{
						ID:           dirPath,
						Name:         dirName,
						Size:         0,
						LastModified: *obj.LastModified,
						ETag:         strings.Trim(*obj.ETag, "\""),
						IsFolder:     true,
					})
					processedDirs[dirPath] = true
				}
				continue
			}

			// Handle files
			fileName := filepath.Base(*obj.Key)
			files = append(files, cloudinterface.CloudFile{
				ID:           relPath,
				Name:         fileName,
				Size:         obj.Size,
				LastModified: *obj.LastModified,
				ETag:         strings.Trim(*obj.ETag, "\""),
				IsFolder:     false,
			})
		}
	}

	return files, nil
}

// DownloadFile downloads a file from S3 to the specified local path
func (s *Service) DownloadFile(ctx context.Context, fileID, localPath string) error {
	if s.client == nil {
		return fmt.Errorf("S3 client not initialized")
	}

	// Build the S3 key
	s3Key := s.prefix
	if !strings.HasSuffix(s3Key, "/") {
		s3Key += "/"
	}
	s3Key += fileID

	// Create destination directory if it doesn't exist
	destDir := filepath.Dir(localPath)
	if err := os.MkdirAll(destDir, 0755); err != nil {
		return fmt.Errorf("failed to create destination directory: %w", err)
	}

	// Get the object from S3
	result, err := s.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(s3Key),
	})
	if err != nil {
		return fmt.Errorf("failed to get S3 object: %w", err)
	}
	defer result.Body.Close()

	// Create the local file
	file, err := os.Create(localPath)
	if err != nil {
		return fmt.Errorf("failed to create local file: %w", err)
	}
	defer file.Close()

	// Copy the content
	_, err = io.Copy(file, result.Body)
	if err != nil {
		return fmt.Errorf("failed to copy file content: %w", err)
	}

	return nil
} 