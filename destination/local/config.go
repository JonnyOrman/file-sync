package local

// Config represents the configuration for local filesystem destination service
type Config struct {
	LocalPath string `yaml:"local_path" mapstructure:"local_path"`
} 