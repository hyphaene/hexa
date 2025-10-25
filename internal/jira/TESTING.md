# Guide des Tests Go - internal/jira/tickets_test.go

Ce document explique les patterns de testing Go utilisés dans ce fichier de tests.

## Patterns Go idiomatiques utilisés

### 1. Table-Driven Tests

**Pattern le plus courant en Go** - Un slice de structs anonymes contenant cas de test.

```go
tests := []struct {
    name    string
    input   []Ticket
    want    int
}{
    {
        name: "filter finds matching status",
        input: []Ticket{...},
        want: 2,
    },
    // Plus de cas...
}

for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
        got := FilterByStatus(tt.tickets, tt.statusName)
        if len(got) != tt.want {
            t.Errorf("got %d, want %d", len(got), tt.want)
        }
    })
}
```

**Avantages:**

✅ Ajouter un cas de test = ajouter une entrée au slice
✅ Chaque test est isolé avec `t.Run`
✅ Facile à lire et maintenir

### 2. Subtests avec t.Run()

```go
t.Run(tt.name, func(t *testing.T) {
    // Test logic
})
```

**Ce que ça fait:**

- Chaque subtest s'exécute indépendamment
- Affichage hiérarchique: `TestFilterByStatus/filter_finds_matching_status`
- Peut run un seul subtest: `go test -run TestFilterByStatus/filter_finds_matching`

### 3. Helper Functions avec t.Helper()

```go
func setupViperForTest(t *testing.T, jiraURL, jiraToken string) func() {
    t.Helper()  // ✅ Les erreurs pointent vers l'appelant, pas cette fonction

    // Setup code
    viper.Set("jira.url", jiraURL)

    // Retourne cleanup function
    return func() {
        // Cleanup code
    }
}
```

**Usage:**

```go
func TestSomething(t *testing.T) {
    cleanup := setupViperForTest(t, "url", "token")
    defer cleanup()  // ✅ Cleanup automatique à la fin du test

    // Test code
}
```

**Pourquoi `defer`?**

- S'exécute TOUJOURS, même si test panic
- S'exécute à la fin de la fonction
- Multiple defers = stack (LIFO)

### 4. Fixture Functions

```go
func makeTicket(key string, status string, assigneeEmail string) Ticket {
    ticket := Ticket{
        Key: key,
        Fields: Fields{
            Status: Status{Name: status},
        },
    }

    if assigneeEmail != "" {
        ticket.Fields.Assignee = &Assignee{
            EmailAddress: assigneeEmail,
        }
    }

    return ticket
}
```

**Avantages:**

✅ Évite duplication de création d'objets test
✅ Centralise la logique de test data
✅ Facile à modifier si les structs changent

### 5. Mock HTTP Server avec httptest

```go
server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    // Vérifier la request
    if r.Method != "GET" {
        t.Errorf("expected GET, got %s", r.Method)
    }

    // Envoyer mock response
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(mockResponse)
}))
defer server.Close()

// Utiliser server.URL dans les tests
```

**Ce qui se passe:**

1. `httptest.NewServer` démarre un serveur HTTP local
2. `server.URL` donne l'URL (ex: `http://127.0.0.1:54321`)
3. Ton code fait des vraies requêtes HTTP vers ce server
4. Le handler vérifie les requests et envoie les responses mockées
5. `defer server.Close()` arrête le serveur après le test

**Pourquoi c'est puissant:**

✅ Pas de dépendance externe (pas de vrai Jira)
✅ Tests rapides et déterministes
✅ Peut tester erreurs réseau, timeouts, etc.

### 6. Error Testing Pattern

```go
got, err := FunctionUnderTest()

// Pattern 1: Vérifier si erreur attendue
if (err != nil) != tt.wantErr {
    t.Errorf("error = %v, wantErr %v", err, tt.wantErr)
    return  // ✅ Stop si erreur inattendue
}

// Pattern 2: Vérifier message d'erreur
if err != nil && err.Error() != tt.wantErrMsg {
    t.Errorf("error message = %q, want %q", err.Error(), tt.wantErrMsg)
}
```

## Structure d'un test Go

```go
func TestFunctionName(t *testing.T) {
    // 1. Setup
    input := makeFixture()
    cleanup := setupDependencies(t)
    defer cleanup()

    // 2. Execute
    got, err := FunctionUnderTest(input)

    // 3. Assert
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }

    if got != want {
        t.Errorf("got %v, want %v", got, want)
    }
}
```

## t.Errorf vs t.Fatalf vs t.Error

```go
t.Errorf("format", args)  // ✅ Log erreur, continue le test
t.Fatalf("format", args)  // ❌ Log erreur, STOP le test immédiatement
t.Error("message")        // Comme Errorf mais sans format
```

**Quand utiliser quoi:**

- `t.Errorf`: Vérifications multiples, tu veux voir toutes les failures
- `t.Fatalf`: Failure critique qui rend le reste du test inutile

```go
got, err := Parse(input)
if err != nil {
    t.Fatalf("Parse failed: %v", err)  // ✅ Stop, on peut pas continuer
}

if got.Name != want.Name {
    t.Errorf("Name = %q, want %q", got.Name, want.Name)  // ✅ Continue
}
if got.Age != want.Age {
    t.Errorf("Age = %d, want %d", got.Age, want.Age)  // ✅ Voit les 2 failures
}
```

## Running Tests

```bash
# Tous les tests
go test ./...

# Package spécifique
go test ./internal/jira/...

# Avec verbose output
go test ./internal/jira/... -v

# Un seul test
go test ./internal/jira/... -run TestFilterByStatus

# Un seul subtest
go test ./internal/jira/... -run TestFilterByStatus/filter_finds_matching

# Avec coverage
go test ./internal/jira/... -cover

# Coverage report HTML
go test ./internal/jira/... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

## Concepts Go à explorer dans les tests

### 1. Defer stack

```go
func example(t *testing.T) {
    defer fmt.Println("1")  // S'exécute 3ème
    defer fmt.Println("2")  // S'exécute 2ème
    defer fmt.Println("3")  // S'exécute 1er (LIFO)
}
```

### 2. Anonymous structs pour test cases

```go
tests := []struct {  // ✅ Pas besoin de définir un type
    name string
    want int
}{...}
```

### 3. Pointer semantics dans les structs

```go
type Fields struct {
    Assignee *Assignee  // ✅ Pointer = peut être nil (unassigned)
}

// Dans les tests:
if ticket.Fields.Assignee == nil {  // Check nil avant d'accéder
    // Unassigned
}
```

### 4. Context dans les tests

```go
ctx := context.Background()  // Context par défaut pour tests

// Pour tester timeouts:
ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
defer cancel()
```

## Questions à explorer

1. **Que se passe si un test panic?** → Go catch le panic, marque le test FAIL
2. **Comment tester du code concurrent?** → Use race detector: `go test -race`
3. **Comment mocker des interfaces?** → Crée un type qui implémente l'interface
4. **Faut-il tester les fonctions privées?** → Non, teste le comportement public
5. **Quelle couverture viser?** → 70-80% est raisonnable, 100% est over-engineering

## Pattern avancé: Test parallelization

```go
func TestSomething(t *testing.T) {
    t.Parallel()  // ✅ Ce test peut run en parallèle avec d'autres

    tests := []struct{...}{...}

    for _, tt := range tests {
        tt := tt  // ❌ Plus nécessaire en Go 1.22+
        t.Run(tt.name, func(t *testing.T) {
            t.Parallel()  // ✅ Subtests aussi en parallèle
            // Test logic
        })
    }
}
```

**Attention:** N'utilise pas `t.Parallel()` si tes tests modifient un état global (comme viper).

## Fichier créé

Le fichier `tickets_test.go` contient :

✅ **113 lignes de helper functions** (setupViperForTest, makeTicket)
✅ **Tests unitaires** pour FilterByStatus, FilterByAssignee, CalculatePageRequests
✅ **Tests d'intégration** avec mock HTTP pour GetTickets
✅ **Test de comportement** pour le sorting dans FetchSprintTickets

**Total: 6 tests, 18 subtests, tous PASS ✅**
