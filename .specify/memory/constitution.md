<!--
SYNC IMPACT REPORT
==================
Version Change: No prior version → 1.0.0

Modified Principles: N/A (initial constitution)

Added Sections:
- Core Principles (5 principles)
- Development Standards
- Quality Gates
- Governance

Removed Sections: N/A (initial constitution)

Templates Requiring Updates:
- ✅ plan-template.md: Updated constitution reference to v1.0.0
- ✅ spec-template.md: No changes needed (tech-agnostic requirements focus compatible)
- ✅ tasks-template.md: No changes needed (TDD and parallel execution align with principles)
- ✅ agent-file-template.md: Not found - likely not needed (AGENTS.md and CLAUDE.md already exist in repo root)

Follow-up TODOs: None
-->

# Hexa CLI Constitution

## Core Principles

### I. Go Learning First
**Every feature serves the learning objective of mastering Go fundamentals.**

- Code demonstrates Go primitives, idioms, and standard patterns
- Complexity introduced only when it teaches valuable Go concepts
- Direct pedagogy: explain WHY Go works this way, not just HOW
- Prefer Go standard library over third-party when both teach equally
- Document learning focus in commit messages and PR descriptions

**Rationale**: This is explicitly a Go learning project. All technical decisions must align with the pedagogical objective of building Go expertise through practical CLI development.

### II. CLI-First Architecture
**Every feature exposes functionality through a clean CLI interface using Cobra patterns.**

- Commands follow Cobra structure: `cmd/root.go` + subcommands
- Flags use persistent vs local appropriately
- Help text is automatically generated and descriptive
- STDIN/args → STDOUT pattern, errors → STDERR
- Support both human-readable and machine-parseable output formats

**Rationale**: The tool's primary interface is the command line. A well-structured CLI using Go's established Cobra framework both serves users and teaches CLI design patterns.

### III. Test-Driven Development (NON-NEGOTIABLE)
**Tests written → User approved → Tests fail → Then implement.**

- Write tests first, verify they fail, then implement to make them pass
- Use Go's standard testing package: `go test ./...`
- Coverage targets: Core logic >80%, CLI commands >60%
- Integration tests for Jira/GitHub/external API interactions
- Red-Green-Refactor cycle strictly enforced

**Rationale**: TDD is non-negotiable because it forces design clarity, provides regression safety, and teaches test-first thinking—a core discipline in professional Go development.

### IV. Distribution Simplicity
**Single binary distributed via Homebrew, no runtime dependencies.**

- GoReleaser handles builds, releases, and Homebrew tap updates
- Binary includes embedded assets via `//go:embed` when needed
- Version command reflects git tags (semantic versioning)
- Alias `hw` → `hexa` for ergonomics
- Cross-platform builds (macOS, Linux) from day one

**Rationale**: Simplicity in distribution reduces friction for users and teaches Go's compilation model. Single-binary deployment is a Go superpower worth emphasizing.

### V. Configuration Hierarchy
**Viper-driven config with clear precedence: CLI flags → env vars → files → defaults.**

- Environment variables prefixed with `HEXA_`
- Auto-load `.env` files with `godotenv`
- YAML config files: `.hexa.local.yml` (gitignored secrets) + `.hexa.yml` (project)
- Placeholders like `${HEXA_JIRA_TOKEN}` resolved at runtime
- `hexa config` command shows effective configuration

**Rationale**: Multi-level configuration teaches Viper patterns and keeps secrets out of repos while maintaining developer ergonomics.

## Development Standards

### Code Conventions
- Run `go fmt ./...`, `go vet ./...`, and `golangci-lint run` before commits
- Preserve existing Cobra patterns in `cmd/` structure
- No comments unless explicitly teaching a non-obvious Go concept
- Package names: lowercase, single-word, descriptive
- Error handling: explicit, not ignored (`if err != nil` checks required)

### Repository Hygiene
- Commit messages: conventional commits format (`feat:`, `fix:`, `docs:`, `refactor:`)
- Branch naming: `feat/description`, `fix/issue-number`
- PR descriptions must explain what Go concept is being learned/applied
- No generated files committed (binaries, build artifacts)
- `.gitignore` includes: `hexa` binary, `.hexa.local.yml`, `.env`

### External Dependencies
- Minimize dependencies: Go standard library preferred
- Existing approved libraries: Cobra (CLI), Viper (config), godotenv (env), go-jira (API)
- New dependencies require justification: "What Go concept does this teach?"
- `go mod tidy` after dependency changes

## Quality Gates

### Pre-Commit Checks
1. `go test ./...` passes
2. `go build` succeeds
3. `golangci-lint run` reports no errors
4. No TODOs or FIXME comments without linked issues

### Pre-Release Checks
1. All tests pass
2. Version bumped in git tag (semantic versioning)
3. `goreleaser release --snapshot --rm-dist` succeeds locally
4. CHANGELOG updated with user-facing changes

### Learning Validation
- Each feature must document what Go concepts it teaches
- Code reviews verify pedagogical value, not just correctness
- Prefer readability over cleverness (this is a learning project)

## Governance

### Amendment Process
1. Proposed changes documented in PR with rationale
2. Discussion in PR comments (async review acceptable)
3. Approval from project maintainer required
4. Constitution version bumped (semantic versioning applies)
5. Dependent templates updated in same commit

### Versioning Policy
- **MAJOR**: Removing or redefining core principles (e.g., dropping TDD requirement)
- **MINOR**: Adding new principles or materially expanding guidance
- **PATCH**: Clarifications, typo fixes, non-semantic improvements

### Compliance Review
- All PRs must verify alignment with Core Principles
- Template updates (plan, spec, tasks) checked when constitution changes
- Complexity deviations documented in `plan.md` Complexity Tracking section
- Use `AGENTS.md` and `CLAUDE.md` for agent-specific runtime guidance

### Conflict Resolution
When principles conflict (e.g., learning value vs. distribution simplicity), prioritize:
1. **Test-Driven Development** (principle III) - non-negotiable
2. **Go Learning First** (principle I) - primary objective
3. Other principles weighted by feature context

**Version**: 1.0.0 | **Ratified**: 2025-10-01 | **Last Amended**: 2025-10-01
