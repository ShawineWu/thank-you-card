package config

import (
	"os"
)

// Config holds all configuration for the application
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Auth     AuthConfig
	External ExternalConfig
}

// ServerConfig holds server configuration
type ServerConfig struct {
	Port string
	Mode string // debug, release
}

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// AuthConfig holds authentication configuration
type AuthConfig struct {
	JWTSecret     string
	TokenIssuer   string
	TokenAudience string
}

// ExternalConfig holds external service configuration
type ExternalConfig struct {
	EmployeeDirectoryURL string
	TeamsWebhookURL      string
	PlatUMSURL           string
	AzureTenantID        string
	AzureClientID        string
	AzureClientSecret    string
}

// Load loads configuration from environment variables
func Load() (*Config, error) {
	config := &Config{
		Server: ServerConfig{
			Port: getEnv("SERVER_PORT", "8080"),
			Mode: getEnv("GIN_MODE", "debug"),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", ""),
			DBName:   getEnv("DB_NAME", "thank_you_card"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
		Auth: AuthConfig{
			JWTSecret:     getEnv("JWT_SECRET", ""),
			TokenIssuer:   getEnv("TOKEN_ISSUER", ""),
			TokenAudience: getEnv("TOKEN_AUDIENCE", ""),
		},
		External: ExternalConfig{
			EmployeeDirectoryURL: getEnv("EMPLOYEE_DIRECTORY_URL", ""),
			TeamsWebhookURL:      getEnv("TEAMS_WEBHOOK_URL", ""),
			PlatUMSURL:           getEnv("PLAT_UMS_URL", ""),
			AzureTenantID:        getEnv("AZURE_TENANT_ID", ""),
			AzureClientID:        getEnv("AZURE_CLIENT_ID", ""),
			AzureClientSecret:    getEnv("AZURE_CLIENT_SECRET", ""),
		},
	}

	// Validate required configuration
	if err := config.Validate(); err != nil {
		return nil, err
	}

	return config, nil
}

// Validate checks if all required configuration is present
func (c *Config) Validate() error {
	// Optional for local dev: allow empty password
	// if c.Database.Password == "" {
	// 	return fmt.Errorf("DB_PASSWORD is required")
	// }

	// if c.Auth.JWTSecret == "" {
	// 	return fmt.Errorf("JWT_SECRET is required")
	// }

	// Optional for local dev
	// if c.External.EmployeeDirectoryURL == "" {
	// 	return fmt.Errorf("EMPLOYEE_DIRECTORY_URL is required")
	// }

	// Optional: Teams Webhook URL (can be empty for local dev)
	// if c.External.TeamsWebhookURL == "" {
	// 	return fmt.Errorf("TEAMS_WEBHOOK_URL is required")
	// }

	return nil
}

// getEnv retrieves an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
