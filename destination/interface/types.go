package storageinterface

import (
	"context"
	"io"
)

// StorageService defines the interface that all storage services must implement
type StorageService interface {
	// GetStorageName returns the name of the storage service
	GetStorageName() string
	
	// Initialize initializes the storage service
	Initialize(ctx context.Context) error
	
	// WriteFile writes a file to the storage service
	WriteFile(ctx context.Context, relativePath string, content io.Reader, size int64) error
	
	// FileExists checks if a file exists in the storage service
	FileExists(ctx context.Context, relativePath string) (bool, error)
	
	// CreateDirectory creates a directory in the storage service
	CreateDirectory(ctx context.Context, relativePath string) error
}

// LocalConfig represents the configuration for local filesystem destination service
type LocalConfig struct {
	LocalPath string `yaml:"local_path" mapstructure:"local_path"`
}

// S3Config represents the configuration for AWS S3 destination service
type S3Config struct {
	S3Bucket string `yaml:"s3_bucket" mapstructure:"s3_bucket"`
	S3Region string `yaml:"s3_region" mapstructure:"s3_region"`
	S3Prefix string `yaml:"s3_prefix" mapstructure:"s3_prefix"`
}

// GCSConfig represents the configuration for Google Cloud Storage destination service
type GCSConfig struct {
	GCSBucket  string `yaml:"gcs_bucket" mapstructure:"gcs_bucket"`
	GCSProject string `yaml:"gcs_project" mapstructure:"gcs_project"`
	GCSPrefix  string `yaml:"gcs_prefix" mapstructure:"gcs_prefix"`
	GCSKeyFile string `yaml:"gcs_key_file" mapstructure:"gcs_key_file"`
}

// AzureBlobConfig represents the configuration for Azure Blob Storage destination service
type AzureBlobConfig struct {
	AzureBlobAccount   string `yaml:"azure_blob_account" mapstructure:"azure_blob_account"`
	AzureBlobContainer string `yaml:"azure_blob_container" mapstructure:"azure_blob_container"`
	AzureBlobPrefix    string `yaml:"azure_blob_prefix" mapstructure:"azure_blob_prefix"`
	AzureBlobKey       string `yaml:"azure_blob_key" mapstructure:"azure_blob_key"`
}

// OneDriveConfig represents the configuration for OneDrive destination service
type OneDriveConfig struct {
	ClientID     string `yaml:"client_id" mapstructure:"client_id"`
	ClientSecret string `yaml:"client_secret" mapstructure:"client_secret"`
	TenantID     string `yaml:"tenant_id" mapstructure:"tenant_id"`
}

// GoogleDriveConfig represents the configuration for Google Drive destination service
type GoogleDriveConfig struct {
	ClientID     string `yaml:"client_id" mapstructure:"client_id"`
	ClientSecret string `yaml:"client_secret" mapstructure:"client_secret"`
}

// DropboxConfig represents the configuration for Dropbox destination service
type DropboxConfig struct {
	AccessToken string `yaml:"access_token" mapstructure:"access_token"`
}
