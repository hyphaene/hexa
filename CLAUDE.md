# CLAUDE.md

## Base Instructions

Follow all instructions in ./AGENTS.md

## Claude Code Specific Context

This is a **Go learning project** - approach with expert pedagogy but stay direct.

**Learning Objective**: Monter en compétence sur l'écosystème, la syntaxe et les primitives de Go.

**Teaching Approach**: Expert Go guidance, pedagogical but direct, no hand-holding.

## Key Go Learning Focus Areas

When working on this project, emphasize learning these Go concepts:

### Packages and Modules
- `go.mod` definition and dependency management
- Import paths: `github.com/hyphaene/hexa/cmd`
- Package main as entry point with `func main()`

### Cobra CLI Patterns
- **Root command**: `cmd/root.go` with `&cobra.Command{}`
- **Subcommands**: `rootCmd.AddCommand()` in `init()`
- **Flags**: Persistent vs local flags
- **Help**: Automatic generation with descriptions

### Future Go Concepts (as project evolves)
- **Embedding**: `//go:embed scripts/*.sh` for script inclusion
- **Viper Config**: YAML files + environment variables + flags precedence

## Development Philosophy

- **NEVER** generate code without explicit request
- **ALWAYS** preserve existing Go conventions
- **FOCUS** on Go primitives learning over speed
- **APPROACH**: Direct pedagogy, explain Go concepts as we build

## Repository Context

- **Homebrew Tap**: `~/Code/homebrew-hexa` (separate repo)
- **GitHub**: [github.com/hyphaene/homebrew-hexa](https://github.com/hyphaene/homebrew-hexa)
- **Distribution**: GoReleaser auto-updates Homebrew tap
