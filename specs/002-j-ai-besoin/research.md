# Research: Jira Ticket Fetch Command

**Feature**: `002-j-ai-besoin` | **Date**: 2025-10-01

## Research Questions Resolved

### 1. Jira REST API Endpoint Selection

**Question**: Which Jira API endpoint should be used to fetch tickets from a sprint?

**Decision**: `/rest/agile/1.0/sprint/{sprintId}/issue`

**Rationale**:
- Direct sprint-to-tickets relationship (no JQL complexity)
- Returns all required fields: `key`, `summary`, `status`, `assignee`, `priority`
- Built-in pagination support via `startAt` and `maxResults` query parameters
- Consistent with existing `GetCurrentSprintId()` which uses Agile API v1.0

**Alternatives Considered**:
- **JQL via `/rest/api/3/search`**: More flexible but requires constructing JQL query string (`sprint = {id}`). Unnecessary complexity for sprint-scoped queries. Rejected: adds error surface for JQL syntax.
- **Board issues endpoint `/rest/agile/1.0/board/{boardId}/issue`**: Requires additional sprint filtering client-side. Rejected: inefficient, fetches irrelevant tickets.

**Implementation Notes**:
- Query params: `startAt=0&maxResults=100` for pagination (max 100 per Jira API limit)
- Response format: `{ "maxResults": 100, "startAt": 0, "total": 250, "issues": [...] }`
- Requires `Authorization: Bearer {token}` header (from `HEXA_JIRA_TOKEN`)

---

### 2. Cache File Format

**Question**: What format should be used for the sprint ticket cache file?

**Decision**: JSON with metadata envelope

**Rationale**:
- Native Go `encoding/json` support (stdlib, no dependencies)
- Human-readable for debugging (`cat ~/.hexa/cache/sprint_123.json`)
- Direct unmarshal from Jira API response structure
- Metadata (timestamp, sprint ID) embedded in same file for atomic reads

**Alternatives Considered**:
- **MessagePack**: Binary format, faster parsing. Rejected: adds `github.com/vmihailenco/msgpack` dependency, violates "minimize dependencies" principle from constitution.
- **Gob encoding**: Go-native binary format. Rejected: not human-readable, makes cache debugging difficult during Go learning phase.
- **SQLite**: Queryable cache. Rejected: overkill for simple key-value (sprint ID → tickets), adds `modernc.org/sqlite` dependency and SQL learning curve distraction.

**File Structure**:
```json
{
  "sprintId": 123,
  "cachedAt": "2025-10-01T14:30:00Z",
  "ttlSeconds": 300,
  "total": 47,
  "issues": [
    {
      "key": "PROJ-123",
      "fields": {
        "summary": "...",
        "status": { "name": "In Progress" },
        "assignee": { "emailAddress": "...", "displayName": "..." },
        "priority": { "name": "High" }
      }
    }
  ]
}
```

**Implementation Notes**:
- Use `time.Time` for `cachedAt`, JSON marshal handles RFC3339 format
- TTL check: `time.Since(cachedAt) < 5 * time.Minute`
- File permissions: `0644` (user read/write, others read-only)

---

### 3. User Profile Fetch Endpoint

**Question**: How should the system retrieve the authenticated user's email for "me" filter?

**Decision**: `/rest/api/3/myself`

**Rationale**:
- Single API call returns authenticated user profile
- Includes `emailAddress` field directly (no parsing required)
- REST API v3 endpoint (current Jira Cloud standard)
- Same authorization as other calls (Bearer token)

**Alternatives Considered**:
- **Parse from ticket assignee**: Extract email from first assigned ticket. Rejected: fragile (requires existing assigned tickets), extra API call dependency.
- **Hardcode in config**: User manually sets `HEXA_JIRA_USER_EMAIL`. Rejected: defeats auto-discovery goal from FR-003a, poor UX.
- **Infer from git config**: Use `git config user.email`. Rejected: assumes git email matches Jira email (often false in corporate environments).

**Response Example**:
```json
{
  "self": "https://...",
  "accountId": "...",
  "emailAddress": "user@example.com",
  "displayName": "John Doe"
}
```

**Implementation Notes**:
- Call on first "me" filter use when `viper.GetString("jira.userEmail")` is empty
- Persist to config via `viper.Set("jira.userEmail", email)` + `viper.WriteConfig()`
- Config level selection: prompt user or default to project-local (`.hexa.local.yml`)

---

### 4. Cache Directory Location

**Question**: Where should sprint ticket cache files be stored?

**Decision**: `~/.hexa/cache/sprint_{sprintId}.json`

**Rationale**:
- Standard user-scoped cache pattern (XDG-like without XDG complexity)
- Cross-platform via `os.UserHomeDir()` (works on macOS, Linux, Windows)
- Persists across terminal sessions (enables 5-minute TTL goal)
- Go stdlib `os.MkdirAll` handles directory creation atomically

**Alternatives Considered**:
- **`os.TempDir()`** (e.g., `/tmp/`): Rejected: cleared on reboot, defeats 5-minute TTL persistence goal.
- **Current working directory** (`.hexa/cache/`): Rejected: pollutes project directories, per-project caching unnecessary for sprint data.
- **XDG Base Directory** (`$XDG_CACHE_HOME` or `~/.cache/hexa/`): Rejected: XDG not universally adopted (macOS uses `~/Library/Caches/`), adds OS detection complexity.

**Directory Structure**:
```
~/.hexa/
├── cache/
│   ├── sprint_123.json
│   ├── sprint_124.json
│   └── sprint_125.json
└── config files (.hexa.yml, etc.)
```

**Implementation Notes**:
- Create directory: `os.MkdirAll(filepath.Join(home, ".hexa", "cache"), 0755)`
- Error handling: If home dir unavailable, fallback to `os.TempDir()` with warning log
- Cache file naming: `sprint_{id}.json` (simple, predictable, avoids timestamp drift)

---

### 5. Status Key Mapping Strategy

**Question**: How should user-friendly CLI status keys map to exact Jira status names?

**Decision**: Hardcoded `map[string]string` in Go

**Rationale**:
- Jira API requires exact status name match (case-sensitive: "To test" not "to-test")
- CLI should accept kebab-case keys (user-friendly: `to-test`, `in-progress`)
- Mapping is stable (status names rarely change in Jira workflow)
- Go map lookup is O(1), no performance penalty

**Alternatives Considered**:
- **Direct passthrough**: User types exact Jira status name. Rejected: poor UX (case-sensitive, spaces in CLI args require quoting).
- **API call to fetch all statuses**: Query `/rest/api/3/status` on startup. Rejected: unnecessary HTTP roundtrip (200ms+ latency), complicates offline cache-hit scenario.
- **Fuzzy matching**: Lowercase + normalize spaces. Rejected: ambiguous ("totest" → "To Test" or "To test"?), adds string processing complexity.

**Mapping Table**:
| CLI Key       | Jira Status Name |
|---------------|------------------|
| `to-do`       | To Do            |
| `in-progress` | In Progress      |
| `to-test`     | To test          |
| `uat`         | UAT              |
| `deploy-uat`  | DEPLOY IN UAT    |
| `to-deploy`   | To deploy        |
| `blocked`     | Blocked          |
| `prep`        | Prep             |
| `new`         | New              |
| `closed`      | Closed           |
| `archived`    | Archived         |

**Implementation Notes**:
```go
var statusMap = map[string]string{
    "to-do":       "To Do",
    "in-progress": "In Progress",
    // ...
}

func mapStatusKey(cliKey string) (string, error) {
    if jiraName, ok := statusMap[cliKey]; ok {
        return jiraName, nil
    }
    return "", fmt.Errorf("invalid status key '%s', valid keys: %v", cliKey, validKeys())
}
```

---

## Go Learning Highlights

### Concepts Reinforced by Research
1. **HTTP client patterns**: `http.NewRequest`, header manipulation, response handling
2. **JSON marshaling**: `encoding/json` with struct tags, `time.Time` RFC3339 handling
3. **File I/O**: `os.UserHomeDir()`, `os.MkdirAll()`, `ioutil.ReadFile` (or `os.ReadFile` in Go 1.16+)
4. **Error wrapping**: `fmt.Errorf("context: %w", err)` for error chains
5. **Map operations**: Constant-time lookups, existence checking with comma-ok idiom

### Standard Library Usage
- `net/http`: API client (no third-party HTTP library needed)
- `encoding/json`: Serialization (no MessagePack/protobuf)
- `os`: Filesystem operations (no `afero` abstraction)
- `time`: TTL calculations with `time.Since()`, `time.Duration`
- `path/filepath`: Cross-platform path joining

---

## Dependencies: No New Additions Required

All research decisions use existing dependencies:
- ✅ `net/http`: Go stdlib
- ✅ `encoding/json`: Go stdlib
- ✅ `os`, `time`, `path/filepath`: Go stdlib
- ✅ `github.com/spf13/viper`: Already in `go.mod` (config management)
- ✅ `github.com/spf13/cobra`: Already in `go.mod` (CLI framework)

**Constitution Compliance**: Minimized dependencies ✅ (Principle I: Go Learning First)

---

## Next Steps

Phase 0 complete. Proceed to Phase 1:
1. Generate `data-model.md` (entities + Go structs)
2. Generate `contracts/` (OpenAPI + JSON schema)
3. Generate `quickstart.md` (manual test scenarios)
4. Update `AGENTS.md` (add Jira API context)
5. Write failing tests (contract + unit + integration)
