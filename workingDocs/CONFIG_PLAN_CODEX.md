# hexa Configuration Implementation Plan (Codex Edition)

## 1. Préparation de l'atelier
1. Relire `CONFIG_REQUIREMENTS.md` et lister les structures Go impliquées (configuration globale, overrides, options CLI) pour clarifier les responsabilités.
2. Identifier les modules existants de la CLI (packages Cobra, répertoires `cmd/`, `internal/` ou `pkg/`) afin de planifier où placer les nouvelles fonctions (`config`, helpers de fichiers, migration, etc.).
3. Définir une structure de tests (Go test + éventuels scénarios end-to-end) qui accompagnera le développement à chaque étape.

## 2. Chargeur de configuration multi-niveaux
1. Implémenter un helper `type ConfigLoader` (ou équivalent) qui encapsule Viper et accepte une liste ordonnée de sources.
2. Écrire `loadGlobalConfig()` qui lit `~/.hexa.yml`. Injecter la valeur de `home` via `os.UserHomeDir()` et gérer les erreurs (log + retours clairs, pas de panique).
3. Implémenter `findProjectConfig(basePath string) (string, error)` :
   - Parcourir les ancêtres du répertoire courant jusqu'à `home`.
   - Gérer liens symboliques (normalisation `filepath.EvalSymlinks`).
   - Introduire un cache simple pour éviter le walk à chaque commande (ex : memoisation dans le package).
4. Ajouter `loadConfigFrom(path string)` qui retourne une structure vide si le chemin est vide ou manquant.
5. Brancher les trois couches (`global`, `project`, `local`) dans l'ordre défini, avec un merge contrôlé :
   - Pour les slices, définir une stratégie personnalisée (ex : copier puis fusionner) car Viper ne merge pas les listes nativement.
6. Configurer les overrides runtime :
   - Poser `viper.SetEnvPrefix("HEXA")` et `SetEnvKeyReplacer` avant `AutomaticEnv()`.
   - Exposer un point d'extension pour alimenter la config depuis les flags Cobra.
7. Produire des tests unitaires qui valident la hiérarchie et les cas d'échec (fichier inexistant, permission refusée, variables d'environnement, overrides flag).

## 3. Cycle de vie de la configuration par défaut
1. Créer `templates/hexa-template.yml` avec placeholders explicites (`${USER_EMAIL}`, `${HOME}`, etc.) et commenter les sections sensibles.
2. Introduire `type ConfigVersion struct` qui stocke version courante, info de migration et éventuellement checksum du template.
3. Écrire `func ensureConfigUpToDate(ctx context.Context) error` qui :
   - Charge l'état du fichier global s'il existe.
   - Compare la version stockée avec `CurrentConfigVersion`.
   - Déclenche `runMigrations(oldVersion, newVersion)` si nécessaire.
4. Implémenter `createOrUpgradeConfig` :
   - Remplacer les placeholders via une structure dédiée (ex : `TemplateContext`).
   - Écrire dans un fichier temporaire puis faire `os.Rename` pour garantir l'atomicité.
   - Sauvegarder l'ancienne version (ex : `~/.hexa.yml.bak`) en cas d'échec.
5. Construire un système de migrations incrémentales (tableau de fonctions) pour gérer l'évolution du schéma sans écraser les personnalisations.
6. Couvrir par des tests :
   - Première installation.
   - Upgrade avec modifications utilisateur.
   - Échecs disque pleine / répertoire read-only (mock FS ou `io/fs` factice).

## 4. Intégration distribution & Homebrew
1. Ajouter une commande `hexa config doctor` qui vérifie la configuration sans la modifier (utile pour CI et post-install).
2. Mettre à jour la formule Homebrew dans GoReleaser :
   - Ajouter le hook `post_install` qui exécute `hexa config init --if-missing` ou `doctor`.
   - Documenter les limitations sandbox et prévoir un fallback silencieux (log vers stdout, code retour 0 si échec permissif).
3. Dans la CLI, rendre `ensureConfigUpToDate` tolérant aux erreurs :
   - Si l’écriture échoue, afficher un warning mais continuer la commande.
   - Fournir un code d’erreur spécifique pour permettre aux scripts d’agir.
4. Préparer un guide dans la doc (README / site) sur l’intégration Homebrew et les étapes post-install.

## 5. Commandes CLI
1. Créer la commande racine `hexa config` puis les sous-commandes :
   - `init` : options `--force`, `--if-missing`, mode interactif facultatif.
   - `show` : flags `--global`, `--project`, `--effective`.
   - `validate` : exécute les checks structurels et affiche les erreurs.
   - `edit` : ouvre l’éditeur avec support `--scope`.
   - `doctor` : diagnostic complet (cf. section précédente).
2. Centraliser la logique de sortie (fmt/JSON). Prévoir `--format json` pour scripts si pertinent.
3. Exposer l’autocomplétion Cobra :
   - Implémenter `hexa completion [bash|zsh|fish|powershell]`.
   - Ajouter doc d’installation rapide.
4. Brancher les hooks `cobra.OnInitialize(initConfig)` en s’assurant qu’ils sont idempotents et testés.
5. Tester chaque commande via tests Cobra (`cmd.ExecuteC`) + tests end-to-end (scenario: `init`, modification manuelle, `show`, `validate`).

## 6. Expérience utilisateur & validation
1. Écrire un validateur dédié (structure ou package) qui :
   - Charge la config, vérifie le schéma (types, valeurs obligatoires, valeurs supportées).
   - Retourne une liste d’avertissements/erreurs compréhensibles.
2. Créer une table de correspondance `config key → env var` dans la doc et l’afficher dans `hexa config show --help`.
3. Mettre en place la génération d’un changelog :
   - Automatiser via GoReleaser ou script pour maintenir `CHANGELOG.md`.
   - Lors d’un upgrade, afficher les sections pertinentes (diff version installée → nouvelle) et pointer vers les actions à effectuer (migrations config, nouvelles commandes).
4. Préparer des messages utilisateur cohérents (stdout/stderr), localiser si nécessaire, et éviter les emoji en contexte non interactif (CI, logs Homebrew).

## 7. Tests, CI et validation finale
1. Étendre la configuration CI pour exécuter :
   - Tests unitaires Go.
   - Tests d’intégration CLI (via `go test` ou script `./scripts/test-config.sh`).
   - Vérification de la génération autocomplétion (lint des scripts générés).
2. Ajouter des checks linters (golangci-lint, staticcheck) sur les nouveaux packages.
3. Mettre en place un job qui simule l’installation Homebrew en local (utiliser `brew install --build-from-source` dans un env temporaire).
4. Rédiger un guide de QA manuel (scénarios premier démarrage, upgrade, override projet, environnements read-only).
5. Avant release, mettre à jour la documentation et vérifier que `CONFIG_REQUIREMENTS.md` est entièrement satisfait.

---

> **Note développeur** : Chaque section est conçue pour être implémentée puis validée via tests avant de passer à la suivante. Utilise ce plan comme checklist : coche les étapes réalisées, ajoute des notes sur les décisions ou écarts, et enrichis-le au fil du développement.
