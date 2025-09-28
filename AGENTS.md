# AGENTS.md

## Project Overview

CLI Go tool "hexa" (alias "hw") - Learning project for Go ecosystem fundamentals.
Goal: Replace 22+ bash scripts with a single distributable binary via Homebrew.

**Approach**: Expert Go guidance, direct pedagogy. Do NOT generate code unless explicitly requested.

## Development Environment Setup

```bash
# Prerequisites
go version  # Requires Go 1.24.4+

# Dependencies
go mod tidy
go mod download
```

## Code Style Guidelines

- Follow Go conventions: `go fmt`, `go vet`
- Preserve existing Cobra patterns in cmd/ structure
- No comments unless explicitly requested
- Focus on Go primitives learning
- Use existing libraries already in project (check imports first)

## Build and Test Commands

```bash
# Build
go build                    # Simple build
go build -o hexa           # Named binary

# Run
go run main.go [args]      # Direct run
./hexa --help              # Test binary
./hexa version             # Test version command

# Test
go test ./...              # Run all tests
go test -v ./...           # Verbose tests
go test -cover ./...       # With coverage
```

## Quality Checks

```bash
# Go toolchain
go fmt ./...               # Format code
go vet ./...               # Static analysis
go mod tidy               # Clean dependencies

# Optional (if available)
golangci-lint run         # Advanced linting
```

## Testing Instructions

- Always run `go test ./...` before commits
- Verify build works: `go build -o hexa && ./hexa --help`
- Test GoReleaser locally: `goreleaser release --snapshot --rm-dist`

## Release Process

```bash
# Local snapshot (testing)
goreleaser release --snapshot --rm-dist

# Production release (requires git tag)
goreleaser release --rm-dist
```

## Project Architecture

### Current Structure (Minimal)
```
hexa/
├── main.go                 # Entry point → cmd.Execute()
├── cmd/                    # Cobra commands
│   ├── root.go            # Root command (placeholder)
│   └── version.go         # Version command
├── scripts/               # Bash scripts to embed (empty)
├── internal/              # Internal packages (empty)
└── .goreleaser.yaml       # GoReleaser config
```

### Target Domains
- JIRA, GIT, SETUP, AI commands
- Framework: Cobra (commands + flags + help)
- Config: Viper (YAML + env vars) - future
- Embedding: `//go:embed` scripts - future
- Distribution: Homebrew with `hw` → `hexa` symlink

## Development Phases

1. **Phase 1**: Cobra structure + basic commands
2. **Phase 2**: Embed existing bash scripts
3. **Phase 3**: Progressive rewrite to pure Go
4. **Phase 4**: Viper config + optimizations

## Cobra Patterns

- Root command: `cmd/root.go` with `&cobra.Command{}`
- Subcommands: `rootCmd.AddCommand()` in `init()`
- Flags: Persistent vs local flags
- Help: Automatic with descriptions

## Distribution

### Development
```bash
go build -o hexa
./hexa --help
./hexa version
```

### Homebrew (Production)
```bash
brew tap hyphaene/hexa
brew install hexa
hw --help    # Alias for hexa
```

## Important Rules

- **NEVER** generate code without explicit request
- **ALWAYS** preserve existing Go conventions
- **FOCUS** on Go primitives learning
- **APPROACH**: Direct pedagogy, no hand-holding