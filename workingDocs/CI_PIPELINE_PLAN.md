# Plan

Ce plan vise à enrichir la CI pour couvrir lint, tests, build et contrôles de qualité/sécurité sur chaque pull request active, afin de prévenir les régressions avant fusion.

## TODO

- Configurer un workflow GitHub Actions déclenché sur `pull_request` (opened, synchronize, reopened) et `push` sur les branches protégées.
- Ajouter un job `lint` : installer Go toolchain, exécuter `go fmt` (vérification), `go vet` et, si disponible, `golangci-lint` avec configuration stricte.
- Ajouter un job `test` : exécuter `go test ./...` avec couverture et uploader le rapport (utiliser Actions cache pour accélérer `go build` et modules).
- Ajouter un job `build` : compiler le binaire (`go build ./...`), valider que les artefacts se compilent sur Linux et macOS (matrices `GOOS`/`GOARCH` minimales).
- Intégrer des checks qualité supplémentaires optionnels :
  - Analyse statique de sécurité (gosec ou `govulncheck`).
  - Analyse de dépendances (`govulncheck -json`, `go list -m -u all`).
  - Lint Markdown/Docs (`markdownlint`) si la doc change.
- Configurer les règles de protection de branche pour exiger ces jobs avant merge et ajouter badges de statut dans le README.
