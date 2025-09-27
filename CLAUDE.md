# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Context & Objectif

Ce projet est une **base d'apprentissage Go** pour monter en compétence sur l'écosystème, la syntaxe et les primitives de Go.

**Approche** : Expert Go, pédagogue mais direct, sans prendre de gants. Ne pas générer de code sauf demande explicite.

**Vision** : CLI Go unifié "hexa" (alias "hw") pour remplacer 22+ scripts bash par un single binary distributable avec Homebrew.

## Architecture

### Structure actuelle (minimale)

```
hexa/
├── main.go                 # Point d'entrée → cmd.Execute()
├── cmd/                    # Commandes Cobra
│   ├── root.go            # Commande racine (placeholder)
│   └── version.go         # Version command
├── scripts/               # Scripts bash à embarquer (vide)
├── internal/              # Packages internes (vide)
└── .goreleaser.yaml       # Configuration GoReleaser
```

### Architecture cible (selon CLI_GO_PROPOSAL.md)

- **Domaines** : JIRA, GIT, SETUP, AI
- **Framework** : Cobra (commandes + flags + help)
- **Configuration** : Viper (YAML + env vars)
- **Embedding** : `//go:embed` (scripts dans binaire)
- **Distribution** : Homebrew avec symlink `hw` → `hexa`

### Stack technique

- **Go** : 1.24.4
- **CLI Framework** : Cobra v1.10.1
- **Distribution** : GoReleaser + Homebrew tap

## Commandes de développement

### Build & Run

```bash
# Build simple
go build

# Build avec nom spécifique
go build -o hexa

# Run direct
go run main.go [args]

# Run du binaire
./hexa --help
./hexa version
```

### Tests

```bash
# Run tous les tests
go test ./...

# Tests avec verbose
go test -v ./...

# Tests avec couverture
go test -cover ./...
```

### GoReleaser

```bash
# Test build local (snapshot)
goreleaser release --snapshot --rm-dist

# Release (nécessite tag git)
goreleaser release --rm-dist
```

### Vérifications Go

```bash
# Modules management
go mod tidy
go mod download

# Formatting
go fmt ./...

# Linting (si golangci-lint installé)
golangci-lint run

# Vet
go vet ./...
```

## Concepts Go clés pour l'apprentissage

### Packages et modules

- `go.mod` : Définition du module et dépendances
- Import paths : `github.com/hyphaene/hexa/cmd`
- Package main : Point d'entrée avec `func main()`

### Cobra CLI patterns

- **Root command** : `cmd/root.go` avec `&cobra.Command{}`
- **Subcommands** : `rootCmd.AddCommand()` dans `init()`
- **Flags** : Persistent vs local flags
- **Help** : Automatique avec descriptions

### Embedding avec `//go:embed`

```go
//go:embed scripts/*.sh
var scriptsFS embed.FS
```

### Configuration avec Viper

- YAML config files
- Environment variables
- Flags precedence

## Versioning & Release

### Stratégie Semantic Versioning

Le projet utilise **semantic versioning automatique** basé sur conventional commits :

#### Types de commits et impacts sur la version

| Type de commit | Impact version | Exemple |
|---|---|---|
| `feat:` | **Minor** (v0.1.0 → v0.2.0) | `feat: add jira domain commands` |
| `fix:` | **Patch** (v0.1.0 → v0.1.1) | `fix: resolve authentication bug` |
| `perf:`, `refactor:` | **Patch** (v0.1.0 → v0.1.1) | `perf: optimize git operations` |
| `BREAKING CHANGE:` | **Major** (v0.1.0 → v1.0.0) | Avec footer `BREAKING CHANGE:` |
| `docs:`, `chore:`, `ci:` | **Skip** (pas de release) | `docs: update README` |

#### Workflow automatique

```bash
# À chaque push sur main :
1. Analyse des commits depuis le dernier tag
2. Calcul automatique de la prochaine version (semver)
3. Création du tag Git (si nécessaire)
4. Release GoReleaser + update Homebrew (si nouvelle version)
5. Skip complet si commits non-releaseables (docs, chore, etc.)
```

#### Exemples concrets

```bash
# Scénario 1: Nouvelles features
git commit -m "feat: add git worktree commands"
git commit -m "feat: add slack notifications"
git push origin main
# → v0.1.9 → v0.2.0 (minor bump)

# Scénario 2: Bug fixes
git commit -m "fix: resolve jira auth timeout"
git push origin main
# → v0.2.0 → v0.2.1 (patch bump)

# Scénario 3: Documentation
git commit -m "docs: improve installation guide"
git push origin main
# → v0.2.1 → v0.2.1 (aucune release)

# Scénario 4: Breaking change
git commit -m "feat!: restructure CLI commands

BREAKING CHANGE: Command syntax changed from 'hexa jira-sprint' to 'hexa jira sprint'"
git push origin main
# → v0.2.1 → v1.0.0 (major bump)
```

## Distribution & Installation

### Homebrew (production)

```bash
# Add tap first
brew tap hyphaene/hexa

# Install hexa
brew install hexa

# Usage avec alias
hw --help    # Équivalent à hexa --help
```

#### Homebrew Tap Repository

- **Path local** : `~/Code/homebrew-hexa` (repo séparé)
- **GitHub** : [github.com/hyphaene/homebrew-hexa](https://github.com/hyphaene/homebrew-hexa)
- **Auto-update** : GoReleaser push automatiquement les nouvelles versions

### Development

```bash
# Build local
go build -o hexa

# Test installation
./hexa --help
./hexa version
```

## Développement progressif recommandé

1. **Phase 1** : Cobra structure + commandes de base
2. **Phase 2** : Embedding scripts bash existants
3. **Phase 3** : Réécriture progressive en Go pur
4. **Phase 4** : Configuration Viper + optimisations

## Règles importantes

- **Ne jamais** générer de code sans demande explicite
- **Toujours** préserver les conventions Go existantes
- **Focus** sur l'apprentissage des primitives Go
- **Approche** pédagogique directe, sans couvrir
