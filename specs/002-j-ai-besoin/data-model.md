# Data Model: Jira Ticket Fetch Command

**Feature**: `002-j-ai-besoin` | **Date**: 2025-10-01

## Entity Overview

This feature introduces 4 core entities and 1 value object:

1. **Ticket** - Jira issue with status, assignee, priority
2. **CacheEntry** - Wrapper for cached sprint ticket data with TTL metadata
3. **UserProfile** - Authenticated Jira user information
4. **Sprint** - Existing entity (reused from `internal/jira/getCurrentSprintId.go`)
5. **Status** - Value object (string constant mapping)

---

## 1. Ticket

### Purpose
Represents a Jira issue within a sprint, containing fields required for display and filtering.

### Go Struct Definition
```go
package jira

import "time"

// Ticket represents a single Jira issue with relevant fields for display and filtering
type Ticket struct {
	Key    string `json:"key"`    // e.g., "PROJ-123"
	Fields Fields `json:"fields"`
}

// Fields contains the nested field structure from Jira API response
type Fields struct {
	Summary  string    `json:"summary"`
	Status   Status    `json:"status"`
	Assignee *Assignee `json:"assignee"` // Pointer: null when unassigned
	Priority *Priority `json:"priority"` // Pointer: null when no priority set
}

// Status represents the workflow status of a ticket
type Status struct {
	Name string `json:"name"` // e.g., "In Progress", "To test"
}

// Assignee represents the user assigned to a ticket
type Assignee struct {
	DisplayName  string `json:"displayName"`  // e.g., "John Doe"
	EmailAddress string `json:"emailAddress"` // e.g., "john.doe@example.com"
}

// Priority represents the priority level of a ticket
type Priority struct {
	Name string `json:"name"` // e.g., "High", "Medium", "Low"
}
```

### Validation Rules
- **Key**: Non-empty, matches pattern `[A-Z]+-\d+` (validated by Jira API, trust response)
- **Summary**: Non-empty (Jira enforces this)
- **Status.Name**: Non-empty, matches one of known Jira workflow statuses
- **Assignee**: Optional (nil when unassigned)
- **Priority**: Optional (nil when no priority)

### Relationships
- Belongs to **Sprint** (implicitly via API query, not stored in struct)
- Has one **Status** (embedded value object)
- Has zero or one **Assignee** (optional pointer)
- Has zero or one **Priority** (optional pointer)

### State Transitions
N/A - Tickets are read-only in this feature. Status changes handled by separate `hexa jira ticket move` command.

---

## 2. CacheEntry

### Purpose
Wraps a collection of tickets with metadata for TTL-based cache invalidation.

### Go Struct Definition
```go
package cache

import "time"

// CacheEntry represents cached sprint ticket data with expiry metadata
type CacheEntry struct {
	SprintID   int                `json:"sprintId"`   // Sprint ID this cache represents
	CachedAt   time.Time          `json:"cachedAt"`   // Timestamp when cache was created
	TTLSeconds int                `json:"ttlSeconds"` // Time-to-live in seconds (default: 300)
	Total      int                `json:"total"`      // Total ticket count in sprint
	Issues     []jira.Ticket      `json:"issues"`     // Ticket array (see Ticket struct above)
}

// IsExpired checks if cache has exceeded its TTL
func (c *CacheEntry) IsExpired() bool {
	return time.Since(c.CachedAt) > time.Duration(c.TTLSeconds) * time.Second
}

// Age returns human-readable cache age (e.g., "2m30s")
func (c *CacheEntry) Age() time.Duration {
	return time.Since(c.CachedAt)
}
```

### Validation Rules
- **SprintID**: Must be positive integer
- **CachedAt**: Must not be in the future
- **TTLSeconds**: Must be positive (default: 300)
- **Total**: Must be non-negative
- **Issues**: Length must match `Total` when fully paginated

### File Storage
- Path: `~/.hexa/cache/sprint_{sprintId}.json`
- Permissions: `0644` (user rw, others r)
- Format: JSON with RFC3339 timestamps

### Lifecycle
1. **Create**: On first API call or cache miss
2. **Read**: On subsequent calls within TTL window
3. **Invalidate**: After 5 minutes (300 seconds) or when `--no-cache` flag used
4. **Cleanup**: Manual (future: add `hexa cache clean` command)

---

## 3. UserProfile

### Purpose
Stores authenticated Jira user information for "me" filter functionality.

### Go Struct Definition
```go
package jira

// UserProfile represents the authenticated Jira user
type UserProfile struct {
	AccountID    string `json:"accountId"`    // Unique Jira account ID
	EmailAddress string `json:"emailAddress"` // User's email (used for "me" filter)
	DisplayName  string `json:"displayName"`  // Full name for display
}
```

### Validation Rules
- **AccountID**: Non-empty, Jira-assigned unique identifier
- **EmailAddress**: Must be valid email format (validated by Jira API)
- **DisplayName**: Non-empty

### Persistence
- **First fetch**: Call `/rest/api/3/myself`, extract `emailAddress`
- **Storage**: Persist to Viper config via `viper.Set("jira.userEmail", email)`
- **Config level**: User choice (project-local `.hexa.local.yml` recommended)
- **Reuse**: Subsequent "me" filter calls read from config (no API call)

### Relationships
- Matches **Ticket.Assignee.EmailAddress** for "me" filter logic

---

## 4. Sprint (Existing, Reused)

### Purpose
Represents a Jira sprint (already defined in `internal/jira/getCurrentSprintId.go`).

### Go Struct Definition (Reference)
```go
package jira

import "time"

// Sprint represents a Jira sprint (existing struct)
type Sprint struct {
	ID            int       `json:"id"`
	Self          string    `json:"self"`
	State         string    `json:"state"` // "active", "future", "closed"
	Name          string    `json:"name"`
	StartDate     time.Time `json:"startDate"`
	EndDate       time.Time `json:"endDate"`
	ActivatedDate time.Time `json:"activatedDate"`
	OriginBoardID int       `json:"originBoardId"`
	Goal          string    `json:"goal"`
	Synced        bool      `json:"synced"`
	AutoStartStop bool      `json:"autoStartStop"`
}
```

### Usage in This Feature
- **Reuse**: Call existing `jira.GetCurrentSprintId()` to get active sprint ID
- **No modifications**: Sprint struct remains unchanged
- **Relationship**: One Sprint has many Tickets (via API query, not modeled in struct)

---

## 5. Status Mapping (Value Object)

### Purpose
Maps CLI-friendly status keys to exact Jira status names (case-sensitive).

### Implementation
```go
package jira

// StatusMap provides CLI key → Jira status name mapping
var StatusMap = map[string]string{
	"to-do":       "To Do",
	"in-progress": "In Progress",
	"to-test":     "To test",
	"uat":         "UAT",
	"deploy-uat":  "DEPLOY IN UAT",
	"to-deploy":   "To deploy",
	"blocked":     "Blocked",
	"prep":        "Prep",
	"new":         "New",
	"closed":      "Closed",
	"archived":    "Archived",
}

// ValidStatusKeys returns all valid CLI status keys for help text
func ValidStatusKeys() []string {
	keys := make([]string, 0, len(StatusMap))
	for k := range StatusMap {
		keys = append(keys, k)
	}
	return keys
}

// MapStatusKey converts CLI key to Jira status name
func MapStatusKey(cliKey string) (string, error) {
	if jiraName, ok := StatusMap[cliKey]; ok {
		return jiraName, nil
	}
	return "", fmt.Errorf("invalid status key '%s', valid keys: %v", cliKey, ValidStatusKeys())
}
```

### Validation Rules
- CLI key must exist in `StatusMap`
- Jira status names are exact matches (case-sensitive)
- Invalid key returns error with valid key list

---

## Entity Diagram

```
┌─────────────────┐
│     Sprint      │ (existing, reused)
│   - ID: int     │
└────────┬────────┘
         │ 1
         │
         │ N
┌────────▼────────┐       ┌──────────────┐
│     Ticket      │───────│    Status    │ (value object, embedded)
│   - Key: string │ 1   1 │  - Name: str │
│   - Fields      │       └──────────────┘
└────────┬────────┘
         │ Fields
         ├── Summary: string
         ├── Status: Status
         ├── Assignee: *Assignee (optional)
         └── Priority: *Priority (optional)

┌─────────────────┐       ┌──────────────┐
│   CacheEntry    │──────▶│   Ticket[]   │
│  - SprintID     │ 1   N │              │
│  - CachedAt     │       └──────────────┘
│  - TTLSeconds   │
│  - Total        │
│  - Issues       │
└─────────────────┘

┌─────────────────┐       ┌──────────────┐
│  UserProfile    │       │   Assignee   │
│  - EmailAddress │◀──────│  - EmailAddr │ (matches for "me" filter)
│  - DisplayName  │       │  - DispName  │
└─────────────────┘       └──────────────┘
```

---

## Go Packages Organization

```
internal/
├── jira/
│   ├── ticket.go       # Ticket, Fields, Status, Assignee, Priority structs
│   ├── user.go         # UserProfile struct + FetchCurrentUser()
│   ├── status_map.go   # StatusMap + MapStatusKey()
│   └── (existing files: getCurrentSprintId.go, Sprint struct)
│
└── cache/
    ├── entry.go        # CacheEntry struct + IsExpired(), Age() methods
    └── manager.go      # Cache file I/O (Read/Write/Delete)
```

**Rationale**: Separation by domain concern. `jira` package handles API types, `cache` package handles persistence. Tests mirror this structure.

---

## Learning Highlights

### Go Concepts Demonstrated
1. **Struct composition**: `Ticket` embeds `Fields`, `Fields` embeds `Status`
2. **Pointer semantics**: `*Assignee`, `*Priority` for optional fields (nil = absent)
3. **JSON tags**: Automatic marshaling/unmarshaling with `json:"fieldName"`
4. **Methods on structs**: `IsExpired()`, `Age()` on `CacheEntry`
5. **Package-level variables**: `StatusMap` as constant lookup table
6. **Error handling**: `MapStatusKey` returns `(string, error)` tuple

### Standard Library Usage
- `time.Time`: Timestamps and duration calculations
- `encoding/json`: Struct serialization with tags
- `fmt.Errorf`: Error message formatting

---

## Next Steps

1. Generate contracts (`contracts/jira-api.yaml`, `contracts/cache-format.json`)
2. Generate quickstart scenarios (`quickstart.md`)
3. Write failing tests for each struct and method
4. Update `AGENTS.md` with data model context
