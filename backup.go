package cloudfilebackup

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/JonnyOrman/cloud-file-backup/source/interface"
	"github.com/JonnyOrman/cloud-file-backup/destination/interface"
)

// FileState tracks the state of a file for incremental sync
type FileState struct {
	LastModified time.Time `json:"last_modified"`
	ETag         string    `json:"etag"`
	Size         int64     `json:"size"`
}

// SyncState tracks the overall sync state
type SyncState struct {
	Files map[string]FileState `json:"files"`
}

// Config represents the application configuration
type Config struct {
	Path string `yaml:"path" mapstructure:"path"` // Cloud path to sync
}

// BackupClient provides the main API for cloud file backup operations
type BackupClient struct {
	config         Config
	cloudService   cloudinterface.CloudService
	storageService storageinterface.StorageService
	syncState      *SyncState
}

// NewBackupClient creates a new backup client with the provided services and configuration
func NewBackupClient(cloudService cloudinterface.CloudService, storageService storageinterface.StorageService, config Config) (*BackupClient, error) {
	if cloudService == nil {
		return nil, fmt.Errorf("cloud service cannot be nil")
	}
	
	if storageService == nil {
		return nil, fmt.Errorf("storage service cannot be nil")
	}
	
	if config.Path == "" {
		return nil, fmt.Errorf("path must be specified")
	}

	return &BackupClient{
		config:         config,
		cloudService:   cloudService,
		storageService: storageService,
		syncState:      &SyncState{Files: make(map[string]FileState)},
	}, nil
}

// Initialize authenticates with cloud and storage services
func (c *BackupClient) Initialize(ctx context.Context) error {
	// Initialize cloud service
	if err := c.cloudService.Authenticate(ctx); err != nil {
		return fmt.Errorf("failed to authenticate with cloud service: %w", err)
	}

	// Initialize storage service
	if err := c.storageService.Initialize(ctx); err != nil {
		return fmt.Errorf("failed to initialize storage service: %w", err)
	}

	return nil
}

// LoadSyncState loads the synchronization state from a file
func (c *BackupClient) LoadSyncState(stateFilePath string) error {
	if stateFilePath == "" {
		stateFilePath = ".cloud-file-backup-state.json"
	}

	syncState, err := loadSyncState(stateFilePath)
	if err != nil {
		// If state file doesn't exist, start with empty state
		c.syncState = &SyncState{Files: make(map[string]FileState)}
		return nil
	}

	c.syncState = syncState
	return nil
}

// SaveSyncState saves the synchronization state to a file
func (c *BackupClient) SaveSyncState(stateFilePath string) error {
	if stateFilePath == "" {
		stateFilePath = ".cloud-file-backup-state.json"
	}
	return saveSyncState(c.syncState, stateFilePath)
}

// Sync performs the backup operation, syncing files from cloud to storage
func (c *BackupClient) Sync(ctx context.Context) error {
	return c.syncFiles(ctx, c.config.Path, "")
}

// GetSyncState returns a copy of the current sync state
func (c *BackupClient) GetSyncState() SyncState {
	stateCopy := SyncState{Files: make(map[string]FileState)}
	for k, v := range c.syncState.Files {
		stateCopy.Files[k] = v
	}
	return stateCopy
}

// syncFiles performs the actual file synchronization
func (c *BackupClient) syncFiles(ctx context.Context, cloudPath, relativePath string) error {
	return c.syncFilesRecursive(ctx, cloudPath, relativePath)
}

func (c *BackupClient) syncFilesRecursive(ctx context.Context, cloudPath, relativePath string) error {
	files, err := c.cloudService.ListFiles(ctx, cloudPath)
	if err != nil {
		return fmt.Errorf("failed to list files in %s: %w", cloudPath, err)
	}

	for _, file := range files {
		fileRelativePath := filepath.Join(relativePath, file.Name)

		if file.IsFolder {
			// Create directory in storage
			if err := c.storageService.CreateDirectory(ctx, fileRelativePath); err != nil {
				return fmt.Errorf("could not create directory %s: %w", fileRelativePath, err)
			}

			// Recursively sync folder contents
			folderCloudPath := cloudPath
			if cloudPath == "/" {
				folderCloudPath = "/" + file.Name
			} else {
				folderCloudPath = cloudPath + "/" + file.Name
			}

			if err := c.syncFilesRecursive(ctx, folderCloudPath, fileRelativePath); err != nil {
				return err
			}
		} else {
			// Check if file needs to be downloaded
			lastState, exists := c.syncState.Files[fileRelativePath]
			needsDownload := !exists ||
				lastState.LastModified.Before(file.LastModified) ||
				lastState.ETag != file.ETag ||
				lastState.Size != file.Size

			if needsDownload {
				if err := c.downloadFile(ctx, file, fileRelativePath); err != nil {
					return fmt.Errorf("failed to download %s: %w", file.Name, err)
				}

				// Update sync state
				c.syncState.Files[fileRelativePath] = FileState{
					LastModified: file.LastModified,
					ETag:         file.ETag,
					Size:         file.Size,
				}
			}
		}
	}

	return nil
}

func (c *BackupClient) downloadFile(ctx context.Context, file cloudinterface.CloudFile, relativePath string) error {
	// Create a temporary file for downloading
	tempFile, err := os.CreateTemp("", "cloud-backup-*")
	if err != nil {
		return fmt.Errorf("failed to create temp file: %w", err)
	}
	defer os.Remove(tempFile.Name())
	defer tempFile.Close()

	// Download to temp file
	if err := c.cloudService.DownloadFile(ctx, file.ID, tempFile.Name()); err != nil {
		return fmt.Errorf("failed to download file: %w", err)
	}

	// Reopen temp file for reading
	tempFile.Close()
	downloadedFile, err := os.Open(tempFile.Name())
	if err != nil {
		return fmt.Errorf("failed to reopen temp file: %w", err)
	}
	defer downloadedFile.Close()

	// Upload to storage service
	return c.storageService.WriteFile(ctx, relativePath, downloadedFile, file.Size)
}

func loadSyncState(stateFilePath string) (*SyncState, error) {
	data, err := os.ReadFile(stateFilePath)
	if err != nil {
		return nil, err
	}

	var state SyncState
	if err := json.Unmarshal(data, &state); err != nil {
		return nil, err
	}

	return &state, nil
}

func saveSyncState(state *SyncState, stateFilePath string) error {
	data, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(stateFilePath, data, 0644)
} 