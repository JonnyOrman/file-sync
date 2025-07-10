package cloudinterface

import (
	"context"
	"time"
)

// CloudFile represents a file in cloud storage
type CloudFile struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Size         int64     `json:"size"`
	LastModified time.Time `json:"last_modified"`
	ETag         string    `json:"etag"`
	IsFolder     bool      `json:"is_folder"`
	DownloadURL  string    `json:"download_url,omitempty"`
}

// CloudService defines the interface that all cloud services must implement
type CloudService interface {
	// GetServiceName returns the name of the cloud service
	GetServiceName() string
	
	// Authenticate performs authentication with the cloud service
	Authenticate(ctx context.Context) error
	
	// ListFiles lists all files in the specified path
	ListFiles(ctx context.Context, path string) ([]CloudFile, error)
	
	// DownloadFile downloads a file to the specified local path
	DownloadFile(ctx context.Context, fileID, localPath string) error
}

// OneDriveConfig represents the configuration for OneDrive source service
type OneDriveConfig struct {
	ClientID string `yaml:"client_id" mapstructure:"client_id"`
	TenantID string `yaml:"tenant_id" mapstructure:"tenant_id"`
}

// GoogleDriveConfig represents the configuration for Google Drive source service
type GoogleDriveConfig struct {
	ClientID     string `yaml:"client_id" mapstructure:"client_id"`
	ClientSecret string `yaml:"client_secret" mapstructure:"client_secret"`
}

// DropboxConfig represents the configuration for Dropbox source service
type DropboxConfig struct {
	AccessToken string `yaml:"access_token" mapstructure:"access_token"`
}

// LocalConfig represents the configuration for local filesystem source service
type LocalConfig struct {
	LocalPath string `yaml:"local_path" mapstructure:"local_path"`
}

// S3Config represents the configuration for AWS S3 source service
type S3Config struct {
	S3Bucket string `yaml:"s3_bucket" mapstructure:"s3_bucket"`
	S3Region string `yaml:"s3_region" mapstructure:"s3_region"`
	S3Prefix string `yaml:"s3_prefix" mapstructure:"s3_prefix"`
}

// GCSConfig represents the configuration for Google Cloud Storage source service
type GCSConfig struct {
	GCSBucket  string `yaml:"gcs_bucket" mapstructure:"gcs_bucket"`
	GCSProject string `yaml:"gcs_project" mapstructure:"gcs_project"`
	GCSPrefix  string `yaml:"gcs_prefix" mapstructure:"gcs_prefix"`
	GCSKeyFile string `yaml:"gcs_key_file" mapstructure:"gcs_key_file"`
}

// AzureBlobConfig represents the configuration for Azure Blob Storage source service
type AzureBlobConfig struct {
	AzureBlobAccount   string `yaml:"azure_blob_account" mapstructure:"azure_blob_account"`
	AzureBlobContainer string `yaml:"azure_blob_container" mapstructure:"azure_blob_container"`
	AzureBlobPrefix    string `yaml:"azure_blob_prefix" mapstructure:"azure_blob_prefix"`
	AzureBlobKey       string `yaml:"azure_blob_key" mapstructure:"azure_blob_key"`
}
