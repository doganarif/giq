package cmd

import (
	"github.com/doganarif/giq/internal/app"
	"github.com/spf13/cobra"
)

// NewRootCommand returns the root cobra.Command.
func NewRootCommand(a *app.App) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:           "giq",
		Short:         "giq - Quick Git operations enhanced with AI",
		SilenceErrors: true,
		SilenceUsage:  true,
		// The RunE here is only used when no subcommand is provided.
		RunE: func(cmd *cobra.Command, args []string) error {
			// With no subcommand, delegate to git (e.g., "giq --help").
			if len(args) == 0 {
				return a.ExecGit("--help")
			}
			return a.ExecGit(args...)
		},
	}

	// Register giq-specific subcommands.
	rootCmd.AddCommand(NewCommitCommand(a))
	rootCmd.AddCommand(NewStatusCommand(a))
	rootCmd.AddCommand(NewSetupCommand())

	return rootCmd
}
