# OneDrive to OneDrive Example

This example demonstrates how to sync files between OneDrive services. This can be useful for:

- **Cross-account backup**: Backing up files from one OneDrive account to another
- **Folder synchronization**: Syncing between different folders in the same OneDrive account
- **Data migration**: Moving files from one OneDrive location to another
- **Redundant storage**: Creating backup copies of important OneDrive folders

## Use Cases

### 1. Personal to Business OneDrive
Sync important personal files to your business OneDrive account for work access.

### 2. Folder Backup
Create backup copies of critical OneDrive folders to a separate location.

### 3. Account Migration
Move files from an old OneDrive account to a new one.

### 4. Cross-Region Sync
Sync files between OneDrive accounts in different regions for redundancy.

## Setup

### 1. Environment Variables

Set the following environment variables:

```bash
# Source OneDrive (where files come from)
export SOURCE_ONEDRIVE_CLIENT_ID="your-source-client-id"
export SOURCE_ONEDRIVE_TENANT_ID="your-source-tenant-id"

# Destination OneDrive (where files go to)
export DEST_ONEDRIVE_CLIENT_ID="your-dest-client-id"
export DEST_ONEDRIVE_TENANT_ID="your-dest-tenant-id"
export DEST_ONEDRIVE_CLIENT_SECRET="your-dest-client-secret"
```

### 2. Azure AD App Registration

You'll need to register applications in Azure AD for both source and destination OneDrive accounts:

#### Source OneDrive App:
- **Redirect URI**: `http://localhost:8080/auth`
- **API Permissions**: `Files.Read.All` (for reading files)

#### Destination OneDrive App:
- **Redirect URI**: `http://localhost:8080/auth`
- **API Permissions**: `Files.ReadWrite.All` (for reading and writing files)
- **Client Secret**: Generate a client secret for authentication

### 3. Authentication

The example uses different authentication methods:

- **Source**: Uses interactive authentication (browser-based)
- **Destination**: Uses client secret authentication

## Running the Example

```bash
cd onedrive-onedrive
go build
./onedrive-onedrive
```

## Configuration

### Source Path
Modify the source path in `main.go`:

```go
backupConfig := backup.Config{
    Path: "/Documents/BackupSource", // Change this to your source folder
}
```

### Sync State
The example saves sync state to `onedrive-onedrive-sync-state.json` to track what files have been synced.

## Important Notes

⚠️ **Note**: The OneDrive destination service is currently a placeholder implementation. To make this fully functional, you would need to:

1. Implement the Microsoft Graph API integration for file uploads
2. Handle OAuth2 token management for the destination service
3. Implement proper error handling for API rate limits
4. Add support for large file uploads with session uploads

## Expected Output

```
Initializing OneDrive services...
Loading sync state...
Starting OneDrive to OneDrive backup...
OneDrive to OneDrive backup completed successfully! Synced 15 files.
Files have been copied from source OneDrive to destination OneDrive
```

## Troubleshooting

### Common Issues

1. **Authentication Errors**: Ensure your Azure AD app registrations are correct
2. **Permission Errors**: Verify API permissions are granted for both apps
3. **Rate Limiting**: OneDrive has API rate limits; the sync handles this automatically
4. **Large Files**: Files larger than 4MB may require special handling

### Debug Mode

Add debug logging by setting the log level:

```go
log.SetFlags(log.LstdFlags | log.Lshortfile)
```

## Security Considerations

- Store client secrets securely (use environment variables or secret management)
- Use least-privilege permissions for Azure AD apps
- Regularly rotate client secrets
- Monitor API usage to detect unusual activity 