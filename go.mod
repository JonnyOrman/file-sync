module github.com/JonnyOrman/cloud-file-backup

go 1.23.0

toolchain go1.24.4

require (
	github.com/JonnyOrman/cloud-file-backup/source/interface v0.0.0
	github.com/JonnyOrman/cloud-file-backup/destination/interface v0.0.0
)

// Local module replacements
replace github.com/JonnyOrman/cloud-file-backup/source/interface => ./source/interface

replace github.com/JonnyOrman/cloud-file-backup/destination/interface => ./destination/interface

// Source service modules (conditionally included via build tags)
replace github.com/JonnyOrman/cloud-file-backup/source/onedrive => ./source/onedrive

replace github.com/JonnyOrman/cloud-file-backup/source/googledrive => ./source/googledrive

replace github.com/JonnyOrman/cloud-file-backup/source/dropbox => ./source/dropbox

replace github.com/JonnyOrman/cloud-file-backup/source/local => ./source/local

replace github.com/JonnyOrman/cloud-file-backup/source/s3 => ./source/s3

replace github.com/JonnyOrman/cloud-file-backup/source/gcs => ./source/gcs

replace github.com/JonnyOrman/cloud-file-backup/source/azblob => ./source/azblob

// Destination service modules (conditionally included via build tags)
replace github.com/JonnyOrman/cloud-file-backup/destination/local => ./destination/local

replace github.com/JonnyOrman/cloud-file-backup/destination/s3 => ./destination/s3

replace github.com/JonnyOrman/cloud-file-backup/destination/gcs => ./destination/gcs

replace github.com/JonnyOrman/cloud-file-backup/destination/azblob => ./destination/azblob

replace github.com/JonnyOrman/cloud-file-backup/destination/onedrive => ./destination/onedrive

replace github.com/JonnyOrman/cloud-file-backup/destination/googledrive => ./destination/googledrive

replace github.com/JonnyOrman/cloud-file-backup/destination/dropbox => ./destination/dropbox
