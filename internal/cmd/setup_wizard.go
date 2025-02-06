package cmd

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/doganarif/giq/internal/config"
)

type setupStep int

const (
	stepProviderSelection setupStep = iota
	stepInputCredential
	stepConfirm
	stepDone
)

type setupModel struct {
	step         setupStep
	provider     string
	currentField string
	fieldIndex   int
	input        textinput.Model
	answers      map[string]string
	done         bool
	err          error
}

func getFieldsForProvider(provider string) []string {
	if provider == "openai" {
		return []string{"OpenAI API Key"}
	}
	return []string{
		"Azure Endpoint",
		"Azure Deployment ID",
		"Azure API Key",
		"Azure API Version",
	}
}

func initialSetupModel() setupModel {
	input := textinput.New()
	input.Width = 60
	input.Focus()

	return setupModel{
		step:    stepProviderSelection,
		input:   input,
		answers: make(map[string]string),
	}
}

func (m setupModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m setupModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch m.step {
		case stepProviderSelection:
			switch msg.String() {
			case "1":
				m.provider = "openai"
				m.step = stepInputCredential
				fields := getFieldsForProvider(m.provider)
				m.currentField = fields[0]
				m.input.Placeholder = "Enter " + m.currentField
				return m, nil
			case "2":
				m.provider = "azure_openai"
				m.step = stepInputCredential
				fields := getFieldsForProvider(m.provider)
				m.currentField = fields[0]
				m.input.Placeholder = "Enter " + m.currentField
				return m, nil
			case "q", "ctrl+c":
				return m, tea.Quit
			}

		case stepInputCredential:
			switch msg.String() {
			case "enter":
				if strings.TrimSpace(m.input.Value()) == "" {
					return m, nil
				}

				// Save the current input
				m.answers[m.currentField] = m.input.Value()

				// Clear the input for the next field
				m.input.SetValue("")

				// Move to next field or confirm
				fields := getFieldsForProvider(m.provider)
				m.fieldIndex++

				if m.fieldIndex >= len(fields) {
					m.step = stepConfirm
					return m, nil
				}

				m.currentField = fields[m.fieldIndex]
				m.input.Placeholder = "Enter " + m.currentField
				return m, nil

			case "esc":
				m.step = stepProviderSelection
				m.fieldIndex = 0
				m.input.SetValue("")
				return m, nil
			}

		case stepConfirm:
			switch msg.String() {
			case "y", "Y", "enter":
				m.step = stepDone
				m.done = true
				return m, tea.Quit
			case "n", "N":
				return initialSetupModel(), nil
			case "esc":
				m.step = stepProviderSelection
				m.fieldIndex = 0
				m.input.SetValue("")
				return m, nil
			}
		}
	}

	if m.step == stepInputCredential {
		m.input, cmd = m.input.Update(msg)
	}

	return m, cmd
}

func (m setupModel) View() string {
	var s strings.Builder

	switch m.step {
	case stepProviderSelection:
		s.WriteString("Select AI Provider:\n\n")
		s.WriteString("1. OpenAI\n")
		s.WriteString("2. Azure OpenAI\n\n")
		s.WriteString("Press 1 or 2 to select a provider (ESC to cancel)\n")

	case stepInputCredential:
		s.WriteString(fmt.Sprintf("Setting up %s\n\n", strings.Title(m.provider)))
		fields := getFieldsForProvider(m.provider)

		// Show progress
		for i, field := range fields {
			if i < m.fieldIndex {
				s.WriteString(fmt.Sprintf("✓ %s\n", field))
			} else if i == m.fieldIndex {
				s.WriteString(fmt.Sprintf("\n> %s:\n", field))
				s.WriteString(m.input.View())
				s.WriteString("\n\nPress ENTER to confirm (ESC to start over)\n")
			} else {
				s.WriteString(fmt.Sprintf("  %s\n", field))
			}
		}

	case stepConfirm:
		s.WriteString("Please confirm your configuration:\n\n")
		s.WriteString(fmt.Sprintf("AI Provider: %s\n", strings.Title(m.provider)))
		for field, value := range m.answers {
			masked := strings.Contains(strings.ToLower(field), "key") || strings.Contains(strings.ToLower(field), "token")
			if masked {
				s.WriteString(fmt.Sprintf("%s: %s%s\n", field, value[:4], strings.Repeat("*", len(value)-4)))
			} else {
				s.WriteString(fmt.Sprintf("%s: %s\n", field, value))
			}
		}
		s.WriteString("\nConfirm? (Y/n): ")

	case stepDone:
		s.WriteString("Configuration saved successfully!")
	}

	return s.String()
}

func RunSetup() (*config.Config, error) {
	p := tea.NewProgram(initialSetupModel())
	m, err := p.Run()
	if err != nil {
		return nil, fmt.Errorf("setup error: %w", err)
	}

	finalModel, ok := m.(setupModel)
	if !ok {
		return nil, fmt.Errorf("unexpected model type")
	}

	if !finalModel.done {
		return nil, fmt.Errorf("setup cancelled")
	}

	// Build the config from the collected answers
	cfg := &config.Config{
		AIProvider: finalModel.provider,
	}

	if cfg.AIProvider == "openai" {
		cfg.AIKey = finalModel.answers["OpenAI API Key"]
	} else {
		cfg.AzureEndpoint = finalModel.answers["Azure Endpoint"]
		cfg.AzureDeploymentID = finalModel.answers["Azure Deployment ID"]
		cfg.AzureAPIKey = finalModel.answers["Azure API Key"]
		cfg.AzureAPIVersion = finalModel.answers["Azure API Version"]
	}

	return cfg, nil
}
