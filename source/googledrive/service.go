package googledrive

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
	"github.com/JonnyOrman/cloud-file-backup/cloud/interface"
)

// Service implements the Google Drive cloud service
type Service struct {
	config       cloudinterface.Config
	driveService *drive.Service
}

// NewService creates a new Google Drive service instance
func NewService(config cloudinterface.Config) cloudinterface.CloudService {
	return &Service{
		config: config,
	}
}

// GetServiceName returns the service name
func (s *Service) GetServiceName() string {
	return "Google Drive"
}

// Authenticate performs OAuth2 authentication for Google Drive
func (s *Service) Authenticate(ctx context.Context) error {
	// Create OAuth2 config
	oauthConfig := &oauth2.Config{
		ClientID:     s.config.ClientID,
		ClientSecret: s.config.ClientSecret,
		RedirectURL:  "http://localhost:8080/callback",
		Scopes:       []string{drive.DriveReadonlyScope},
		Endpoint:     google.Endpoint,
	}

	// Get token
	token, err := s.getToken(ctx, oauthConfig)
	if err != nil {
		return fmt.Errorf("error getting Google token: %w", err)
	}

	// Create Drive service
	driveService, err := drive.NewService(ctx, option.WithTokenSource(oauthConfig.TokenSource(ctx, token)))
	if err != nil {
		return fmt.Errorf("error creating Drive service: %w", err)
	}

	s.driveService = driveService
	return nil
}

// ListFiles lists files and folders in the specified path
func (s *Service) ListFiles(ctx context.Context, path string) ([]cloudinterface.CloudFile, error) {
	var parentID string
	
	if path == "" || path == "/" {
		parentID = "root"
	} else {
		// Find folder by path
		var err error
		parentID, err = s.findFolderByPath(ctx, path)
		if err != nil {
			return nil, err
		}
	}

	query := fmt.Sprintf("'%s' in parents and trashed=false", parentID)
	fileList, err := s.driveService.Files.List().Q(query).Fields("files(id,name,size,modifiedTime,mimeType,md5Checksum)").Do()
	if err != nil {
		return nil, err
	}

	var files []cloudinterface.CloudFile
	for _, file := range fileList.Files {
		cloudFile := cloudinterface.CloudFile{
			ID:       file.Id,
			Name:     file.Name,
			ETag:     file.Md5Checksum,
			IsFolder: file.MimeType == "application/vnd.google-apps.folder",
		}
		
		if file.Size != 0 {
			cloudFile.Size = file.Size
		}
		
		if file.ModifiedTime != "" {
			if modTime, err := time.Parse(time.RFC3339, file.ModifiedTime); err == nil {
				cloudFile.LastModified = modTime
			}
		}
		
		files = append(files, cloudFile)
	}

	return files, nil
}

// DownloadFile downloads a file by ID to the local path
func (s *Service) DownloadFile(ctx context.Context, fileID, localPath string) error {
	resp, err := s.driveService.Files.Get(fileID).Download()
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	file, err := os.Create(localPath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	return err
}

// findFolderByPath finds a folder by its path
func (s *Service) findFolderByPath(ctx context.Context, path string) (string, error) {
	parts := strings.Split(strings.Trim(path, "/"), "/")
	currentID := "root"
	
	for _, part := range parts {
		if part == "" {
			continue
		}
		
		query := fmt.Sprintf("'%s' in parents and name='%s' and mimeType='application/vnd.google-apps.folder' and trashed=false", currentID, part)
		fileList, err := s.driveService.Files.List().Q(query).Fields("files(id)").Do()
		if err != nil {
			return "", err
		}
		
		if len(fileList.Files) == 0 {
			return "", fmt.Errorf("folder not found: %s", part)
		}
		
		currentID = fileList.Files[0].Id
	}
	
	return currentID, nil
}

// getToken handles the OAuth2 token flow
func (s *Service) getToken(ctx context.Context, oauthConfig *oauth2.Config) (*oauth2.Token, error) {
	// Try to load token from file first
	tokenFile := "google-token.json"
	if token, err := s.loadTokenFromFile(tokenFile); err == nil {
		// Check if token is still valid
		if token.Valid() {
			return token, nil
		}
		// Try to refresh if we have a refresh token
		if token.RefreshToken != "" {
			tokenSource := oauthConfig.TokenSource(ctx, token)
			newToken, err := tokenSource.Token()
			if err == nil {
				s.saveTokenToFile(tokenFile, newToken)
				return newToken, nil
			}
		}
	}

	// Start OAuth2 flow
	authURL := oauthConfig.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser:\n%v\n", authURL)
	fmt.Print("Enter the authorization code: ")
	
	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		return nil, fmt.Errorf("unable to read authorization code: %w", err)
	}

	token, err := oauthConfig.Exchange(ctx, authCode)
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve token from web: %w", err)
	}

	// Save token for future use
	s.saveTokenToFile(tokenFile, token)
	return token, nil
}

// loadTokenFromFile loads a token from a file
func (s *Service) loadTokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	
	token := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(token)
	return token, err
}

// saveTokenToFile saves a token to a file
func (s *Service) saveTokenToFile(file string, token *oauth2.Token) error {
	f, err := os.Create(file)
	if err != nil {
		return err
	}
	defer f.Close()
	
	return json.NewEncoder(f).Encode(token)
} 