#!/usr/bin/env bash

set -euo pipefail

# 🏷️ GitHub Label Sync Script
# Fetches all labels from the current GitHub repository and stores them in repo.config.json

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

error() {
    echo -e "${RED}❌ $1${NC}" >&2
    exit 1
}

success() {
    echo -e "${GREEN}✅ $1${NC}"
}

info() {
    echo -e "${YELLOW}🔍 $1${NC}"
}

# Vérifier qu'on est dans un repo git
if ! git rev-parse --git-dir > /dev/null 2>&1; then
    error "Pas dans un repository git"
fi

# Vérifier que gh CLI est installé
if ! command -v gh &> /dev/null; then
    error "gh CLI n'est pas installé\n💡 Installez-le: https://cli.github.com"
fi

# Vérifier qu'on est authentifié
if ! gh auth status &> /dev/null; then
    error "Non authentifié avec GitHub CLI\n💡 Lancez: gh auth login"
fi

# Vérifier qu'il y a un remote GitHub
if ! gh repo view &> /dev/null; then
    error "Aucun remote GitHub détecté"
fi

# Récupérer la racine du repo git
repo_root=$(git rev-parse --show-toplevel)
config_file="$repo_root/repo.config.json"

info "Fetching labels from GitHub..."

# Fetch labels
labels_json=$(gh label list --json name,description,color --jq '.')

if [ $? -ne 0 ]; then
    error "Échec du fetch des labels GitHub"
fi

# Compter le nombre de labels
label_count=$(echo "$labels_json" | jq 'length')

# Générer le timestamp ISO 8601
synced_at=$(date -u +"%Y-%m-%dT%H:%M:%SZ")

# Construire le JSON final
final_json=$(jq -n \
    --argjson labels "$labels_json" \
    --arg synced_at "$synced_at" \
    '{labels: $labels, synced_at: $synced_at}')

# Écrire dans le fichier
echo "$final_json" > "$config_file"

if [ $? -ne 0 ]; then
    error "Échec de l'écriture dans $config_file"
fi

success "Labels synchronisés avec succès!"
echo ""
echo "📊 $label_count labels sauvegardés dans:"
echo "   $config_file"
echo ""
echo "🕐 Dernière synchronisation: $synced_at"
