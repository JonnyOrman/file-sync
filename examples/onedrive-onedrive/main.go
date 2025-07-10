package main

import (
	"context"
	"log"
	"os"

	backup "github.com/JonnyOrman/cloud-file-backup"
	cloudinterface "github.com/JonnyOrman/cloud-file-backup/source/interface"
	onedrivesource "github.com/JonnyOrman/cloud-file-backup/source/onedrive"
	storageinterface "github.com/JonnyOrman/cloud-file-backup/destination/interface"
	onedrivedest "github.com/JonnyOrman/cloud-file-backup/destination/onedrive"
)

func main() {
	// Create source OneDrive service configuration
	// This could be a different OneDrive account or folder
	sourceConfig := cloudinterface.Config{
		Service:  "onedrive",
		ClientID: os.Getenv("SOURCE_ONEDRIVE_CLIENT_ID"),
		TenantID: os.Getenv("SOURCE_ONEDRIVE_TENANT_ID"),
	}

	// Create destination OneDrive service configuration
	// This could be a different OneDrive account or folder
	destConfig := storageinterface.Config{
		Type:      "onedrive",
		ClientID:  os.Getenv("DEST_ONEDRIVE_CLIENT_ID"),
		TenantID:  os.Getenv("DEST_ONEDRIVE_TENANT_ID"),
		ClientSecret: os.Getenv("DEST_ONEDRIVE_CLIENT_SECRET"),
	}

	// Create the source OneDrive service
	sourceService := onedrivesource.NewService(sourceConfig)

	// Create the destination OneDrive service
	destService := onedrivedest.NewService(destConfig)

	// Create backup client configuration
	backupConfig := backup.Config{
		Path: "/Documents/BackupSource", // Source OneDrive path to sync
	}

	// Create backup client with the services
	client, err := backup.NewBackupClient(sourceService, destService, backupConfig)
	if err != nil {
		log.Fatalf("Failed to create backup client: %v", err)
	}

	ctx := context.Background()

	// Initialize services
	log.Println("Initializing OneDrive services...")
	if err := client.Initialize(ctx); err != nil {
		log.Fatalf("Failed to initialize: %v", err)
	}

	// Load previous sync state
	log.Println("Loading sync state...")
	if err := client.LoadSyncState("onedrive-onedrive-sync-state.json"); err != nil {
		log.Printf("Warning: Could not load sync state: %v", err)
	}

	// Perform backup
	log.Println("Starting OneDrive to OneDrive backup...")
	if err := client.Sync(ctx); err != nil {
		log.Fatalf("Backup failed: %v", err)
	}

	// Save sync state
	log.Println("Saving sync state...")
	if err := client.SaveSyncState("onedrive-onedrive-sync-state.json"); err != nil {
		log.Printf("Warning: Could not save sync state: %v", err)
	}

	syncState := client.GetSyncState()
	log.Printf("OneDrive to OneDrive backup completed successfully! Synced %d files.", len(syncState.Files))
	log.Println("Files have been copied from source OneDrive to destination OneDrive")
} 