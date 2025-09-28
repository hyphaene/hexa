# Plan d'Implémentation - Système de Configuration hexa CLI

## Objectif Pédagogique Go

Implémenter un système de configuration hiérarchique pour apprendre :
- **Viper** : gestion configuration YAML + env vars + flags
- **go:embed** : intégration de templates de configuration
- **Cobra integration** : flags et subcommandes config
- **File operations** : lecture/écriture sécurisée, migration

---

## Phase 1 : Foundation & Structure

### Étape 1.1 : Package configuration
```go
// internal/config/config.go
```
**Concepts Go à apprendre :**
- Package organization (internal/)
- Struct tags pour Viper (`mapstructure`, `yaml`)
- Zero values et validation

**À implémenter :**
- Struct `Config` avec tags appropriés
- Validation des champs obligatoires
- Méthodes de sérialisation/désérialisation

### Étape 1.2 : Template par défaut
```go
// internal/config/template.go
```
**Concepts Go à apprendre :**
- `//go:embed` directive
- `text/template` package
- Variables d'environnement de build

**À implémenter :**
- Template YAML embarqué
- Fonction de génération avec placeholders
- Version tracking du template

---

## Phase 2 : Hierarchical Loading

### Étape 2.1 : Configuration Manager
```go
// internal/config/manager.go
```
**Concepts Go à apprendre :**
- Viper configuration precedence
- File system operations (`os`, `filepath`)
- Error wrapping avec `fmt.Errorf`

**À implémenter :**
- Fonction `LoadConfig()` avec ordre de priorité :
  1. Flags CLI (highest)
  2. Environment variables
  3. Project config (walking up directories)
  4. User global config
  5. Default template (lowest)

### Étape 2.2 : Directory Walking
```go
// internal/config/discovery.go
```
**Concepts Go à apprendre :**
- `filepath.Walk` ou custom walking
- `os.UserHomeDir()`
- Path manipulation cross-platform

**À implémenter :**
- `FindProjectConfig()` : remonte jusqu'à $HOME
- `GetUserConfigPath()` : config globale utilisateur
- Detection de `.hexa/` directories

---

## Phase 3 : CLI Integration

### Étape 3.1 : Subcommandes config
```go
// cmd/config.go
```
**Concepts Go à apprendre :**
- Cobra subcommands pattern
- Flag binding avec Viper
- Command organization

**À implémenter :**
```bash
hexa config init [--force] [--global]
hexa config show [--merged] [--scope=global|project]
hexa config validate
hexa config edit [--global]
hexa config upgrade [--dry-run]
```

### Étape 3.2 : Environment Variable Mapping
```go
// internal/config/env.go
```
**Concepts Go à apprendre :**
- Viper automatic env vars (`AutomaticEnv`, `SetEnvPrefix`)
- String manipulation et reflection
- Key transformation (dot notation → HEXA_KEY_SUBKEY)

**À implémenter :**
- Mapping automatique config → env vars
- Documentation des variables disponibles

---

## Phase 4 : Lifecycle Management

### Étape 4.1 : Safe File Operations
```go
// internal/config/writer.go
```
**Concepts Go à apprendre :**
- Atomic file writes (temp file + rename)
- File permissions handling
- Backup strategies

**À implémenter :**
- `WriteConfig()` avec atomic operations
- Backup avant modification
- Recovery en cas d'échec

### Étape 4.2 : Migration System
```go
// internal/config/migration.go
```
**Concepts Go à apprendre :**
- Version comparison
- YAML parsing et modification
- Interface design pour migrations

**À implémenter :**
- Detection de version obsolète
- Migration step-by-step
- Préservation des customizations utilisateur

---

## Phase 5 : Integration & Polish

### Étape 5.1 : Startup Integration
```go
// cmd/root.go - modification
```
**Concepts Go à apprendre :**
- Cobra PreRun hooks
- Graceful degradation
- Performance considerations

**À implémenter :**
- Chargement config au démarrage
- Fallback si config indisponible
- Warning messages informatifs

### Étape 5.2 : Shell Completion
```go
// cmd/completion.go - enhancement
```
**Concepts Go à apprendre :**
- Cobra completion system
- Dynamic completion (config keys)
- Shell integration

**À implémenter :**
- Completion pour config keys
- Completion pour scopes disponibles
- Validation en temps réel

---

## Ordre d'Implémentation Recommandé

1. **Start Simple** : Phase 1 (struct + template)
2. **Build Core** : Phase 2 (loading logic)
3. **Add CLI** : Phase 3 (user interface)
4. **Make Robust** : Phase 4 (lifecycle)
5. **Polish** : Phase 5 (integration)

## Points d'Attention Go

- **Error handling** : always wrap errors with context
- **Testing** : table-driven tests pour configurations
- **Documentation** : godoc comments pour API publique
- **Modules** : proper import paths depuis le module root

## Validation à Chaque Étape

```bash
go test ./internal/config/...
go run . config show --merged
hexa config validate
```

---

*Ce plan respecte l'approche pédagogique : concepts Go concrets + application pratique + progression logique.*