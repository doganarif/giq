package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	_ "github.com/charmbracelet/bubbletea"
	_ "github.com/doganarif/giq/internal/config"
	"github.com/spf13/cobra"
)

// NewSetupCommand creates a new command for interactive configuration setup.
func NewSetupCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "setup",
		Short: "Interactive setup for giq configuration",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := RunSetup()
			if err != nil {
				return err
			}

			// Write config file to $HOME/.config/giq/config.yaml
			home, err := os.UserHomeDir()
			if err != nil {
				return err
			}
			configDir := filepath.Join(home, ".config", "giq")
			if err := os.MkdirAll(configDir, 0755); err != nil {
				return err
			}
			configFile := filepath.Join(configDir, "config.yaml")
			// Create YAML content based on provider.
			var content string
			if cfg.AIProvider == "openai" {
				content = fmt.Sprintf("ai_provider: openai\nai_key: %s\n", cfg.AIKey)
			} else {
				content = fmt.Sprintf(
					"ai_provider: azure_openai\nazure_endpoint: %s\nazure_deployment_id: %s\nazure_api_key: %s\nazure_api_version: %s\n",
					cfg.AzureEndpoint, cfg.AzureDeploymentID, cfg.AzureAPIKey, cfg.AzureAPIVersion,
				)
			}
			if err := os.WriteFile(configFile, []byte(content), 0644); err != nil {
				return err
			}
			fmt.Println("Configuration saved to", configFile)
			return nil
		},
	}
}
