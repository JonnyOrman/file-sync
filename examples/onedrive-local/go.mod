module github.com/JonnyOrman/cloud-file-backup/examples/onedrive-local

go 1.20

require (
	github.com/JonnyOrman/cloud-file-backup v0.0.0
	github.com/JonnyOrman/cloud-file-backup/cloud/interface v0.0.0
	github.com/JonnyOrman/cloud-file-backup/cloud/onedrive v0.0.0
	github.com/JonnyOrman/cloud-file-backup/storage/interface v0.0.0
	github.com/JonnyOrman/cloud-file-backup/storage/local v0.0.0
)

replace github.com/JonnyOrman/cloud-file-backup => ../../

replace github.com/JonnyOrman/cloud-file-backup/cloud/interface => ../../cloud/interface

replace github.com/JonnyOrman/cloud-file-backup/cloud/onedrive => ../../cloud/onedrive

replace github.com/JonnyOrman/cloud-file-backup/storage/interface => ../../storage/interface

replace github.com/JonnyOrman/cloud-file-backup/storage/local => ../../storage/local 