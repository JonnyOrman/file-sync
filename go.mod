module github.com/JonnyOrman/cloud-file-backup

go 1.23.0

toolchain go1.24.4

require (
	github.com/JonnyOrman/cloud-file-backup/cloud/interface v0.0.0
	github.com/JonnyOrman/cloud-file-backup/storage/interface v0.0.0
)

// Local module replacements
replace github.com/JonnyOrman/cloud-file-backup/cloud/interface => ./cloud/interface

replace github.com/JonnyOrman/cloud-file-backup/storage/interface => ./storage/interface

// Cloud service modules (conditionally included via build tags)
replace github.com/JonnyOrman/cloud-file-backup/cloud/onedrive => ./cloud/onedrive

replace github.com/JonnyOrman/cloud-file-backup/cloud/googledrive => ./cloud/googledrive

replace github.com/JonnyOrman/cloud-file-backup/cloud/dropbox => ./cloud/dropbox

// Storage service modules (conditionally included via build tags)
replace github.com/JonnyOrman/cloud-file-backup/storage/local => ./storage/local

replace github.com/JonnyOrman/cloud-file-backup/storage/s3 => ./storage/s3

replace github.com/JonnyOrman/cloud-file-backup/storage/gcs => ./storage/gcs

replace github.com/JonnyOrman/cloud-file-backup/storage/azblob => ./storage/azblob
