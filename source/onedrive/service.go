package onedrive

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	cloudinterface "github.com/JonnyOrman/cloud-file-backup/source/interface"
)

// Service implements the OneDrive cloud service
type Service struct {
	config cloudinterface.OneDriveConfig
	token  string
}

// DriveItem represents a OneDrive item from the API
type DriveItem struct {
	ID               string    `json:"id"`
	Name             string    `json:"name"`
	Size             *int64    `json:"size,omitempty"`
	LastModifiedTime time.Time `json:"lastModifiedDateTime"`
	ETag             string    `json:"eTag"`
	File             *struct{} `json:"file,omitempty"`
	Folder           *struct{} `json:"folder,omitempty"`
	DownloadURL      string    `json:"@microsoft.graph.downloadUrl,omitempty"`
}

// DriveItemCollection represents a collection of drive items
type DriveItemCollection struct {
	Value []DriveItem `json:"value"`
}

// NewService creates a new OneDrive service instance
func NewService(config cloudinterface.OneDriveConfig) cloudinterface.CloudService {
	return &Service{
		config: config,
	}
}

// GetServiceName returns the service name
func (s *Service) GetServiceName() string {
	return "OneDrive"
}

// Authenticate performs Azure AD authentication
func (s *Service) Authenticate(ctx context.Context) error {
	cred, err := azidentity.NewInteractiveBrowserCredential(&azidentity.InteractiveBrowserCredentialOptions{
		TenantID: s.config.TenantID,
		ClientID: s.config.ClientID,
	})
	if err != nil {
		return fmt.Errorf("error creating credential: %w", err)
	}

	token, err := cred.GetToken(ctx, policy.TokenRequestOptions{
		Scopes: []string{"https://graph.microsoft.com/.default"},
	})
	if err != nil {
		return fmt.Errorf("error getting token: %w", err)
	}

	s.token = token.Token
	return nil
}

// ListFiles lists files and folders in the specified path
func (s *Service) ListFiles(ctx context.Context, path string) ([]cloudinterface.CloudFile, error) {
	// Convert path to OneDrive API format
	apiPath := strings.TrimPrefix(path, "/")
	if apiPath == "" {
		apiPath = "root"
	} else {
		apiPath = fmt.Sprintf("root:/%s:", apiPath)
	}

	url := fmt.Sprintf("https://graph.microsoft.com/v1.0/me/drive/%s/children", apiPath)
	
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+s.token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var collection DriveItemCollection
	if err := json.NewDecoder(resp.Body).Decode(&collection); err != nil {
		return nil, err
	}

	var files []cloudinterface.CloudFile
	for _, item := range collection.Value {
		file := cloudinterface.CloudFile{
			ID:           item.ID,
			Name:         item.Name,
			LastModified: item.LastModifiedTime,
			ETag:         item.ETag,
			IsFolder:     item.Folder != nil,
			DownloadURL:  item.DownloadURL,
		}
		if item.Size != nil {
			file.Size = *item.Size
		}
		files = append(files, file)
	}

	return files, nil
}

// DownloadFile downloads a file by ID to the local path
func (s *Service) DownloadFile(ctx context.Context, fileID, localPath string) error {
	url := fmt.Sprintf("https://graph.microsoft.com/v1.0/me/drive/items/%s/content", fileID)
	
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+s.token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("download failed with status %d", resp.StatusCode)
	}

	file, err := os.Create(localPath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	return err
}
