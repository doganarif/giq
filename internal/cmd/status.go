package cmd

import (
	"fmt"
	"strings"

	"github.com/doganarif/giq/internal/ai"
	"github.com/doganarif/giq/internal/app"
	"github.com/spf13/cobra"
)

// NewStatusCommand creates the status command which shows the working tree status
// and AI-generated insights regarding the changes.
func NewStatusCommand(a *app.App) *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "Show working tree status with AI insights",
		RunE: func(cmd *cobra.Command, args []string) error {
			// If not in a git repository, delegate to system git.
			if a.Repo == nil {
				return a.ExecGit("status")
			}

			// Get the standard git status.
			w, err := a.Repo.Worktree()
			if err != nil {
				return err
			}
			status, err := w.Status()
			if err != nil {
				return err
			}
			statusStr := status.String()
			fmt.Println("Git status:")
			fmt.Println(statusStr)

			// Get the diff of staged changes.
			diff, err := a.GetDiff()
			if err != nil {
				return err
			}
			if len(strings.TrimSpace(diff)) == 0 {
				fmt.Println("\nNo staged changes to analyze for AI insights.")
				return nil
			}

			// Generate AI insights based on the diff output.
			insights, err := ai.GenerateStatusInsights(a.Config, diff)
			if err != nil {
				fmt.Println("\n[Warning: Could not generate AI insights]")
			} else {
				fmt.Println("\nAI insights:")
				fmt.Println(insights)
			}
			return nil
		},
	}
}
