package main

import (
	"context"
	"log"
	"os"

	backup "github.com/JonnyOrman/cloud-file-backup"
	cloudinterface "github.com/JonnyOrman/cloud-file-backup/source/interface"
	googledrive "github.com/JonnyOrman/cloud-file-backup/source/googledrive"
	storageinterface "github.com/JonnyOrman/cloud-file-backup/destination/interface"
	s3storage "github.com/JonnyOrman/cloud-file-backup/destination/s3"
)

func main() {
	// Create Google Drive source service configuration
	sourceConfig := cloudinterface.GoogleDriveConfig{
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
	}

	// Create S3 destination service configuration
	destConfig := storageinterface.S3Config{
		S3Bucket: os.Getenv("S3_BUCKET"),
		S3Region: os.Getenv("S3_REGION"),
		S3Prefix: os.Getenv("S3_PREFIX"),
	}

	// Set default S3 prefix if not specified
	if destConfig.S3Prefix == "" {
		destConfig.S3Prefix = "gdrive-backup/"
	}

	// Create the source service
	sourceService := googledrive.NewService(sourceConfig)

	// Create the destination service
	destService := s3storage.NewService(destConfig)

	// Create backup client configuration
	backupConfig := backup.Config{
		Path: "/Photos", // Cloud path to sync
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
	log.Printf("Files are stored in S3 bucket: %s with prefix: %s", 
		destConfig.S3Bucket, destConfig.S3Prefix)
} 