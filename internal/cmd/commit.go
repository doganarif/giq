package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/doganarif/giq/internal/ai"
	"github.com/doganarif/giq/internal/app"
	"github.com/spf13/cobra"

	tea "github.com/charmbracelet/bubbletea"
)

// fallbackModel is a simple model for when AI is not configured
type fallbackModel struct {
	cursor   int
	options  []string
	selected int
}

func initialFallbackModel() fallbackModel {
	return fallbackModel{
		options: []string{
			"Write custom commit message",
			"Setup AI configuration",
		},
		selected: -1,
	}
}

func (m fallbackModel) Init() tea.Cmd {
	return nil
}

func (m fallbackModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.options)-1 {
				m.cursor++
			}
		case "enter":
			m.selected = m.cursor
			return m, tea.Quit
		case "q", "ctrl+c":
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m fallbackModel) View() string {
	s := "API key is not configured. Please choose an option:\n\n"
	for i, option := range m.options {
		cursor := "  "
		if m.cursor == i {
			cursor = "> "
		}
		s += fmt.Sprintf("%s%s\n", cursor, option)
	}
	s += "\nUse ↑/↓ arrows to navigate, enter to select"
	return s
}

// handleUnconfiguredAPI presents options to the user when API is not configured
func handleUnconfiguredAPI(a *app.App) (string, error) {
	p := tea.NewProgram(initialFallbackModel())
	m, err := p.Run()
	if err != nil {
		return "", err
	}

	fm, ok := m.(fallbackModel)
	if !ok {
		return "", fmt.Errorf("unexpected model type")
	}

	if fm.selected == -1 {
		return "", fmt.Errorf("no option selected")
	}

	if fm.selected == 0 {
		// Get custom message
		fmt.Print("\nEnter your commit message: ")
		reader := bufio.NewReader(os.Stdin)
		message, err := reader.ReadString('\n')
		if err != nil {
			return "", err
		}
		// Use system git for commit to ensure signing config is respected
		message = strings.TrimSpace(message)
		if err := a.ExecGit("commit", "-m", message); err != nil {
			return "", err
		}
		return "", nil
	} else {
		// Run setup
		if _, err := RunSetup(); err != nil {
			return "", fmt.Errorf("setup failed: %w", err)
		}
		return "", fmt.Errorf("setup completed, please try committing again")
	}
}

// commitModel is the model for selecting AI-generated commit messages
type commitModel struct {
	choices  []string
	cursor   int
	selected int
}

func initialCommitModel(choices []string) commitModel {
	return commitModel{
		choices:  choices,
		selected: -1,
	}
}

func (m commitModel) Init() tea.Cmd {
	return nil
}

func (m commitModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}
		case "enter":
			m.selected = m.cursor
			return m, tea.Quit
		case "q", "ctrl+c":
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m commitModel) View() string {
	s := "Select a commit message:\n\n"
	for i, choice := range m.choices {
		cursor := "  "
		if m.cursor == i {
			cursor = "> "
		}
		s += fmt.Sprintf("%s%s\n", cursor, choice)
	}
	s += "\nUse ↑/↓ arrows to navigate, enter to select"
	return s
}

// NewCommitCommand creates the commit command with AI-enhanced commit message support
func NewCommitCommand(a *app.App) *cobra.Command {
	var message string
	cmd := &cobra.Command{
		Use:   "commit",
		Short: "Create a commit with an AI-generated message from staged changes",
		RunE: func(cmd *cobra.Command, args []string) error {
			// If message flag is provided, use it directly with system git
			if message != "" {
				return a.ExecGit("commit", "-m", message)
			}

			// Show staged files
			stagedFiles, err := a.GetStagedFiles()
			if err != nil {
				return err
			}
			fmt.Println("Staged files:")
			fmt.Println(strings.TrimSpace(stagedFiles))
			fmt.Println("----------")

			// Get the diff
			diff, err := a.GetDiff()
			if err != nil {
				return err
			}
			if strings.TrimSpace(diff) == "" {
				return fmt.Errorf("no staged changes detected")
			}

			// Try to generate AI suggestions
			prompt := fmt.Sprintf(
				"Generate a single line, concise, and descriptive git commit message summarizing the staged changes on the following files: %s. "+
					"Do not include bullet points, extra formatting, or multiple lines. Diff:\n%s",
				strings.TrimSpace(stagedFiles), diff,
			)

			suggestions, err := ai.GenerateCommitMessages(a.Config, prompt)
			if err != nil {
				// Handle unconfigured API case
				if strings.Contains(err.Error(), "API key is not configured") {
					_, err := handleUnconfiguredAPI(a)
					return err
				}
				return err
			}

			// Add custom message option
			suggestions = append(suggestions, "Write custom message")

			// Show commit message selection UI
			p := tea.NewProgram(initialCommitModel(suggestions))
			m, err := p.Run()
			if err != nil {
				return err
			}

			cm, ok := m.(commitModel)
			if !ok {
				return fmt.Errorf("unexpected model type")
			}

			if cm.selected == -1 {
				return fmt.Errorf("no commit message selected")
			}

			commitMsg := suggestions[cm.selected]

			// Handle custom message option
			if commitMsg == "Write custom message" {
				fmt.Print("\nEnter your commit message: ")
				reader := bufio.NewReader(os.Stdin)
				customMsg, err := reader.ReadString('\n')
				if err != nil {
					return err
				}
				commitMsg = strings.TrimSpace(customMsg)
			}

			return a.ExecGit("commit", "-m", commitMsg)
		},
	}

	cmd.Flags().StringVarP(&message, "message", "m", "", "Commit message (overrides AI suggestions)")
	return cmd
}
