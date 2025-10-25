# Guide de Testing - Hexa

## Scripts npm disponibles

Le projet utilise npm comme task runner pour lancer les commandes Go facilement.

### Tests

```bash
# Tests basiques - quick feedback
npm test
# Équivalent: go test ./...

# Tests avec output détaillé
npm run test:verbose
# Équivalent: go test ./... -v

# Tests avec race detector (détecte race conditions)
npm run test:race
# Équivalent: go test ./... -race

# Tests complets (verbose + race + coverage)
npm run test:all
# Équivalent: go test ./... -v -race -coverprofile=coverage.out

# Coverage avec rapport HTML
npm run test:cover
# Génère coverage.out et coverage.html
# Ouvre coverage.html dans ton navigateur pour voir le détail
```

### Linting

```bash
# Lance golangci-lint
npm run check
# Équivalent: golangci-lint run
```

### Build

```bash
# Build le binaire hexa
npm run build
# Équivalent: go build -o hexa
```

## Workflow CI

Le workflow `.github/workflows/check.yml` s'exécute sur :

- ✅ Push sur `main`
- ✅ Pull requests (opened, synchronize, reopened)
- ✅ Déclenchement manuel (workflow_dispatch)

### Jobs

**1. lint** (bloquant)

- `gofmt -l .` - Vérifie le formatting
- `go vet ./...` - Analyse statique Go
- `golangci-lint` - Linter complet avec timeout 3min

**2. test** (bloquant, needs: lint)

- `go test ./... -v -coverprofile=coverage.out -covermode=atomic` - Tests avec coverage
- `go test ./... -race` - Tests avec race detector
- Vérification coverage ≥ 5% (à augmenter progressivement)
- Upload coverage artifact

**3. build** (bloquant, needs: test)

- Build multi-platform: linux/amd64, darwin/amd64
- Matrice de build pour tester les cross-compilations

**4. required-checks** (JOB GATE - required status check)

- `if: always()` - S'exécute même si d'autres jobs échouent
- `needs: [lint, test, build]` - Dépend de tous les jobs
- Vérifie explicitement que chaque job a réussi
- **C'est ce job qui doit être marqué "required" dans les branch protection rules**

## Configuration Branch Protection

Pour bloquer les merges si les tests échouent :

1. Va sur https://github.com/hyphaene/hexa/settings/branches
2. Add branch protection rule pour `main`
3. Coche "Require status checks to pass before merging"
4. Ajoute **uniquement** `required-checks` dans les status checks requis

Voir `.github/BRANCH_PROTECTION.md` pour plus de détails.

## Développement local

### Workflow recommandé

```bash
# 1. Fait tes modifications
vim internal/jira/tickets.go

# 2. Lance les tests rapides
npm test

# 3. Si ça passe, lance le linter
npm run check

# 4. Avant de commit, lance tout
npm run test:all

# 5. Si tout passe, commit et push
git add .
git commit -m "feat: amélioration de FetchSprintTickets"
git push
```

### Debug test failures

```bash
# Run un test spécifique
go test ./internal/jira/... -run TestFilterByStatus

# Run avec verbose pour voir les logs
go test ./internal/jira/... -v -run TestFilterByStatus

# Run un subtest spécifique
go test ./internal/jira/... -run TestFilterByStatus/filter_finds_matching

# Run avec race detector pour détecter race conditions
go test ./internal/jira/... -race -run TestFilterByStatus
```

### Voir la coverage

```bash
# Génère le rapport
npm run test:cover

# Ouvre coverage.html dans ton navigateur
open coverage.html  # macOS
xdg-open coverage.html  # Linux
```

Le rapport HTML montre ligne par ligne quel code est couvert (vert) ou non (rouge).

## Écrire des tests

### Structure d'un test Go

```go
func TestFunctionName(t *testing.T) {
    // 1. Setup
    input := "test data"

    // 2. Execute
    result, err := FunctionUnderTest(input)

    // 3. Assert
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    if result != expected {
        t.Errorf("got %v, want %v", result, expected)
    }
}
```

### Table-driven tests (recommandé)

```go
func TestSomething(t *testing.T) {
    tests := []struct {
        name  string
        input string
        want  int
    }{
        {name: "empty string", input: "", want: 0},
        {name: "single word", input: "hello", want: 1},
        {name: "multiple words", input: "hello world", want: 2},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got := CountWords(tt.input)
            if got != tt.want {
                t.Errorf("got %d, want %d", got, tt.want)
            }
        })
    }
}
```

Voir `internal/jira/tickets_test.go` pour des exemples concrets.

### Mock HTTP avec httptest

```go
server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(mockResponse)
}))
defer server.Close()

// Utilise server.URL dans ton test
client := &http.Client{}
resp, err := client.Get(server.URL + "/api/endpoint")
```

### Helper functions

```go
func setupTest(t *testing.T) func() {
    t.Helper()  // Les erreurs pointent vers l'appelant

    // Setup code
    viper.Set("config.key", "test-value")

    // Retourne cleanup function
    return func() {
        // Cleanup code
        viper.Reset()
    }
}

func TestSomething(t *testing.T) {
    cleanup := setupTest(t)
    defer cleanup()  // S'exécute toujours à la fin

    // Test code
}
```

## Coverage targets

**Actuel:** 8.8% global (27.5% sur internal/jira)

**Progression cible:**

- 🎯 **5%** - Baseline actuel (CI vérifie minimum)
- 🎯 **10%** - Court terme (ajouter tests cache, config)
- 🎯 **20%** - Moyen terme (ajouter tests cmd)
- 🎯 **50%** - Long terme (coverage solide)

**Note:** 100% de coverage n'est pas l'objectif. Focus sur le code critique et les edge cases.

## Race detector

Le race detector de Go détecte les **race conditions** (accès concurrent non synchronisé à la même mémoire).

```bash
# Lance avec race detector
npm run test:race

# Si race détectée, output détaillé:
# WARNING: DATA RACE
# Write at 0x... by goroutine X:
#   file.go:123
# Previous read at 0x... by goroutine Y:
#   file.go:456
```

**Important:** Le race detector a un overhead ~10x, c'est pourquoi on le run séparément en CI.

## Troubleshooting

### Tests passent localement mais échouent en CI

1. **Différence de timing:** Les tests concurrent peuvent être sensibles au timing
2. **Race conditions:** Lance `npm run test:race` localement
3. **Variables d'environnement:** Vérifie que CI a les mêmes configs
4. **Dépendances:** Vérifie go.mod/go.sum sont commités

### go test prend trop de temps

```bash
# Run tests en parallèle (défaut, mais peut être forcé)
go test ./... -parallel 4

# Timeout plus court pour détecter les tests qui hang
go test ./... -timeout 30s

# Skip les tests lents si marqués avec -short
go test ./... -short
```

### "no test files" pour certains packages

C'est normal. Go dit juste qu'il n'y a pas de fichiers `*_test.go` dans ces packages.

Pour l'instant, seul `internal/jira` a des tests. Au fur et à mesure, ajoute des tests pour :

- `internal/cache` - Tests de read/write cache
- `internal/config` - Tests de config loading
- `cmd/*` - Tests des commandes CLI (plus difficile)

## Resources

- [Go Testing Package](https://pkg.go.dev/testing)
- [Table Driven Tests](https://go.dev/wiki/TableDrivenTests)
- [httptest Package](https://pkg.go.dev/net/http/httptest)
- [Go Race Detector](https://go.dev/doc/articles/race_detector)
- Voir aussi: `internal/jira/TESTING.md` pour patterns détaillés
