package app

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/doganarif/giq/internal/ai"
	"github.com/doganarif/giq/internal/config"
	"github.com/go-git/go-git/v5"
)

// App holds configuration, repository and the git executable path.
type App struct {
	Config *config.Config
	Repo   *git.Repository
	GitCmd string
}

// New creates a new App instance by locating the git executable,
// loading configuration, and opening the git repository (if any).
func New() (*App, error) {
	// Find the git executable.
	gitPath, err := exec.LookPath("git")
	if err != nil {
		return nil, fmt.Errorf("git not found in PATH: %w", err)
	}

	// Load configuration (defaults are used if no config file is found).
	cfg, _ := config.Load()

	// Attempt to open the git repository.
	repo, err := git.PlainOpenWithOptions(".", &git.PlainOpenOptions{
		DetectDotGit: true,
	})
	if err != nil && !errors.Is(err, git.ErrRepositoryNotExists) {
		return nil, fmt.Errorf("opening repository: %w", err)
	}

	return &App{
		Config: cfg,
		Repo:   repo,
		GitCmd: gitPath,
	}, nil
}

// ExecGit delegates execution to the system's git executable.
func (a *App) ExecGit(args ...string) error {
	cmd := exec.Command(a.GitCmd, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd.Run()
}

// GetDiff collects the diff for staged changes using the system git.
func (a *App) GetDiff() (string, error) {
	if a.Repo == nil {
		return "", fmt.Errorf("not a git repository")
	}

	// Use "git diff --cached" to get the diff of staged changes.
	cmd := exec.Command(a.GitCmd, "diff", "--cached")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	return string(output), nil
}

// GetStagedFiles retrieves a list of staged file names.
func (a *App) GetStagedFiles() (string, error) {
	if a.Repo == nil {
		return "", fmt.Errorf("not a git repository")
	}

	// Use "git diff --cached --name-only" to get the staged file names.
	cmd := exec.Command(a.GitCmd, "diff", "--cached", "--name-only")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	return string(output), nil
}

// GenerateCommitMessage creates a commit message based on the diff and staged file names
// by calling an AI service.
func (a *App) GenerateCommitMessage() (string, error) {
	diff, err := a.GetDiff()
	if err != nil {
		return "", err
	}

	files, err := a.GetStagedFiles()
	if err != nil {
		return "", err
	}

	trimmedFiles := strings.TrimSpace(files)
	if trimmedFiles == "" {
		return "No changes detected.", nil
	}
	if strings.TrimSpace(diff) == "" {
		return "No changes detected.", nil
	}

	// Build a prompt that includes the list of staged files.
	prompt := fmt.Sprintf(
		"Generate a single line, concise, and descriptive git commit message summarizing the staged changes on the following files: %s. "+
			"Do not include bullet points, extra formatting, or multiple lines. Diff:\n%s",
		trimmedFiles, diff,
	)

	// Use the AI package to generate the commit message.
	commitMsg, err := ai.GenerateCommitMessage(a.Config, prompt)
	if err != nil {
		return "", err
	}
	return commitMsg, nil
}
