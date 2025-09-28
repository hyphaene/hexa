# Proposition CLI Go : hexa

## Vue d'ensemble

Transformation des 22+ scripts bash en un CLI Go unifié et distributable, organisé en domaines fonctionnels avec une hiérarchie claire.

**Approche didactique** : Ce plan sert de guide pour développer le CLI ensemble, avec un support technique sans implémentation directe.

```bash
hexa [DOMAIN] [COMMAND] [SUBCOMMAND] [FLAGS]
hw [DOMAIN] [COMMAND] [SUBCOMMAND] [FLAGS]    # Alias court
```

## Naming & Distribution

### Nom validé : `hexa` (alias `hw`)

- **Nom principal** : `hexa` (lien avec Hexactitude, unique, mémorable)
- **Alias court** : `hw` (usage quotidien, 2 caractères)
- **Setup** : Symlink automatique `hw` → `hexa`

## Stratégie de packaging validée

### Phase 1 : Wrapper avec scripts embarqués (embed)

- **Single binary** : Scripts bash intégrés via `//go:embed`
- **Distribution simplifiée** : Un seul fichier exécutable
- **Migration progressive** : Remplacer script par script par du Go pur selon priorité

## Architecture proposée

### 1. Domaine JIRA (`jira`)

#### 1.1 Sprint Management (`hexa jira sprint`)

```bash
# Aperçu du sprint actuel
hexa jira sprint overview [--user me|all|unassigned] [--format table|json]

# Sprint actif
hexa jira sprint current [--format id|name|full]

# Tous les tickets du sprint
hexa jira sprint tickets [--status todo|progress|done|all]

# ID du sprint actuel
hexa jira sprint id
```

**Mapping des scripts** :

- `jira_overview.sh` → `hexa jira sprint overview`
- `get_current_sprint.sh` → `hexa jira sprint current`
- `fetch_all_sprint_tickets.sh` → `hexa jira sprint tickets`
- `get_sprint_id.sh` → `hexa jira sprint id`

#### 1.2 Ticket Operations (`hexa jira ticket`)

```bash
# Gestion des tickets
hexa jira ticket get TEAM-12345 [--format json|summary]
hexa jira ticket move TEAM-12345 --status "In Progress"
hexa jira ticket comment TEAM-12345 "Votre commentaire"

# Pièces jointes
hexa jira ticket attachments TEAM-12345 [--list|--download]
hexa jira ticket attachments TEAM-12345 --download [--output-dir ./downloads]

# Catégorisation intelligente
hexa jira ticket categorize TEAM-12345
```

**Mapping des scripts** :

- `jira_fetch.sh` → `hexa jira ticket get`
- `move_ticket.sh` → `hexa jira ticket move`
- `add_comment.sh` → `hexa jira ticket comment`
- `view_attachments.sh` → `hexa jira ticket attachments --list`
- `download_attachments.sh` → `hexa jira ticket attachments --download`
- `jira_categorize.sh` → `hexa jira ticket categorize`

#### 1.3 Release Management (`hexa jira release`)

```bash
# Gestion des release notes
hexa jira release notes [DATA_FILE] [--command report|update-all|update-mine]
hexa jira release notes [DATA_FILE] --two-steps

# Validation des versions
hexa jira release verify [--version VERSION]

# Récupération des tickets par version
hexa jira release tickets --version VERSION [--format json|table]
```

**Mapping des scripts** :

- `release_notes_manager.sh` → `hexa jira release notes`
- `release_notes_two_steps.sh` → `hexa jira release notes --two-steps`
- `verify_version_completeness.sh` → `hexa jira release verify`
- `fetch_version_tickets.sh` → `hexa jira release tickets`

#### 1.4 Development Workflow (`hexa jira dev`)

```bash
# Setup de tâche (worktree + contexte)
hexa jira dev setup TEAM-12345 [--suffix SUFFIX] [--type feat|fix|hotfix]

# Synchronisation centrale
hexa jira dev sync --from-central [TICKET_KEY]
hexa jira dev sync --to-central
```

**Mapping des scripts** :

- `setup_task.sh` → `hexa jira dev setup`
- `sync_from_central.sh` → `hexa jira dev sync --from-central`
- `sync_to_central.sh` → `hexa jira dev sync --to-central`

#### 1.5 Testing & Utilities (`hexa jira test`)

```bash
# Tests rapides
hexa jira test comment [TICKET_KEY]
hexa jira test tag [TICKET_KEY]
hexa jira test slack [MESSAGE]
```

**Mapping des scripts** :

- `test_comment.sh` → `hexa jira test comment`
- `test_tag.sh` → `hexa jira test tag`
- `test_slack.sh` → `hexa jira test slack`

### 2. Domaine GIT (`git`)

#### 2.1 Repository Management (`hexa git repo`)

```bash
# Status multi-repos
hexa git repo status [--format table|json]

# Collection des repos
hexa git repo collect [--path PATH]

# Worktree management
hexa git worktree create BRANCH_NAME [--from BRANCH]
```

**Mapping des scripts** :

- `multi-repo-status.sh` → `hexa git repo status`
- `git/collect-git-repos.sh` → `hexa git repo collect`
- `create-worktree.sh` → `hexa git worktree create`

#### 2.2 Analysis (`hexa git analyze`)

```bash
# Analyse de cherry-pick
hexa git analyze cherry-pick [--commits COMMITS]

# Recherche raw
hexa git analyze search [PATTERN] [--raw]
```

**Mapping des scripts** :

- `analyse-cherry-pick.sh` → `hexa git analyze cherry-pick`
- `search-raw.sh` → `hexa git analyze search`

### 3. Domaine SETUP (`setup`)

#### 3.1 Environment Setup (`hexa setup env`)

```bash
# Configuration des symlinks
hexa setup env claude-md
hexa setup env alexandria
hexa setup env commands
hexa setup env agents
hexa setup env settings

# Setup complet
hexa setup env all
```

**Mapping des scripts** :

- `setup-claude-md-symlink.sh` → `hexa setup env claude-md`
- `setup-alexandria-symlink.sh` → `hexa setup env alexandria`
- `setup-commands-symlinks.sh` → `hexa setup env commands`
- `setup-agents-symlinks.sh` → `hexa setup env agents`
- `setup-settings-symlink.sh` → `hexa setup env settings`

### 4. Domaine AI (`ai`)

#### 4.1 Claude Code Surcouches (`hexa ai cc`)

```bash
# Surcouches sympas pour Claude Code
hexa ai cc pull-request
hexa ai cc workflow

```

#### 4.2 Trucs IA à creuser (`hexa ai explore`)

```bash
# Idées à développer
hexa ai code-analysis [FILES...]

hexa ai bug-hunter [--auto-fix]
hexa ai docs-gen [--intelligent]
```

## Structure du projet Go

```
hexa/
├── cmd/                        # Commandes Cobra
│   ├── root.go                 # Commande racine
│   ├── jira/                   # Domaine Jira
│   ├── git/                    # Domaine Git
│   └── setup/                  # Domaine Setup
├── scripts/                    # Scripts bash source (pour embed)
│   ├── jira/                   # Scripts Jira existants
│   ├── git/                    # Scripts Git existants
│   └── setup/                  # Scripts Setup existants
├── internal/
│   ├── executor/               # Exécution scripts embarqués
│   ├── config/                 # Configuration (Viper)
│   └── utils/                  # Utilitaires partagés
└── main.go
```

## Technologies validées

- **CLI Framework** : Cobra (commandes + flags + help)
- **Configuration** : Viper (YAML + env vars)
- **Embedding** : `//go:embed` (scripts dans binaire)
- **Execution** : os/exec (scripts temporaires)

## Configuration centralisée multi-niveaux

### Hiérarchie de configuration (spécifique override global)

Le système de configuration utilise **Viper avec MergeConfig()** pour un support multi-niveaux natif :

```
1. ~/.hexa.yml              (global user)
2. ~/Code/project/.hexa.yml (project-specific) ← override
3. ./hexa.yml               (current directory) ← override
4. ENV variables            (runtime) ← override
5. Command line flags       (runtime) ← override
```

**Ordre de précédence** : Plus spécifique = plus prioritaire.

### Implémentation Viper

```go
// Lecture et merge automatique des configs
func loadConfig() error {
    viper.SetConfigType("yaml")

    // 1. Base : global config (~/.hexa.yml)
    if homeConfig, err := os.ReadFile(filepath.Join(home, ".hexa.yml")); err == nil {
        viper.ReadConfig(bytes.NewBuffer(homeConfig))
    }

    // 2. Override : project config (remontée vers ~/)
    projectPath := findProjectConfig() // cherche .hexa.yml jusqu'à ~/
    if projectConfig, err := os.ReadFile(projectPath); err == nil {
        viper.MergeConfig(bytes.NewBuffer(projectConfig)) // ← Merge avec précédence
    }

    // 3. Override : current dir (./hexa.yml)
    if currentConfig, err := os.ReadFile("./hexa.yml"); err == nil {
        viper.MergeConfig(bytes.NewBuffer(currentConfig))
    }

    // 4. ENV variables automatiques (HEXA_JIRA_URL, etc.)
    viper.AutomaticEnv()
    viper.SetEnvPrefix("HEXA")
    viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

    return nil
}
```

### Structure YAML supportée

```yaml
# ~/.hexa.yml (global)
jira:
  url: "https://jira.company.com"
  token: "${JIRA_PAT}"
  default_project: "GLOBAL"

git:
  default_branch: "main"
  worktree_base: "~/worktrees"

user:
  me: "maximilien.garenne@company.com"
```

```yaml
# ~/Code/specific-project/.hexa.yml (project override)
jira:
  default_project: "PROJ"  # Override global
  # url héritée de global

git:
  default_branch: "develop"  # Override global
  worktree_base: "./local-worktrees"  # Override global

# user héritée de global

# Config spécifique au projet
database:
  host: "localhost"
  port: 5432
```

### Recherche de configuration projet

```go
// Remonte depuis PWD jusqu'à home pour trouver .hexa.yml
func findProjectConfig() string {
    dir, _ := os.Getwd()
    home, _ := os.UserHomeDir()

    for dir != home && dir != "/" {
        configPath := filepath.Join(dir, ".hexa.yml")
        if _, err := os.Stat(configPath); err == nil {
            return configPath
        }
        dir = filepath.Dir(dir)
    }
    return ""
}
```

### Aliases et workflows

```yaml
aliases:
  # Shortcuts quotidiens
  overview: "jira sprint overview --user me"
  setup: "jira dev setup"
  comment: "jira ticket comment"
  move: "jira ticket move"

  # Workflows composés
  daily: ["jira sprint overview --user me", "git repo status"]
  wip: "jira ticket move {ticket} --status 'In Progress' && jira ticket comment {ticket} 'Work in progress'"
```

### Support des variables d'environnement

```bash
# Override à l'exécution
HEXA_JIRA_URL="https://jira.dev.com" hexa jira sprint overview
export HEXA_GIT_DEFAULT_BRANCH="feature-branch"
hexa git repo status
```

## Smart Completion & Help (Phase 1)

### Help automatique à tous les niveaux

```bash
hexa --help                         # Root help
hexa jira --help                    # Domain help
hexa jira sprint --help             # Command group help
hexa jira sprint overview --help    # Specific command help
```

### Autocomplétion intelligente

```bash
# Setup
hexa completion bash > /etc/bash_completion.d/hexa
hexa completion zsh > ~/.zsh/completions/_hexa

# Completion dynamique (API calls)
hw ticket comment <TAB>             # → Tes tickets actifs (TEAM-12345, TEAM-12346...)
hw ticket move TEAM-12345 --status <TAB>  # → Transitions disponibles ("In Progress", "Done"...)
hw setup env <TAB>                  # → claude-md, alexandria, commands, agents, all
```

### Cache intelligent

```bash
# Cache 5min pour éviter appels répétés
~/.hexa/cache/
├── tickets_active.json     # Tes tickets actifs
├── transition_TEAM.json    # Transitions possibles
└── sprints_current.json    # Sprint actuel
```

### Usage aliases

```bash
# Aliases simples
hw overview                    # → hw jira sprint overview --user me
hw comment TEAM-12345 "msg"     # → hw jira ticket comment TEAM-12345 "msg"
hw setup TEAM-12345             # → hw jira dev setup TEAM-12345

# Workflows composés
hw daily                       # Exécute overview + git status
hw wip TEAM-12345              # → move + comment automatiques
```

## Exemples d'usage

### Workflow quotidien typique

```bash
# 1. Vue d'ensemble du sprint (ou alias)
hexa jira sprint overview --user me
hw overview                            # Alias équivalent

# 2. Setup d'une nouvelle tâche
hexa jira dev setup TEAM-12345 --type feat
hw setup TEAM-12345                     # Alias équivalent

# 3. Ajouter un commentaire
hexa jira ticket comment TEAM-12345 "Work in progress on authentication"
hw comment TEAM-12345 "Work in progress on authentication"  # Alias

# 4. Déplacer le ticket (avec completion)
hexa jira ticket move TEAM-12345 --status "In Progress"
hw wip TEAM-12345                       # Alias avec move + comment automatique

# 5. Workflow complet quotidien
hw daily                               # overview + git status
```

### Gestion des releases

```bash
# 1. Validation de la complétude
hexa jira release verify --version "1.2.3"

# 2. Génération des notes
hexa jira release notes ./release_data.json --command report

# 3. Mise à jour des tickets
hexa jira release notes ./release_data.json --command update-all
```

## Avantages vs scripts actuels

### ✅ **Ce qu'on gagne**

1. **Distribution simplifiée** : Un seul binaire vs 22+ scripts
2. **Autocomplétion** : Cobra fournit l'autocomplétion automatique
3. **Help intégré** : `--help` sur chaque commande
4. **Validation** : Arguments et flags validés avant exécution
5. **Configuration centralisée** : Plus de duplication de config
6. **Cross-platform** : Fonctionne sur macOS, Linux, Windows
7. **Performance** : Go vs Bash pour les opérations complexes
8. **Maintenabilité** : Code structuré vs scripts éparpillés

### 🔄 **Migration progressive**

1. **Phase 1** : CLI qui wrappe les scripts existants
2. **Phase 2** : Réécriture progressive en Go pur
3. **Phase 3** : Optimisations et nouvelles features

## Next Steps

1. **Prototype minimal** avec Cobra + 2-3 commandes principales
2. **Configuration Viper** pour remplacer le sourcing de `.env`
3. **Migration progressive** script par script
4. **Tests** pour s'assurer de la parité fonctionnelle
