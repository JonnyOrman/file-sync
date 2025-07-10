# Cloud File Backup - Modular Build System
# Each cloud service and storage backend is a standalone module

BINARY_NAME = cloud-file-backup
MAIN_FILE = main.go
GO_MOD = go.mod

# Build tag variables
CLOUD_TAG ?= onedrive
STORAGE_TAG ?= local

# Default target
.PHONY: all
all: help

# Initialize all modules
.PHONY: init-modules
init-modules:
	@echo "Initializing all modules..."
	@echo "Cloud interface module:"
	cd modules/cloud-interface && go mod tidy
	@echo "Storage interface module:"
	cd modules/storage-interface && go mod tidy
	@echo "Cloud service modules:"
	cd modules/cloud-onedrive && go mod tidy
	cd modules/cloud-googledrive && go mod tidy
	cd modules/cloud-dropbox && go mod tidy
	@echo "Storage service modules:"
	cd modules/storage-local && go mod tidy
	cd modules/storage-s3 && go mod tidy
	cd modules/storage-gcs && go mod tidy
	cd modules/storage-azblob && go mod tidy
	@echo "Main application:"
	go mod tidy -modfile=$(GO_MOD)

# Generic build target
.PHONY: build
build:
	@echo "Building with tags: cloud=$(CLOUD_TAG), storage=$(STORAGE_TAG)"
	go build -modfile=$(GO_MOD) -tags "$(CLOUD_TAG) $(STORAGE_TAG)" -o $(BINARY_NAME) $(MAIN_FILE)
	@echo "Binary size:"
	@ls -lh $(BINARY_NAME) | awk '{print $$5 " " $$9}'

# Quick build targets for common combinations
.PHONY: build-onedrive-local
build-onedrive-local:
	$(MAKE) build CLOUD_TAG=onedrive STORAGE_TAG=local

.PHONY: build-onedrive-s3
build-onedrive-s3:
	$(MAKE) build CLOUD_TAG=onedrive STORAGE_TAG=s3

.PHONY: build-googledrive-local
build-googledrive-local:
	$(MAKE) build CLOUD_TAG=googledrive STORAGE_TAG=local

.PHONY: build-googledrive-s3
build-googledrive-s3:
	$(MAKE) build CLOUD_TAG=googledrive STORAGE_TAG=s3

.PHONY: build-googledrive-gcs
build-googledrive-gcs:
	$(MAKE) build CLOUD_TAG=googledrive STORAGE_TAG=gcs

.PHONY: build-dropbox-local
build-dropbox-local:
	$(MAKE) build CLOUD_TAG=dropbox STORAGE_TAG=local

.PHONY: build-dropbox-s3
build-dropbox-s3:
	$(MAKE) build CLOUD_TAG=dropbox STORAGE_TAG=s3

.PHONY: build-dropbox-gcs
build-dropbox-gcs:
	$(MAKE) build CLOUD_TAG=dropbox STORAGE_TAG=gcs

.PHONY: build-dropbox-azblob
build-dropbox-azblob:
	$(MAKE) build CLOUD_TAG=dropbox STORAGE_TAG=azblob

# Build all common combinations
.PHONY: build-examples
build-examples: build-onedrive-local build-googledrive-s3 build-dropbox-gcs

# Build matrix - all combinations
.PHONY: build-all
build-all:
	@echo "Building all combinations..."
	@for cloud in onedrive googledrive dropbox; do \
		for storage in local s3 gcs azblob; do \
			echo "Building $$cloud + $$storage"; \
			$(MAKE) build CLOUD_TAG=$$cloud STORAGE_TAG=$$storage || exit 1; \
			mv $(BINARY_NAME) $(BINARY_NAME)-$$cloud-$$storage; \
		done; \
	done
	@echo "All builds completed:"
	@ls -lh $(BINARY_NAME)-* | awk '{print $$5 " " $$9}'

# Test builds (check compilation without creating binaries)
.PHONY: test-build
test-build:
	@echo "Testing build with tags: cloud=$(CLOUD_TAG), storage=$(STORAGE_TAG)"
	go build -modfile=$(GO_MOD) -tags "$(CLOUD_TAG) $(STORAGE_TAG)" -o /dev/null $(MAIN_FILE)
	@echo "Build test successful!"

# Test all combinations
.PHONY: test-all
test-all:
	@echo "Testing all build combinations..."
	@for cloud in onedrive googledrive dropbox; do \
		for storage in local s3 gcs azblob; do \
			echo "Testing $$cloud + $$storage"; \
			$(MAKE) test-build CLOUD_TAG=$$cloud STORAGE_TAG=$$storage || exit 1; \
		done; \
	done
	@echo "All build tests passed!"

# Clean up
.PHONY: clean
clean:
	rm -f $(BINARY_NAME) $(BINARY_NAME)-*
	rm -f .cloud-file-backup-state.json
	@echo "Cleaned up build artifacts"

# Run with specific configuration
.PHONY: run
run:
	@if [ ! -f $(BINARY_NAME) ]; then \
		echo "Binary not found. Building with current tags..."; \
		$(MAKE) build; \
	fi
	./$(BINARY_NAME) $(ARGS)

# Development helpers
.PHONY: fmt
fmt:
	go fmt ./...
	@for dir in modules/*/; do \
		echo "Formatting $$dir"; \
		cd "$$dir" && go fmt ./...; \
		cd ../..; \
	done

.PHONY: vet
vet:
	go vet -modfile=$(GO_MOD) ./...

.PHONY: deps
deps:
	@echo "Downloading dependencies for main application..."
	go mod download -modfile=$(GO_MOD)
	@echo "Downloading dependencies for all modules..."
	@for dir in modules/*/; do \
		echo "Downloading deps for $$dir"; \
		cd "$$dir" && go mod download; \
		cd ../..; \
	done

# Show available build tags and combinations
.PHONY: show-tags
show-tags:
	@echo "Available Cloud Services:"
	@echo "  onedrive    - Microsoft OneDrive"
	@echo "  googledrive - Google Drive"
	@echo "  dropbox     - Dropbox"
	@echo ""
	@echo "Available Storage Backends:"
	@echo "  local       - Local filesystem"
	@echo "  s3          - AWS S3"
	@echo "  gcs         - Google Cloud Storage"
	@echo "  azblob      - Azure Blob Storage"
	@echo ""
	@echo "Example builds:"
	@echo "  make build CLOUD_TAG=onedrive STORAGE_TAG=local"
	@echo "  make build-googledrive-s3"
	@echo "  make build-dropbox-gcs"

# Show module information
.PHONY: show-modules
show-modules:
	@echo "Standalone Modules:"
	@echo ""
	@echo "Interface Modules:"
	@echo "  modules/cloud-interface   - Cloud service interface"
	@echo "  modules/storage-interface - Storage service interface"
	@echo ""
	@echo "Cloud Service Modules:"
	@echo "  modules/cloud-onedrive    - OneDrive implementation"
	@echo "  modules/cloud-googledrive - Google Drive implementation"
	@echo "  modules/cloud-dropbox     - Dropbox implementation"
	@echo ""
	@echo "Storage Service Modules:"
	@echo "  modules/storage-local     - Local filesystem storage"
	@echo "  modules/storage-s3        - AWS S3 storage"
	@echo "  modules/storage-gcs       - Google Cloud Storage"
	@echo "  modules/storage-azblob    - Azure Blob Storage"
	@echo ""
	@echo "Factory System:"
	@echo "  factories/                - Build tag-based service factories"

# Help target
.PHONY: help
help:
	@echo "Cloud File Backup - Modular Build System"
	@echo ""
	@echo "Available targets:"
	@echo "  init-modules     - Initialize all Go modules"
	@echo "  build           - Build with specified tags (CLOUD_TAG, STORAGE_TAG)"
	@echo "  build-examples  - Build common combinations"
	@echo "  build-all       - Build all possible combinations"
	@echo "  test-build      - Test compilation without creating binary"
	@echo "  test-all        - Test all build combinations"
	@echo "  run             - Run the application (ARGS for arguments)"
	@echo "  clean           - Remove build artifacts"
	@echo "  fmt             - Format all code"
	@echo "  vet             - Run go vet"
	@echo "  deps            - Download all dependencies"
	@echo "  show-tags       - Show available build tags"
	@echo "  show-modules    - Show module information"
	@echo "  help            - Show this help"
	@echo ""
	@echo "Quick build targets:"
	@echo "  build-onedrive-local"
	@echo "  build-onedrive-s3"
	@echo "  build-googledrive-local"
	@echo "  build-googledrive-s3"
	@echo "  build-googledrive-gcs"
	@echo "  build-dropbox-local"
	@echo "  build-dropbox-s3"
	@echo "  build-dropbox-gcs"
	@echo "  build-dropbox-azblob"
	@echo ""
	@echo "Examples:"
	@echo "  make init-modules"
	@echo "  make build CLOUD_TAG=dropbox STORAGE_TAG=s3"
	@echo "  make build-googledrive-gcs"
	@echo "  make run ARGS='--help'" 