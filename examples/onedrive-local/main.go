package main

import (
	"context"
	"log"
	"os"
	"path/filepath"

	backup "github.com/JonnyOrman/cloud-file-backup"
	cloudinterface "github.com/JonnyOrman/cloud-file-backup/source/interface"
	onedrive "github.com/JonnyOrman/cloud-file-backup/source/onedrive"
	storageinterface "github.com/JonnyOrman/cloud-file-backup/destination/interface"
	localstorage "github.com/JonnyOrman/cloud-file-backup/destination/local"
)

func main() {
	// Create OneDrive source service configuration
	sourceConfig := cloudinterface.OneDriveConfig{
		ClientID: os.Getenv("ONEDRIVE_CLIENT_ID"),
		TenantID: os.Getenv("ONEDRIVE_TENANT_ID"),
	}

	// Create local storage service configuration
	// Create downloads directory if it doesn't exist
	downloadsDir := "./downloads"
	if err := os.MkdirAll(downloadsDir, 0755); err != nil {
		log.Fatalf("Failed to create downloads directory: %v", err)
	}
	
	absPath, err := filepath.Abs(downloadsDir)
	if err != nil {
		log.Fatalf("Failed to get absolute path: %v", err)
	}
	
	destConfig := storageinterface.LocalConfig{
		LocalPath: absPath,
	}

	// Create the source service
	sourceService := onedrive.NewService(sourceConfig)

	// Create the destination service
	destService := localstorage.NewService(destConfig)

	// Create backup client configuration
	backupConfig := backup.Config{
		Path: "/Documents", // Cloud path to sync
	}

	// Create backup client with the services
	client, err := backup.NewBackupClient(sourceService, destService, backupConfig)
	if err != nil {
		log.Fatalf("Failed to create backup client: %v", err)
	}

	ctx := context.Background()

	// Initialize services
	log.Println("Initializing services...")
	if err := client.Initialize(ctx); err != nil {
		log.Fatalf("Failed to initialize: %v", err)
	}

	// Load previous sync state
	log.Println("Loading sync state...")
	if err := client.LoadSyncState("sync-state.json"); err != nil {
		log.Printf("Warning: Could not load sync state: %v", err)
	}

	// Perform backup
	log.Println("Starting backup...")
	if err := client.Sync(ctx); err != nil {
		log.Fatalf("Backup failed: %v", err)
	}

	// Save sync state
	log.Println("Saving sync state...")
	if err := client.SaveSyncState("sync-state.json"); err != nil {
		log.Printf("Warning: Could not save sync state: %v", err)
	}

	syncState := client.GetSyncState()
	log.Printf("Backup completed successfully! Synced %d files.", len(syncState.Files))
	log.Printf("Files are stored in: %s", absPath)
} 