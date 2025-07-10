module github.com/JonnyOrman/cloud-file-backup/examples/onedrive-onedrive

go 1.20

require (
	github.com/JonnyOrman/cloud-file-backup v0.0.0
	github.com/JonnyOrman/cloud-file-backup/source/interface v0.0.0
	github.com/JonnyOrman/cloud-file-backup/source/onedrive v0.0.0
	github.com/JonnyOrman/cloud-file-backup/destination/interface v0.0.0
	github.com/JonnyOrman/cloud-file-backup/destination/onedrive v0.0.0
)

replace github.com/JonnyOrman/cloud-file-backup => ../../

replace github.com/JonnyOrman/cloud-file-backup/source/interface => ../../source/interface

replace github.com/JonnyOrman/cloud-file-backup/source/onedrive => ../../source/onedrive

replace github.com/JonnyOrman/cloud-file-backup/destination/interface => ../../destination/interface

replace github.com/JonnyOrman/cloud-file-backup/destination/onedrive => ../../destination/onedrive 