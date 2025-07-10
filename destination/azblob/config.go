package azblob

// Config represents the configuration for Azure Blob Storage destination service
type Config struct {
	AzureBlobAccount   string `yaml:"azure_blob_account" mapstructure:"azure_blob_account"`
	AzureBlobContainer string `yaml:"azure_blob_container" mapstructure:"azure_blob_container"`
	AzureBlobPrefix    string `yaml:"azure_blob_prefix" mapstructure:"azure_blob_prefix"`
	AzureBlobKey       string `yaml:"azure_blob_key" mapstructure:"azure_blob_key"`
} 