# giq (Git Intelligence Quickstart)

giq is a Git wrapper that enhances your Git workflow with AI-powered features. It provides intelligent commit message suggestions and status insights while maintaining full compatibility with standard Git commands.

![demo](https://raw.githubusercontent.com/doganarif/giq/main/demo/demo.gif)

## Features

- **AI-Powered Commit Messages**: Automatically generates contextual commit messages based on your staged changes
- **Intelligent Status Insights**: Provides AI-enhanced analysis of your working tree status
- **Multi-Provider Support**: Works with both OpenAI and Azure OpenAI
- **Git Command Passthrough**: Seamlessly delegates unknown commands to your system's Git
- **Interactive Setup**: User-friendly configuration wizard

## Installation

### Using Homebrew (Recommended)

giq is available via Homebrew for macOS (both Apple Silicon and Intel) and Linux:

```bash
# Add the tap
brew tap doganarif/giq

# Install giq
brew install giq
```

To upgrade to the latest version:
```bash
brew upgrade giq
```

To uninstall:
```bash
brew uninstall giq
```

### Manual Installation

#### Prerequisites

- Go 1.19 or later
- Git installed and available in your PATH
- OpenAI API key or Azure OpenAI credentials

#### Building from Source

1. Clone the repository:
```bash
git clone https://github.com/doganarif/giq.git
cd giq
```

2. Build the binary:
```bash
go build
```

3. Move the binary to a location in your PATH:
```bash
# Linux/macOS
sudo mv giq /usr/local/bin/

# Windows
# Move giq.exe to a location in your PATH
```

### Installing Pre-built Binaries

Pre-built binaries for various platforms are available on the [releases page](https://github.com/doganarif/giq/releases).

Supported platforms:
- macOS (Apple Silicon/ARM64)
- macOS (Intel/AMD64)
- Linux (x86_64)
- Linux (ARM64)

## Configuration

Run the interactive setup wizard:

```bash
giq setup
```

This will guide you through:
1. Selecting your AI provider (OpenAI or Azure OpenAI)
2. Entering your API credentials
3. Saving the configuration

Configuration is stored in `~/.config/giq/config.yaml` (or equivalent on Windows).

### Manual Configuration

You can also create the configuration file manually:

For OpenAI:
```yaml
ai_provider: openai
ai_key: your-openai-api-key
```

For Azure OpenAI:
```yaml
ai_provider: azure_openai
azure_endpoint: https://your-resource.openai.azure.com/
azure_deployment_id: your-deployment-id
azure_api_key: your-azure-api-key
azure_api_version: 2022-12-01
```

## Usage

### Committing Changes

```bash
# Stage your changes as usual
git add .

# Generate AI-powered commit message
giq commit

# Or provide your own message
giq commit -m "your message"
```

When using `giq commit` without a message:
1. View staged files
2. Choose from AI-generated commit message suggestions
3. Or enter a custom message

### Checking Status

```bash
giq status
```

Shows:
- Standard Git status output
- AI-generated insights about your changes

### Other Git Commands

giq passes through any unrecognized commands to Git:

```bash
# These work exactly like standard git commands
giq push
giq pull
giq branch
# etc.
```

## Environment Variables

You can configure giq using environment variables:

- `GIQ_AI_PROVIDER`: AI provider (`openai` or `azure_openai`)
- `GIQ_AI_KEY`: OpenAI API key
- `GIQ_AZURE_ENDPOINT`: Azure OpenAI endpoint
- `GIQ_AZURE_DEPLOYMENT_ID`: Azure OpenAI deployment ID
- `GIQ_AZURE_API_KEY`: Azure OpenAI API key
- `GIQ_AZURE_API_VERSION`: Azure OpenAI API version

Environment variables take precedence over configuration file settings.

## Project Structure

```
giq/
├── internal/
│   ├── ai/      # AI service integration
│   ├── app/     # Core application logic
│   ├── cmd/     # Command implementations
│   └── config/  # Configuration management
└── main.go      # Application entry point
```

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Acknowledgments

- [go-git](https://github.com/go-git/go-git) for Git operations
- [cobra](https://github.com/spf13/cobra) for CLI interface
- [bubbletea](https://github.com/charmbracelet/bubbletea) for terminal UI
- [go-openai](https://github.com/sashabaranov/go-openai) for OpenAI integration
