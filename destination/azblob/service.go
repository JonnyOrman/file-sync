package azblob

import (
	"context"
	"fmt"
	"io"
	"path"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/JonnyOrman/cloud-file-backup/storage/interface"
)


// Service implements Azure Blob Storage
type Service struct {
	config storageinterface.Config
	client *azblob.Client
}

// NewService creates a new Azure Blob storage service
func NewService(config storageinterface.Config) *Service {
	return &Service{
		config: config,
	}
}

// GetStorageName returns the storage service name
func (s *Service) GetStorageName() string {
	return "Azure Blob Storage"
}

// Initialize sets up the Azure Blob client and validates the container
func (s *Service) Initialize(ctx context.Context) error {
	var err error
	
	// Build service URL
	serviceURL := "https://" + s.config.AzureBlobAccount + ".blob.core.windows.net/"
	
	// Create client with appropriate credentials
	if s.config.AzureBlobKey != "" {
		// Use storage account key
		credential, err := azblob.NewSharedKeyCredential(s.config.AzureBlobAccount, s.config.AzureBlobKey)
		if err != nil {
			return err
		}
		s.client, err = azblob.NewClientWithSharedKeyCredential(serviceURL, credential, nil)
	} else {
		// Use default Azure credentials (managed identity, environment variables, etc.)
		credential, err := azidentity.NewDefaultAzureCredential(nil)
		if err != nil {
			return err
		}
		s.client, err = azblob.NewClient(serviceURL, credential, nil)
	}
	if err != nil {
		return err
	}
	
	// Verify container exists and is accessible
	_, err = s.client.ServiceClient().NewContainerClient(s.config.AzureBlobContainer).GetProperties(ctx, nil)
	return err
}

// WriteFile uploads content to Azure Blob Storage
func (s *Service) WriteFile(ctx context.Context, relativePath string, content io.Reader, size int64) error {
	blobName := s.buildBlobName(relativePath)
	
	// Upload blob directly via client
	_, err := s.client.UploadStream(ctx, s.config.AzureBlobContainer, blobName, content, &azblob.UploadStreamOptions{})
	return err
}

// FileExists checks if a blob exists in Azure Blob Storage
func (s *Service) FileExists(ctx context.Context, relativePath string) (bool, error) {
	blobName := s.buildBlobName(relativePath)
	
	// Try to download stream to check existence (cheap operation)
	_, err := s.client.DownloadStream(ctx, s.config.AzureBlobContainer, blobName, &azblob.DownloadStreamOptions{
		Range: azblob.HTTPRange{Offset: 0, Count: 1}, // Just check the first byte
	})
	if err != nil {
		if respErr, ok := err.(*azcore.ResponseError); ok && respErr.StatusCode == 404 {
			return false, nil
		}
		return false, err
	}
	
	return true, nil
}

// CreateDirectory is a no-op for Azure Blob Storage (directories don't exist)
func (s *Service) CreateDirectory(ctx context.Context, relativePath string) error {
	// Azure Blob Storage doesn't have directories, so this is a no-op
	return nil
}

// buildBlobName constructs the full blob name from the relative path
func (s *Service) buildBlobName(relativePath string) string {
	// Normalize path separators to forward slashes for Azure Blob Storage
	relativePath = strings.ReplaceAll(relativePath, "\\", "/")
	
	if s.config.AzureBlobPrefix == "" {
		return relativePath
	}
	
	// Ensure prefix ends with a slash
	prefix := strings.TrimSuffix(s.config.AzureBlobPrefix, "/") + "/"
	return path.Join(prefix, relativePath)
} 