package onedrive

// Config represents the configuration for OneDrive source service
type Config struct {
	ClientID string `yaml:"client_id" mapstructure:"client_id"`
	TenantID string `yaml:"tenant_id" mapstructure:"tenant_id"`
} 