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
	// Create cloud service configuration for Google Drive
	cloudConfig := cloudinterface.Config{
		Service:      "googledrive",
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
	}

	// Create storage service configuration for S3
	storageConfig := storageinterface.Config{
		Type:     "s3",
		S3Bucket: os.Getenv("S3_BUCKET"),
		S3Region: os.Getenv("S3_REGION"),
		S3Prefix: os.Getenv("S3_PREFIX"),
	}

	// Set default S3 prefix if not specified
	if storageConfig.S3Prefix == "" {
		storageConfig.S3Prefix = "gdrive-backup/"
	}

	// Create the cloud service
	cloudService := googledrive.NewService(cloudConfig)

	// Create the storage service
	storageService := s3storage.NewService(storageConfig)

	// Create backup client configuration
	backupConfig := backup.Config{
		Path: "/Photos", // Cloud path to sync
	}

	// Create backup client with the services
	client, err := backup.NewBackupClient(cloudService, storageService, backupConfig)
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
		storageConfig.S3Bucket, storageConfig.S3Prefix)
} 