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

- OneDrive (`github.com/JonnyOrman/cloud-file-backup/cloud/onedrive`)
- Google Drive (`github.com/JonnyOrman/cloud-file-backup/cloud/googledrive`)
- Dropbox (`github.com/JonnyOrman/cloud-file-backup/cloud/dropbox`)

## Available Storage Backends

- Local Storage (`github.com/JonnyOrman/cloud-file-backup/storage/local`)
- Amazon S3 (`github.com/JonnyOrman/cloud-file-backup/storage/s3`)
- Google Cloud Storage (`github.com/JonnyOrman/cloud-file-backup/storage/gcs`)
- Azure Blob Storage (`github.com/JonnyOrman/cloud-file-backup/storage/azblob`) 