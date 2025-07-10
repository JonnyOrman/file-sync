# Google Drive to Amazon S3 Backup Example

This example demonstrates how to use the cloud-file-backup package to back up files from Google Drive to Amazon S3.

## Features

- Authenticates with Google Drive
- Backs up files from the `/Photos` folder in Google Drive
- Stores files in an Amazon S3 bucket
- Tracks sync state to support incremental backups

## Prerequisites

- Go 1.20 or later
- Google Cloud Platform project with:
  - OAuth 2.0 Client ID and Secret
  - Google Drive API enabled
- AWS account with:
  - S3 bucket
  - IAM credentials with S3 access

## Setup

1. Set the required environment variables:

```bash
# Google Drive credentials
export GOOGLE_CLIENT_ID="your-client-id"
export GOOGLE_CLIENT_SECRET="your-client-secret"

# AWS S3 configuration
export S3_BUCKET="your-bucket-name"
export S3_REGION="your-aws-region"
export S3_PREFIX="gdrive-backup/"  # Optional, defaults to "gdrive-backup/"
```

2. Configure AWS credentials:
   - Either set the standard AWS environment variables (`AWS_ACCESS_KEY_ID`, `AWS_SECRET_ACCESS_KEY`)
   - Or configure the AWS CLI with `aws configure`

3. Build the example:

```bash
go build
```

## Running

```bash
./googledrive-s3
```

On the first run, you'll be prompted to authenticate with Google. The application will display a URL that you need to open in your browser, authorize the application, and then paste the authorization code back into the terminal.

## Configuration

You can modify the following in `main.go`:

- `backupConfig.Path`: The path in Google Drive to sync (default: `/Photos`)
- `storageConfig.S3Prefix`: The prefix for objects in S3 (default: `gdrive-backup/`)

## How It Works

1. The app authenticates with Google Drive using OAuth2
2. It lists files in the specified Google Drive path
3. For each file, it checks if it needs to be downloaded based on the sync state
4. New or modified files are downloaded to a temporary location and then uploaded to S3
5. The sync state is updated and saved for future incremental backups

## Sync State

The app maintains a sync state in `sync-state.json`. This file tracks:

- Last modified time of each file
- ETag for change detection
- File size

This allows the app to perform incremental backups, only downloading and uploading files that have changed since the last run. 