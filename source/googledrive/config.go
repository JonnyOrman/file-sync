package googledrive

// Config represents the configuration for Google Drive source service
type Config struct {
	ClientID     string `yaml:"client_id" mapstructure:"client_id"`
	ClientSecret string `yaml:"client_secret" mapstructure:"client_secret"`
} 