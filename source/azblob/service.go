package azblob

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	cloudinterface "github.com/JonnyOrman/cloud-file-backup/source/interface"
)

// Service implements the CloudService interface for Azure Blob Storage
type Service struct {
	config    Config
	client    *azblob.Client
	container string
	prefix    string
}

// NewService creates a new Azure Blob source service
func NewService(config Config) *Service {
	return &Service{
		config:    config,
		container: config.AzureBlobContainer,
		prefix:    config.AzureBlobPrefix,
	}
}

// GetServiceName returns the name of the service
func (s *Service) GetServiceName() string {
	return "azblob"
}

// Authenticate initializes the Azure Blob client
func (s *Service) Authenticate(ctx context.Context) error {
	// Create credential
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		return fmt.Errorf("failed to create Azure credential: %w", err)
	}

	// Create client
	accountURL := fmt.Sprintf("https://%s.blob.core.windows.net", s.config.AzureBlobAccount)
	client, err := azblob.NewClient(accountURL, cred, nil)
	if err != nil {
		return fmt.Errorf("failed to create Azure Blob client: %w", err)
	}

	s.client = client
	return nil
}

// ListFiles lists all files in the specified path
func (s *Service) ListFiles(ctx context.Context, path string) ([]cloudinterface.CloudFile, error) {
	if s.client == nil {
		return nil, fmt.Errorf("Azure Blob client not initialized")
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

	// List blobs in the container
	pager := s.client.NewListBlobsFlatPager(s.container, &azblob.ListBlobsFlatOptions{
		Prefix: &listPrefix,
	})

	for pager.More() {
		page, err := pager.NextPage(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to list Azure blobs: %w", err)
		}

		for _, blob := range page.Segment.BlobItems {
			// Skip the prefix itself
			if *blob.Name == listPrefix {
				continue
			}

			// Calculate relative path
			relPath := strings.TrimPrefix(*blob.Name, s.prefix)
			relPath = strings.TrimPrefix(relPath, "/")

			// Handle directories (Azure Blob doesn't have real directories, but we can infer them)
			if strings.HasSuffix(*blob.Name, "/") {
				dirPath := strings.TrimSuffix(relPath, "/")
				if !processedDirs[dirPath] {
					dirName := filepath.Base(dirPath)
					if dirName == "" {
						dirName = filepath.Base(strings.TrimSuffix(*blob.Name, "/"))
					}
					
					files = append(files, cloudinterface.CloudFile{
						ID:           dirPath,
						Name:         dirName,
						Size:         0,
						LastModified: *blob.Properties.LastModified,
						ETag:         *blob.Properties.ETag,
						IsFolder:     true,
					})
					processedDirs[dirPath] = true
				}
				continue
			}

			// Handle files
			fileName := filepath.Base(*blob.Name)
			files = append(files, cloudinterface.CloudFile{
				ID:           relPath,
				Name:         fileName,
				Size:         *blob.Properties.ContentLength,
				LastModified: *blob.Properties.LastModified,
				ETag:         *blob.Properties.ETag,
				IsFolder:     false,
			})
		}
	}

	return files, nil
}

// DownloadFile downloads a file from Azure Blob Storage to the specified local path
func (s *Service) DownloadFile(ctx context.Context, fileID, localPath string) error {
	if s.client == nil {
		return fmt.Errorf("Azure Blob client not initialized")
	}

	// Build the blob name
	blobName := s.prefix
	if !strings.HasSuffix(blobName, "/") {
		blobName += "/"
	}
	blobName += fileID

	// Create destination directory if it doesn't exist
	destDir := filepath.Dir(localPath)
	if err := os.MkdirAll(destDir, 0755); err != nil {
		return fmt.Errorf("failed to create destination directory: %w", err)
	}

	// Get the blob from Azure
	blobClient := s.client.ServiceClient().NewContainerClient(s.container).NewBlobClient(blobName)

	// Download the blob
	response, err := blobClient.Download(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to download Azure blob: %w", err)
	}

	// Create the local file
	file, err := os.Create(localPath)
	if err != nil {
		return fmt.Errorf("failed to create local file: %w", err)
	}
	defer file.Close()

	// Copy the content
	_, err = io.Copy(file, response.Body(azblob.RetryReaderOptions{}))
	if err != nil {
		return fmt.Errorf("failed to copy file content: %w", err)
	}

	return nil
} 