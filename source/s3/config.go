package s3

// Config represents the configuration for AWS S3 source service
type Config struct {
	S3Bucket string `yaml:"s3_bucket" mapstructure:"s3_bucket"`
	S3Region string `yaml:"s3_region" mapstructure:"s3_region"`
	S3Prefix string `yaml:"s3_prefix" mapstructure:"s3_prefix"`
} 