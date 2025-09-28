Ce plan cible un MVP pour aider les utilisateurs à initialiser leurs fichiers de configuration hexa via la CLI, sans mécanisme de versionnage ni migration automatique.

## TODO
- Définir le contenu des templates globaux (`~/.hexa.yml`) et locaux (`.hexa.local.yml`) utilisés pour l'initialisation.
- Documenter la procédure dans README.md : quand lancer `hexa config setup`, quelles sections renseigner manuellement, et précautions sur les secrets.
- Implémenter `hexa config setup` avec prompts interactifs : création conditionnelle des fichiers globaux et locaux, confirmation avant écrasement, message final récapitulatif.
- Gérer l'ajout de `.hexa.local.yml` au `.gitignore` du projet s'il n'est pas déjà listé.
- Ajouter des tests ciblés (unitaires ou commande) pour vérifier la génération des templates et le comportement idempotent.
