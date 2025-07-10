module github.com/JonnyOrman/cloud-file-backup/destination/dropbox

go 1.21

require (
	github.com/JonnyOrman/cloud-file-backup/destination/interface v0.0.0
	github.com/dropbox/dropbox-sdk-go-unofficial/v6 v6.0.5
)

require (
	github.com/golang/protobuf v1.4.2 // indirect
	golang.org/x/net v0.0.0-20200822124328-c89045814202 // indirect
	golang.org/x/oauth2 v0.0.0-20201208152858-08078c50e5b5 // indirect
	google.golang.org/appengine v1.6.6 // indirect
	google.golang.org/protobuf v1.25.0 // indirect
)

replace github.com/JonnyOrman/cloud-file-backup/destination/interface => ../interface 