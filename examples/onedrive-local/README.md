# OneDrive to Local Storage Backup Example

This example demonstrates how to use the cloud-file-backup package to back up files from OneDrive to local storage.

## Features

- Authenticates with Microsoft OneDrive
- Backs up files from the `/Documents` folder in OneDrive
- Stores files in a local `downloads` directory
- Tracks sync state to support incremental backups

## Prerequisites

- Go 1.20 or later
- Microsoft Azure App Registration with the following:
  - Client ID
  - Tenant ID
  - Appropriate permissions for OneDrive access

## Setup

1. Set the required environment variables:

```bash
export ONEDRIVE_CLIENT_ID="your-client-id"
export ONEDRIVE_TENANT_ID="your-tenant-id"
```

2. Build the example:

```bash
go build
```

## Running

```bash
./onedrive-local
```

On the first run, you'll be prompted to authenticate with Microsoft. A browser window will open for you to sign in.

## Configuration

You can modify the following in `main.go`:

- `backupConfig.Path`: The path in OneDrive to sync (default: `/Documents`)
- `downloadsDir`: The local directory to store files (default: `./downloads`)

## How It Works

1. The app authenticates with OneDrive using the Microsoft Authentication Library
2. It lists files in the specified OneDrive path
3. For each file, it checks if it needs to be downloaded based on the sync state
4. New or modified files are downloaded to the local directory
5. The sync state is updated and saved for future incremental backups

## Sync State

The app maintains a sync state in `sync-state.json`. This file tracks:

- Last modified time of each file
- ETag for change detection
- File size

This allows the app to perform incremental backups, only downloading files that have changed since the last run. 