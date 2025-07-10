module github.com/JonnyOrman/cloud-file-backup/examples/googledrive-s3

go 1.20

require (
	github.com/JonnyOrman/cloud-file-backup v0.0.0
	github.com/JonnyOrman/cloud-file-backup/cloud/interface v0.0.0
	github.com/JonnyOrman/cloud-file-backup/cloud/googledrive v0.0.0
	github.com/JonnyOrman/cloud-file-backup/storage/interface v0.0.0
	github.com/JonnyOrman/cloud-file-backup/storage/s3 v0.0.0
)

replace github.com/JonnyOrman/cloud-file-backup => ../../

replace github.com/JonnyOrman/cloud-file-backup/cloud/interface => ../../cloud/interface

replace github.com/JonnyOrman/cloud-file-backup/cloud/googledrive => ../../cloud/googledrive

replace github.com/JonnyOrman/cloud-file-backup/storage/interface => ../../storage/interface

replace github.com/JonnyOrman/cloud-file-backup/storage/s3 => ../../storage/s3 