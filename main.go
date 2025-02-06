package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/doganarif/giq/internal/app"
	"github.com/doganarif/giq/internal/cmd"
)

func main() {
	a, err := app.New()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error initializing giq: %v\n", err)
		os.Exit(1)
	}

	// Define the commands handled by giq.
	handledCommands := map[string]bool{
		"commit": true,
		"status": true,
		"help":   true,
		"setup":  true,
		"--help": true,
		"-h":     true,
	}

	// If there are arguments and the first argument is not one of our custom commands,
	// delegate directly to the system git executable.
	if len(os.Args) > 1 {
		if !handledCommands[os.Args[1]] {
			if err := a.ExecGit(os.Args[1:]...); err != nil {
				cmdStr := strings.Join(os.Args[1:], " ")
				fmt.Fprintf(os.Stderr, "Error executing 'git %s': %v\n", cmdStr, err)
				os.Exit(1)
			}
			return
		}
	}

	// Otherwise, use Cobra to process the command.
	rootCmd := cmd.NewRootCommand(a)
	if err := rootCmd.Execute(); err != nil {
		cmdStr := strings.Join(os.Args[1:], " ")
		fmt.Fprintf(os.Stderr, "Error executing 'git %s': %v\n", cmdStr, err)
		os.Exit(1)
	}
}
