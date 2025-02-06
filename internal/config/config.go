package config

import (
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

// Config holds the configuration values for the giq application.
type Config struct {
	AIProvider        string `mapstructure:"ai_provider"`
	AIKey             string `mapstructure:"ai_key"`
	AzureEndpoint     string `mapstructure:"azure_endpoint"`
	AzureDeploymentID string `mapstructure:"azure_deployment_id"`
	AzureAPIKey       string `mapstructure:"azure_api_key"`
	AzureAPIVersion   string `mapstructure:"azure_api_version"`
}

// Load reads configuration from common config file locations and environment variables.
// If no config file is found, it creates one in $HOME/.config/giq/config.yaml with extended commented instructions.
func Load() (*Config, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	// Search in current directory, then in $HOME/.config/giq and $HOME/.giq.
	v.AddConfigPath(".")
	v.AddConfigPath(filepath.Join(home, ".config", "giq"))
	v.AddConfigPath(filepath.Join(home, ".giq"))

	// Environment variables prefixed with GIQ_ are also read.
	v.SetEnvPrefix("GIQ")
	v.AutomaticEnv()

	// Set default values.
	v.SetDefault("ai_provider", "openai")

	// Attempt to read the config file.
	err = v.ReadInConfig()
	if err != nil {
		// If no config file is found, create one with instructions.
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			configDir := filepath.Join(home, ".config", "giq")
			configFile := filepath.Join(configDir, "config.yaml")
			// Create the config directory if it doesn't exist.
			if mkErr := os.MkdirAll(configDir, 0755); mkErr != nil {
				return nil, mkErr
			}
			// Extended default content with commented instructions.
			defaultContent := `# giq configuration file
#
# This file configures the AI provider for giq.
#
# ai_provider: Specify the AI provider to use. Options include:
#    - openai (default)
#    - azure_openai
#
# For OpenAI, configure the following:
#   ai_key: Your API key for OpenAI.
#
# For Azure OpenAI, configure the following:
#   azure_endpoint: The endpoint for your Azure OpenAI resource
#                   (e.g., https://your-resource-name.openai.azure.com/)
#   azure_deployment_id: The deployment ID for the OpenAI model.
#   azure_api_key: Your API key for Azure OpenAI.
#   azure_api_version: The API version for Azure OpenAI (e.g., 2022-12-01).
#
# Example configuration for OpenAI:
#
#   ai_provider: openai
#   ai_key: your-openai-api-key
#
# Example configuration for Azure OpenAI:
#
#   ai_provider: azure_openai
#   azure_endpoint: https://your-resource-name.openai.azure.com/
#   azure_deployment_id: your-deployment-id
#   azure_api_key: your-azure-api-key
#   azure_api_version: 2022-12-01
#
`
			_ = os.WriteFile(configFile, []byte(defaultContent), 0644)
		} else {
			return nil, err
		}
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
