package gcs

// Config represents the configuration for Google Cloud Storage destination service
type Config struct {
	GCSBucket  string `yaml:"gcs_bucket" mapstructure:"gcs_bucket"`
	GCSProject string `yaml:"gcs_project" mapstructure:"gcs_project"`
	GCSPrefix  string `yaml:"gcs_prefix" mapstructure:"gcs_prefix"`
	GCSKeyFile string `yaml:"gcs_key_file" mapstructure:"gcs_key_file"`
} 