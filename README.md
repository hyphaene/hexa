# Hexa CLI

Hexactitude CLI - Unified automation and scripting toolkit

## Installation

### Via Homebrew (recommended)

```bash
# Add the tap first
brew tap hyphaene/hexa

# Then install hexa
brew install hexa

# Verify installation with both commands
hexa --help
hw --help   # Short alias
```

### Manual Installation

Download the latest release from the [releases page](https://github.com/hyphaene/hexa/releases).

## Configuration

Hexa CLI supports a flexible multi-level configuration system with smart template expansion and automatic security features. Use environment variables for sensitive data (tokens, passwords) without exposing them in config files.

**📖 [Full Configuration Guide](docs/configuration.md)** | **📚 [Viper Documentation](https://github.com/spf13/viper#working-with-environment-variables)**

### Quick Start

```bash
# Create local config (gitignored)
cat > .hexa.local.yml << 'EOF'
user:
  email: "your@email.com"
jira:
  token: "${HEXA_JIRA_TOKEN}"  # Use env var for secrets
EOF

# Set sensitive values via environment variables
export HEXA_JIRA_TOKEN="your-secret-token"

# Or use .env file (automatically loaded, gitignored)
echo 'HEXA_JIRA_TOKEN=your-secret-token' > .env

# View current configuration
hexa config
```

### Configuration Hierarchy (priority order)
1. CLI flags → 2. Environment variables → 3. Project local secrets → 4. Project config → 5. User global config → 6. Defaults

### Environment Variables & .env Support
- **Environment variables**: Use `HEXA_` prefix (e.g., `HEXA_JIRA_TOKEN`)
- **Auto .env loading**: Place `.env` file in working directory
- **Security**: Sensitive data stays out of config files

## Development

### Local Build

```bash
# Build the binary
go build -o hexa

# Run the local binary
./hexa --help
./hexa version
```

### Development Workflow

```bash
# Run directly without building
go run main.go [args]

# Run tests
go test ./...

# Format code
go fmt ./...

# Tidy dependencies
go mod tidy
```

## Homebrew Tap

This project uses a custom Homebrew tap for distribution. The tap repository is maintained at:
[github.com/hyphaene/homebrew-hexa](https://github.com/hyphaene/homebrew-hexa)
