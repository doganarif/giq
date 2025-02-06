package ai

import (
	"context"
	"fmt"
	"strings"

	openai "github.com/sashabaranov/go-openai"

	"github.com/doganarif/giq/internal/config"
)

// GenerateCommitMessage is a convenience wrapper that returns a single commit message suggestion.
// It calls GenerateCommitMessages and returns the first suggestion.
func GenerateCommitMessage(cfg *config.Config, prompt string) (string, error) {
	suggestions, err := GenerateCommitMessages(cfg, prompt)
	fmt.Println(prompt)
	if err != nil {
		return "", err
	}
	if len(suggestions) == 0 {
		return "", fmt.Errorf("no commit message suggestions returned")
	}
	return suggestions[0], nil
}

// GenerateCommitMessages selects the appropriate provider based on configuration
// and returns a slice of generated commit message suggestions.
func GenerateCommitMessages(cfg *config.Config, prompt string) ([]string, error) {
	if strings.ToLower(cfg.AIProvider) == "azure_openai" {
		return generateCommitMessagesAzure(cfg, prompt)
	}
	return generateCommitMessagesOpenAI(cfg, prompt)
}

// generateCommitMessagesOpenAI uses the official OpenAI SDK to generate commit messages
// via the Chat Completion API.
func generateCommitMessagesOpenAI(cfg *config.Config, prompt string) ([]string, error) {
	if cfg.AIKey == "" {
		return nil, fmt.Errorf("OpenAI API key is not configured")
	}

	client := openai.NewClient(cfg.AIKey)
	req := openai.ChatCompletionRequest{
		Model:       openai.GPT3Dot5Turbo, // Change as desired.
		Messages:    []openai.ChatCompletionMessage{{Role: openai.ChatMessageRoleUser, Content: prompt}},
		Temperature: 0.5,
		MaxTokens:   64,
		N:           3, // Request three completions.
	}

	resp, err := client.CreateChatCompletion(context.Background(), req)
	if err != nil {
		return nil, fmt.Errorf("OpenAI API error: %w", err)
	}

	if len(resp.Choices) == 0 {
		return nil, fmt.Errorf("no completions returned from OpenAI")
	}

	suggestions := make([]string, 0, len(resp.Choices))
	for _, choice := range resp.Choices {
		suggestions = append(suggestions, strings.TrimSpace(choice.Message.Content))
	}

	return suggestions, nil
}

// generateCommitMessagesAzure uses the official OpenAI SDK configured for Azure OpenAI.
func generateCommitMessagesAzure(cfg *config.Config, prompt string) ([]string, error) {
	// Ensure that all required Azure configuration values are provided.
	if cfg.AzureAPIKey == "" || cfg.AzureEndpoint == "" || cfg.AzureDeploymentID == "" || cfg.AzureAPIVersion == "" {
		return nil, fmt.Errorf("Azure OpenAI configuration is incomplete")
	}

	azureConfig := openai.DefaultAzureConfig(cfg.AzureAPIKey, cfg.AzureEndpoint)
	azureConfig.AzureModelMapperFunc = func(model string) string {
		azureModelMapping := map[string]string{
			// Map the desired model to your Azure deployment name.
			openai.GPT4o: cfg.AzureDeploymentID,
		}
		return azureModelMapping[model]
	}

	client := openai.NewClientWithConfig(azureConfig)
	req := openai.ChatCompletionRequest{
		Model:       openai.GPT4o, // Change as desired.
		Messages:    []openai.ChatCompletionMessage{{Role: openai.ChatMessageRoleUser, Content: prompt}},
		Temperature: 0.5,
		MaxTokens:   64,
		N:           3, // Request three completions.
	}

	resp, err := client.CreateChatCompletion(context.Background(), req)
	if err != nil {
		return nil, fmt.Errorf("Azure OpenAI API error: %w", err)
	}

	if len(resp.Choices) == 0 {
		return nil, fmt.Errorf("no completions returned from Azure OpenAI")
	}

	suggestions := make([]string, 0, len(resp.Choices))
	for _, choice := range resp.Choices {
		suggestions = append(suggestions, strings.TrimSpace(choice.Message.Content))
	}

	return suggestions, nil
}

// GenerateStatusInsights generates AI-based insights based on the diff output
// of the staged changes. It asks the AI to provide a concise summary describing
// what changed in each file.
func GenerateStatusInsights(cfg *config.Config, diff string) (string, error) {
	// Build a prompt that asks the AI to describe the changes (per file) based on the diff.
	prompt := fmt.Sprintf(
		"Based on the following git diff output for staged changes, provide a single concise sentence that summarizes the changes made to each file. "+
			"Indicate for each file whether code was added, removed, or modified, and if possible, what kind of changes occurred (for example, bug fixes, refactoring, or feature additions). "+
			"Do not simply list the file names. Diff:\n\n%s",
		diff,
	)

	if strings.ToLower(cfg.AIProvider) == "azure_openai" {
		return generateStatusInsightsAzure(cfg, prompt)
	}
	return generateStatusInsightsOpenAI(cfg, prompt)
}

func generateStatusInsightsOpenAI(cfg *config.Config, prompt string) (string, error) {
	if cfg.AIKey == "" {
		return "", fmt.Errorf("OpenAI API key is not configured")
	}
	client := openai.NewClient(cfg.AIKey)
	req := openai.ChatCompletionRequest{
		Model:       openai.GPT3Dot5Turbo,
		Messages:    []openai.ChatCompletionMessage{{Role: openai.ChatMessageRoleUser, Content: prompt}},
		Temperature: 0.5,
		MaxTokens:   128,
	}
	resp, err := client.CreateChatCompletion(context.Background(), req)
	if err != nil {
		return "", fmt.Errorf("OpenAI API error: %w", err)
	}
	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("no explanation returned from OpenAI")
	}
	return strings.TrimSpace(resp.Choices[0].Message.Content), nil
}

func generateStatusInsightsAzure(cfg *config.Config, prompt string) (string, error) {
	if cfg.AzureAPIKey == "" || cfg.AzureEndpoint == "" || cfg.AzureDeploymentID == "" || cfg.AzureAPIVersion == "" {
		return "", fmt.Errorf("Azure OpenAI configuration is incomplete")
	}
	azureConfig := openai.DefaultAzureConfig(cfg.AzureAPIKey, cfg.AzureEndpoint)
	azureConfig.AzureModelMapperFunc = func(model string) string {
		azureModelMapping := map[string]string{
			openai.GPT4o: cfg.AzureDeploymentID,
		}
		return azureModelMapping[model]
	}
	client := openai.NewClientWithConfig(azureConfig)
	req := openai.ChatCompletionRequest{
		Model:       openai.GPT4o,
		Messages:    []openai.ChatCompletionMessage{{Role: openai.ChatMessageRoleUser, Content: prompt}},
		Temperature: 0.5,
		MaxTokens:   128,
	}
	resp, err := client.CreateChatCompletion(context.Background(), req)
	if err != nil {
		return "", fmt.Errorf("Azure OpenAI API error: %w", err)
	}
	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("no explanation returned from Azure OpenAI")
	}
	return strings.TrimSpace(resp.Choices[0].Message.Content), nil
}
