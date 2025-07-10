package dropbox

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox"
	"github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox/files"
	"github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox/users"
	cloudinterface "github.com/JonnyOrman/cloud-file-backup/source/interface"
)

// Service implements the Dropbox cloud service
type Service struct {
	config      cloudinterface.DropboxConfig
	filesClient files.Client
	usersClient users.Client
}

// NewService creates a new Dropbox service instance
func NewService(config cloudinterface.DropboxConfig) cloudinterface.CloudService {
	return &Service{
		config: config,
	}
}

// GetServiceName returns the service name
func (s *Service) GetServiceName() string {
	return "Dropbox"
}

// Authenticate performs Dropbox authentication
func (s *Service) Authenticate(ctx context.Context) error {
	if s.config.AccessToken == "" {
		return fmt.Errorf("Dropbox access token is required")
	}
	
	// Create Dropbox clients with access token
	config := dropbox.Config{
		Token: s.config.AccessToken,
	}
	s.filesClient = files.New(config)
	s.usersClient = users.New(config)
	
	// Test the connection by getting account info
	_, err := s.usersClient.GetCurrentAccount()
	if err != nil {
		return fmt.Errorf("failed to authenticate with Dropbox: %w", err)
	}
	
	return nil
}

// ListFiles lists files and folders in the specified path
func (s *Service) ListFiles(ctx context.Context, path string) ([]cloudinterface.CloudFile, error) {
	// Normalize path for Dropbox (must start with /)
	if path == "" || path == "/" {
		path = ""
	} else if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	
	// List folder contents
	arg := files.NewListFolderArg(path)
	result, err := s.filesClient.ListFolder(arg)
	if err != nil {
		return nil, fmt.Errorf("failed to list folder: %w", err)
	}
	
	var cloudFiles []cloudinterface.CloudFile
	
	// Process entries
	for _, entry := range result.Entries {
		cloudFile := s.convertEntry(entry)
		if cloudFile != nil {
			cloudFiles = append(cloudFiles, *cloudFile)
		}
	}
	
	// Handle pagination if there are more results
	for result.HasMore {
		cursor := files.NewListFolderContinueArg(result.Cursor)
		result, err = s.filesClient.ListFolderContinue(cursor)
		if err != nil {
			return nil, fmt.Errorf("failed to continue listing folder: %w", err)
		}
		
		for _, entry := range result.Entries {
			cloudFile := s.convertEntry(entry)
			if cloudFile != nil {
				cloudFiles = append(cloudFiles, *cloudFile)
			}
		}
	}
	
	return cloudFiles, nil
}

// convertEntry converts a Dropbox entry to CloudFile
func (s *Service) convertEntry(entry files.IsMetadata) *cloudinterface.CloudFile {
	switch e := entry.(type) {
	case *files.FileMetadata:
		return &cloudinterface.CloudFile{
			ID:           e.PathLower, // Use path as ID for Dropbox since we need it for downloads
			Name:         e.Name,
			Size:         int64(e.Size),
			LastModified: e.ClientModified,
			ETag:         e.Rev,
			IsFolder:     false,
			DownloadURL:  "", // Dropbox doesn't provide direct download URLs
		}
	case *files.FolderMetadata:
		return &cloudinterface.CloudFile{
			ID:           e.PathLower, // Use path as ID for consistency
			Name:         e.Name,
			Size:         0,
			LastModified: time.Time{}, // Folders don't have modification time in Dropbox
			ETag:         "",
			IsFolder:     true,
			DownloadURL:  "",
		}
	case *files.DeletedMetadata:
		// Skip deleted entries
		return nil
	default:
		// Unknown entry type
		return nil
	}
}

// DownloadFile downloads a file by ID to the local path
func (s *Service) DownloadFile(ctx context.Context, fileID, localPath string) error {
	// For Dropbox, we need the file path, not the ID
	// We'll use the file path which should be stored in the ID field
	filePath := fileID
	
	// Ensure the file path starts with /
	if !strings.HasPrefix(filePath, "/") {
		filePath = "/" + filePath
	}
	
	// Download the file
	arg := files.NewDownloadArg(filePath)
	_, content, err := s.filesClient.Download(arg)
	if err != nil {
		return fmt.Errorf("failed to download file: %w", err)
	}
	defer content.Close()
	
	// Create local file
	file, err := os.Create(localPath)
	if err != nil {
		return fmt.Errorf("failed to create local file: %w", err)
	}
	defer file.Close()
	
	// Copy content to local file
	_, err = io.Copy(file, content)
	if err != nil {
		return fmt.Errorf("failed to write file content: %w", err)
	}
	
	return nil
} 