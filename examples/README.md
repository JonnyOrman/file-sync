# Cloud File Backup Examples

This directory contains standalone example applications that demonstrate how to use the cloud-file-backup package with different cloud services and storage backends.

## Available Examples

### [OneDrive to Local Storage](./onedrive-local)

A detailed example showing how to back up files from OneDrive to a local directory with custom paths and sync state.

```bash
cd onedrive-local
go build
./onedrive-local
```

### [Google Drive to S3](./googledrive-s3)

An example showing how to back up files from Google Drive to Amazon S3.

```bash
cd googledrive-s3
go build
./googledrive-s3
```

### [OneDrive to OneDrive](./onedrive-onedrive)

An example showing how to sync files between OneDrive services (cross-account backup, folder sync, etc.).

```bash
cd onedrive-onedrive
go build
./onedrive-onedrive
```

## Creating Your Own Example

To create your own example application:

1. Create a new directory for your example
2. Create a `go.mod` file with the necessary dependencies and replacements
3. Create a `main.go` file that:
   - Imports the necessary packages
   - Creates cloud and storage services
   - Configures and initializes the backup client
   - Performs the backup operation

## Available Cloud Services

- OneDrive (`github.com/JonnyOrman/cloud-file-backup/source/onedrive`)
- Google Drive (`github.com/JonnyOrman/cloud-file-backup/source/googledrive`)
- Dropbox (`github.com/JonnyOrman/cloud-file-backup/source/dropbox`)

## Available Storage Backends

- Local Storage (`github.com/JonnyOrman/cloud-file-backup/destination/local`)
- Amazon S3 (`github.com/JonnyOrman/cloud-file-backup/destination/s3`)
- Google Cloud Storage (`github.com/JonnyOrman/cloud-file-backup/destination/gcs`)
- Azure Blob Storage (`github.com/JonnyOrman/cloud-file-backup/destination/azblob`) 