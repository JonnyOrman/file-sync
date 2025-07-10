package dropbox

// Config represents the configuration for Dropbox source service
type Config struct {
	AccessToken string `yaml:"access_token" mapstructure:"access_token"`
} 