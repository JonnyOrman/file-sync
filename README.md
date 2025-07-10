# Cloud File Backup Package

A Go package for backing up files from cloud services (OneDrive, Google Drive, Dropbox) to various storage backends (local, S3, GCS, Azure Blob). This package provides a clean, implementation-agnostic API for integrating cloud file backup functionality into your Go applications.

## Installation

```bash
go get github.com/JonnyOrman/cloud-file-backup
```

## Quick Start

```go
package main

import (
    "context"
    "log"
    
    backup "github.com/JonnyOrman/cloud-file-backup"
    cloudinterface "github.com/JonnyOrman/cloud-file-backup/source/interface"
    onedrive "github.com/JonnyOrman/cloud-file-backup/source/onedrive"
    storageinterface "github.com/JonnyOrman/cloud-file-backup/destination/interface"
    localstorage "github.com/JonnyOrman/cloud-file-backup/destination/local"
)

func main() {
    // Create cloud service configuration
    cloudConfig := cloudinterface.Config{
        Service:  "onedrive",
        ClientID: "your-client-id",
        TenantID: "your-tenant-id",
    }

    // Create storage service configuration
    storageConfig := storageinterface.Config{
        Type:      "local",
        LocalPath: "./downloads",
    }

    // Create the cloud service
    cloudService := onedrive.NewService(cloudConfig)

    // Create the storage service
    storageService := localstorage.NewService(storageConfig)

    // Create backup client configuration
    backupConfig := backup.Config{
        Path: "/Documents", // Cloud path to sync
    }

    // Create backup client with the services
    client, err := backup.NewBackupClient(cloudService, storageService, backupConfig)
    if err != nil {
        log.Fatal(err)
    }
    
    ctx := context.Background()
    
    // Initialize services (authenticate, etc.)
    if err := client.Initialize(ctx); err != nil {
        log.Fatal(err)
    }
    
    // Load previous sync state (optional)
    if err := client.LoadSyncState(""); err != nil {
        log.Printf("Warning: %v", err)
    }
    
    // Perform backup
    if err := client.Sync(ctx); err != nil {
        log.Fatal(err)
    }
    
    // Save sync state
    if err := client.SaveSyncState(""); err != nil {
        log.Printf("Warning: %v", err)
    }
    
    log.Println("Backup completed successfully!")
}
```

For more examples, check out the [examples directory](./examples) which includes:
- [OneDrive to Local Storage](./examples/onedrive-local) - A detailed example with custom paths and sync state
- [Google Drive to S3](./examples/googledrive-s3) - An example showing Google Drive to Amazon S3 backup
```

## API Reference

### Types

#### `Config`
Configuration for the backup client.

```go
type Config struct {
    Path string `yaml:"path" mapstructure:"path"` // Cloud path to sync
}
```

#### `BackupClient`
Main client for performing backup operations.

#### `FileState`
Tracks file metadata for incremental sync.

```go
type FileState struct {
    LastModified time.Time `json:"last_modified"`
    ETag         string    `json:"etag"`
    Size         int64     `json:"size"`
}
```

#### `SyncState`
Tracks the overall synchronization state.

```go
type SyncState struct {
    Files map[string]FileState `json:"files"`
}
```

### Functions

#### `NewBackupClient(cloudService cloudinterface.CloudService, storageService storageinterface.StorageService, config Config) (*BackupClient, error)`
Creates a new backup client with the provided services and configuration.

#### `(*BackupClient) Initialize(ctx context.Context) error`
Authenticates with cloud and storage services. Must be called before `Sync()`.

#### `(*BackupClient) LoadSyncState(stateFilePath string) error`
Loads synchronization state from a file. Pass empty string to use default path.

#### `(*BackupClient) SaveSyncState(stateFilePath string) error`
Saves synchronization state to a file. Pass empty string to use default path.

#### `(*BackupClient) Sync(ctx context.Context) error`
Performs the backup operation, syncing files from cloud to storage.

#### `(*BackupClient) GetSyncState() SyncState`
Returns a copy of the current synchronization state.

## Available Implementations

### Cloud Services

Import the specific cloud service implementation you want to use:

```go
import "github.com/JonnyOrman/cloud-file-backup/source/onedrive"
// or
import "github.com/JonnyOrman/cloud-file-backup/source/googledrive"
// or
import "github.com/JonnyOrman/cloud-file-backup/source/dropbox"
```

#### OneDrive
```go
cloudConfig := cloudinterface.Config{
    Service:  "onedrive",
    ClientID: "your-client-id",
    TenantID: "your-tenant-id",
}
cloudService := onedrive.NewService(cloudConfig)
```

#### Google Drive
```go
cloudConfig := cloudinterface.Config{
    Service:      "googledrive",
    ClientID:     "your-client-id",
    ClientSecret: "your-client-secret",
}
cloudService := googledrive.NewService(cloudConfig)
```

#### Dropbox
```go
cloudConfig := cloudinterface.Config{
    Service:     "dropbox",
    AccessToken: "your-access-token",
}
cloudService := dropbox.NewService(cloudConfig)
```

### Storage Backends

Import the specific storage implementation you want to use:

```go
import "github.com/JonnyOrman/cloud-file-backup/destination/local"
// or
import "github.com/JonnyOrman/cloud-file-backup/destination/s3"
// or
import "github.com/JonnyOrman/cloud-file-backup/destination/gcs"
// or
import "github.com/JonnyOrman/cloud-file-backup/destination/azblob"
```

#### Local Storage
```go
storageConfig := storageinterface.Config{
    Type:      "local",
    LocalPath: "./downloads",
}
storageService := localstorage.NewService(storageConfig)
```

#### Amazon S3
```go
storageConfig := storageinterface.Config{
    Type:     "s3",
    S3Bucket: "my-backup-bucket",
    S3Region: "us-east-1",
    S3Prefix: "backups/",
}
storageService := s3storage.NewService(storageConfig)
```

#### Google Cloud Storage
```go
storageConfig := storageinterface.Config{
    Type:       "gcs",
    GCSBucket:  "my-backup-bucket",
    GCSProject: "my-project-id",
    GCSPrefix:  "backups/",
    GCSKeyFile: "/path/to/service-account.json",
}
storageService := gcsstorage.NewService(storageConfig)
```

#### Azure Blob Storage
```go
storageConfig := storageinterface.Config{
    Type:                "azblob",
    AzureBlobAccount:    "mystorageaccount",
    AzureBlobContainer:  "backups",
    AzureBlobKey:        "account-key",
    AzureBlobPrefix:     "backups/",
}
storageService := azblobstorage.NewService(storageConfig)
```

## Features

- **Implementation-Agnostic**: Core package has no direct dependencies on specific implementations
- **Incremental Sync**: Only downloads new or modified files
- **Multiple Cloud Services**: OneDrive, Google Drive, Dropbox
- **Multiple Storage Backends**: Local, S3, GCS, Azure Blob
- **Modular Architecture**: Each service is a separate module
- **Build Tags**: Compile only the services you need
- **State Tracking**: Persistent sync state to avoid re-downloading

## Build Tags

You can compile with only the services you need using build tags:

```bash
# Only OneDrive and local storage
go build -tags "onedrive local" your-app.go

# Only Google Drive and S3
go build -tags "googledrive s3" your-app.go

# All services
go build -tags "onedrive googledrive dropbox local s3 gcs azblob" your-app.go
```

## Error Handling

All methods return descriptive errors. The package does not log internally, leaving logging decisions to the consuming application.

```go
if err := client.Sync(ctx); err != nil {
    log.Printf("Backup failed: %v", err)
    // Handle error appropriately
}
```

## Threading

The package is not thread-safe. Use separate client instances for concurrent operations or implement your own synchronization.

## License

This project is licensed under the MIT License. 