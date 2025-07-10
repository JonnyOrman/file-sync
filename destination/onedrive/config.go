package onedrive

// Config represents the configuration for OneDrive destination service
type Config struct {
	ClientID     string `yaml:"client_id" mapstructure:"client_id"`
	ClientSecret string `yaml:"client_secret" mapstructure:"client_secret"`
	TenantID     string `yaml:"tenant_id" mapstructure:"tenant_id"`
} 