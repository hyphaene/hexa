# Implementation Plan: Jira Ticket Fetch Command

**Branch**: `002-j-ai-besoin` | **Date**: 2025-10-01 | **Spec**: [spec.md](./spec.md)
**Input**: Feature specification from `/Users/maximilien/Code/hexa/specs/002-j-ai-besoin/spec.md`

## Execution Flow (/plan command scope)
```
1. Load feature spec from Input path
   → ✅ Loaded - Jira fetch command migration from bash script
2. Fill Technical Context (scan for NEEDS CLARIFICATION)
   → ✅ All clarified via Session 2025-10-01
   → Project Type: single (CLI tool)
3. Fill the Constitution Check section based on constitution
   → ✅ Constitution v1.0.0 loaded
4. Evaluate Constitution Check section
   → In progress
5. Execute Phase 0 → research.md
   → Pending
6. Execute Phase 1 → contracts, data-model.md, quickstart.md, AGENTS.md
   → Pending
7. Re-evaluate Constitution Check section
   → Pending
8. Plan Phase 2 → Describe task generation approach
   → Pending
9. STOP - Ready for /tasks command
```

**IMPORTANT**: The /plan command STOPS at step 8. Phases 2-4 are executed by other commands:
- Phase 2: /tasks command creates tasks.md
- Phase 3-4: Implementation execution (manual or via tools)

## Summary
Migrate existing bash script (`jira_fetch.sh`) functionality to native Go CLI command `hexa jira fetch`. Users can query current sprint tickets filtered by status (to-do, in-progress, etc.) and assignee (me, unassigned, all). System caches ticket data for 5 minutes to reduce API calls, with bypass flag for force refresh. User profile email auto-fetched on first "me" filter use and persisted to Viper config.

## Technical Context
**Language/Version**: Go 1.24.4
**Primary Dependencies**: Cobra (CLI), Viper (config), godotenv (env), net/http (Jira API)
**Storage**: Filesystem cache at `~/.hexa/cache/sprint_{id}.json`, Viper config persistence
**Testing**: Go standard testing (`go test ./...`), TDD with Red-Green-Refactor
**Target Platform**: macOS, Linux (cross-platform via GoReleaser)
**Project Type**: single (CLI tool with Cobra structure)
**Performance Goals**: Cache hits <50ms, API calls <2s p95, startup <100ms
**Constraints**: Single binary, no runtime dependencies, 5-minute cache TTL
**Scale/Scope**: Sprint data ~100-500 tickets typical, pagination required for large sprints

## Constitution Check
*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

### I. Go Learning First ✅
- **Teaches**: HTTP client patterns, JSON decoding, file I/O (cache), Viper config, Cobra subcommands
- **Primitives**: `net/http`, `encoding/json`, `os`, `time`, error handling patterns
- **Complexity justified**: Cache logic teaches file operations + time-based TTL logic

### II. CLI-First Architecture ✅
- **Command structure**: `hexa jira fetch <status> --filter=<me|unassigned|all> [--no-cache]`
- **Flags**: `--filter` (local), `--no-cache` (local), inherits global Viper config
- **Output**: Human-readable table format with cache statistics
- **STDIN/STDOUT**: Ticket list → STDOUT, errors → STDERR

### III. Test-Driven Development ✅
- **TDD enforced**: Contract tests for Jira API responses, unit tests for cache logic, integration tests for full command flow
- **Coverage targets**: Core cache/filter logic >80%, CLI command >60%
- **Red-Green-Refactor**: Tests written first, validated to fail, then implement

### IV. Distribution Simplicity ✅
- **Single binary**: Cache directory created at runtime if missing
- **No dependencies**: Uses existing Viper config, no new external tools
- **Homebrew tap**: GoReleaser auto-updates on release

### V. Configuration Hierarchy ✅
- **Config keys**: `HEXA_JIRA_TOKEN`, `HEXA_JIRA_URL`, `HEXA_JIRA_USER_EMAIL` (auto-populated)
- **Precedence**: CLI flags → env vars → `.hexa.local.yml` → `.hexa.yml` → defaults
- **Auto-fetch**: User email fetched from Jira API on first "me" filter, persisted to chosen config level

**Gate Status**: ✅ PASS - All principles aligned, justified complexity

## Project Structure

### Documentation (this feature)
```
specs/002-j-ai-besoin/
├── plan.md              # This file (/plan command output)
├── spec.md              # Feature specification (input)
├── research.md          # Phase 0 output (/plan command)
├── data-model.md        # Phase 1 output (/plan command)
├── quickstart.md        # Phase 1 output (/plan command)
├── contracts/           # Phase 1 output (/plan command)
│   ├── jira-api.yaml   # Jira REST API contract (OpenAPI subset)
│   └── cache-format.json # Cache file schema
└── tasks.md             # Phase 2 output (/tasks command - NOT created by /plan)
```

### Source Code (repository root)
```
cmd/
├── jira/
│   ├── jira.go          # Existing parent command
│   ├── init.go          # Existing init command
│   └── fetch/
│       └── fetch.go     # New fetch subcommand

internal/
├── jira/
│   ├── getCurrentSprintId.go  # Existing (reuse)
│   ├── getBoardIdFromName.go  # Existing (reuse)
│   ├── client.go              # New: Jira API client wrapper
│   ├── tickets.go             # New: Ticket fetch/filter logic
│   └── user.go                # New: User profile fetch
└── cache/
    ├── manager.go       # New: Cache TTL and file operations
    └── sprint.go        # New: Sprint ticket cache type

tests/
├── contract/
│   └── jira_api_test.go        # Contract tests for Jira API
├── integration/
│   └── fetch_command_test.go   # End-to-end command tests
└── unit/
    ├── cache_test.go           # Cache manager unit tests
    ├── tickets_test.go         # Ticket filter logic tests
    └── user_test.go            # User profile fetch tests
```

**Structure Decision**: Single project structure (default). New `fetch` subcommand under existing `cmd/jira/` hierarchy. Internal packages separated by concern (`jira` for API, `cache` for filesystem). Tests organized by type (contract/integration/unit) following TDD principle.

## Phase 0: Outline & Research
**Status**: Not started

### Research Tasks
1. **Jira REST API patterns**:
   - Decision: Use `/rest/agile/1.0/sprint/{sprintId}/issue` endpoint for ticket fetch
   - Rationale: Supports pagination via `startAt`/`maxResults`, includes all required fields (key, summary, status, assignee, priority)
   - Alternatives considered: JQL search (more complex, unnecessary for sprint-scoped queries)

2. **Cache file format**:
   - Decision: JSON with metadata (timestamp, sprint ID, total count) + issues array
   - Rationale: Direct unmarshal from Jira API response, standard Go `encoding/json`
   - Alternatives considered: MessagePack (adds dependency), gob (not human-readable for debugging)

3. **User profile fetch endpoint**:
   - Decision: Use `/rest/api/3/myself` endpoint
   - Rationale: Returns authenticated user email, single API call
   - Alternatives considered: Parse from ticket assignee (fragile, requires ticket data first)

4. **Cache directory location**:
   - Decision: `~/.hexa/cache/` with `os.UserHomeDir()` + `os.MkdirAll`
   - Rationale: Standard XDG-like pattern, cross-platform via Go stdlib
   - Alternatives considered: `os.TempDir()` (doesn't persist across reboots, defeats 5min TTL goal)

5. **Status key mapping**:
   - Decision: Map CLI keys (kebab-case) to Jira display names (title case) via `map[string]string`
   - Rationale: User-friendly CLI (`to-test`), Jira API requires exact match ("To test")
   - Alternatives considered: Direct passthrough (breaks on case mismatch), API call to fetch all statuses (unnecessary roundtrip)

**Output**: `research.md` with consolidated decisions

## Phase 1: Design & Contracts
**Status**: Not started

### Artifacts to Generate
1. **data-model.md**: Entities (Sprint, Ticket, CacheEntry, UserProfile) with Go struct definitions
2. **contracts/jira-api.yaml**: OpenAPI snippet for 3 endpoints (sprint issues, myself, active sprint)
3. **contracts/cache-format.json**: JSON schema for cache file structure
4. **quickstart.md**: Manual test scenarios (fresh cache, expired cache, "me" filter first-time)
5. **AGENTS.md**: Update with Jira API context, cache patterns, test-first reminder

### Test Generation Strategy
- **Contract tests**: Mock Jira API responses, assert JSON structure matches OpenAPI
- **Unit tests**: Cache TTL expiry logic, status key mapping, ticket filtering by assignee
- **Integration tests**: Full command flow with mocked HTTP client, verify output format

**Output**: All Phase 1 artifacts generated, tests failing (no implementation yet)

## Phase 2: Task Planning Approach
*This section describes what the /tasks command will do - DO NOT execute during /plan*

**Task Generation Strategy**:
1. **From contracts**: 3 contract test tasks [P] (sprint issues, myself, active sprint)
2. **From data model**: 4 model creation tasks [P] (Sprint, Ticket, CacheEntry, UserProfile structs)
3. **From cache logic**: 2 cache manager tasks (TTL check, read/write operations)
4. **From CLI**: 1 Cobra subcommand skeleton task
5. **From business logic**: 4 implementation tasks (fetch tickets, filter by status, filter by assignee, user profile fetch)
6. **From quickstart**: 3 integration test tasks (scenarios from quickstart.md)
7. **From output formatting**: 1 table formatter task

**Ordering Strategy**:
- Phase A [P]: Contract tests (3), model structs (4) — parallel, no dependencies
- Phase B: Cache manager (2) — depends on CacheEntry model
- Phase C: User profile fetch (1) — depends on UserProfile model
- Phase D: Ticket fetch + filters (3) — depends on Ticket model, cache manager
- Phase E: CLI command (1) — depends on all business logic
- Phase F: Integration tests (3) — depends on CLI command
- Phase G: Output formatter (1) — refactor, depends on working CLI

**Estimated Output**: 22 numbered, dependency-ordered tasks in tasks.md

**IMPORTANT**: This phase is executed by the /tasks command, NOT by /plan

## Phase 3+: Future Implementation
*These phases are beyond the scope of the /plan command*

**Phase 3**: Task execution (/tasks command creates tasks.md)
**Phase 4**: Implementation (execute tasks.md following constitutional principles)
**Phase 5**: Validation (run tests, execute quickstart.md, verify cache behavior)

## Complexity Tracking
*No constitutional violations - table omitted*

## Progress Tracking
*This checklist is updated during execution flow*

**Phase Status**:
- [x] Phase 0: Research complete (/plan command)
- [x] Phase 1: Design complete (/plan command)
- [x] Phase 2: Task planning complete (/plan command - describe approach only)
- [ ] Phase 3: Tasks generated (/tasks command)
- [ ] Phase 4: Implementation complete
- [ ] Phase 5: Validation passed

**Gate Status**:
- [x] Initial Constitution Check: PASS
- [x] Post-Design Constitution Check: PASS (no new violations)
- [x] All NEEDS CLARIFICATION resolved (via spec clarifications)
- [x] Complexity deviations documented (none)

---
*Based on Constitution v1.0.0 - See `.specify/memory/constitution.md`*
