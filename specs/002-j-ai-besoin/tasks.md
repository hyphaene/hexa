# Tasks: Jira Ticket Fetch Command

**Feature**: `002-j-ai-besoin` - Migrate bash script to native Go CLI command
**Input**: Design documents from `/Users/maximilien/Code/hexa/specs/002-j-ai-besoin/`
**Prerequisites**: plan.md ✅, research.md ✅, data-model.md ✅, contracts/ ✅

## Execution Flow (main)
```
1. Load plan.md from feature directory ✅
   → Tech stack: Go 1.24.4, Cobra, Viper, net/http
   → Structure: Single project (CLI tool)
2. Load optional design documents ✅
   → data-model.md: 4 entities (Ticket, CacheEntry, UserProfile, Sprint)
   → contracts/: jira-api.yaml (3 endpoints), cache-format.json
   → research.md: 5 decisions (API, cache, user profile, directory, status map)
3. Generate tasks by category ✅
   → Setup: 3 tasks (structure, deps, linting)
   → Tests: 9 tasks (3 contract + 6 integration)
   → Core: 10 tasks (4 models + 6 implementation)
   → Integration: 2 tasks (config, error handling)
   → Polish: 3 tasks (unit tests, performance, quickstart)
4. Apply task rules ✅
   → Different files = [P] parallel
   → Same file = sequential
   → Tests before implementation (TDD enforced)
5. Number tasks sequentially (T001-T027) ✅
6. Generate dependency graph ✅
7. Create parallel execution examples ✅
8. Validate task completeness ✅
9. Return: SUCCESS (27 tasks ready for execution)
```

## Format: `[ID] [P?] Description`
- **[P]**: Can run in parallel (different files, no dependencies)
- Include exact file paths in descriptions

## Path Conventions
Single project structure (CLI tool):
- Commands: `cmd/jira/fetch/`
- Internal packages: `internal/jira/`, `internal/cache/`
- Tests: `tests/contract/`, `tests/integration/`, `tests/unit/`

---

## Phase 3.1: Setup

- [x] **T001** Create directory structure for fetch command
  - Create `cmd/jira/fetch/` directory
  - Create `internal/cache/` directory
  - Create `tests/contract/`, `tests/integration/`, `tests/unit/` directories if missing
  - **File impact**: Filesystem only (no code files yet)

- [x] **T002** Verify Go dependencies in `go.mod`
  - Confirm `github.com/spf13/cobra` present
  - Confirm `github.com/spf13/viper` present
  - Run `go mod tidy` to ensure clean state
  - **File impact**: `go.mod`, `go.sum`

- [x] **T003** [P] Configure golangci-lint for project
  - Create `.golangci.yml` with standard Go linting rules if missing
  - Verify `golangci-lint run` passes on existing code
  - **File impact**: `.golangci.yml` (if new)

---

## Phase 3.2: Tests First (TDD) ⚠️ MUST COMPLETE BEFORE 3.3

**CRITICAL: These tests MUST be written and MUST FAIL before ANY implementation**

### Contract Tests (Jira API)

- [ ] **T004** [P] Contract test for GET /rest/agile/1.0/sprint/{sprintId}/issue
  - File: `tests/contract/jira_sprint_issues_test.go`
  - Mock Jira API response matching `contracts/jira-api.yaml` SprintIssuesResponse schema
  - Assert response contains `maxResults`, `startAt`, `total`, `issues[]`
  - Assert each issue has `key`, `fields.summary`, `fields.status.name`
  - Test pagination: startAt=0, maxResults=100
  - **Expected**: Test fails (endpoint not called yet)

- [ ] **T005** [P] Contract test for GET /rest/api/3/myself
  - File: `tests/contract/jira_myself_test.go`
  - Mock Jira API response matching `contracts/jira-api.yaml` UserProfile schema
  - Assert response contains `accountId`, `emailAddress`, `displayName`
  - Test authentication header presence (Bearer token)
  - **Expected**: Test fails (endpoint not called yet)

- [ ] **T006** [P] Contract test for GET /rest/agile/1.0/board/{boardId}/sprint?state=active
  - File: `tests/contract/jira_active_sprint_test.go`
  - Mock Jira API response matching `contracts/jira-api.yaml` SprintListResponse schema
  - Assert response contains `values[]` with active sprint (reuse existing `getCurrentSprintId` logic)
  - Test state filter parameter
  - **Expected**: Test fails (already implemented, but validates contract)

### Integration Tests (Command Flow)

- [ ] **T007** [P] Integration test: Fetch tickets with status filter (cache miss)
  - File: `tests/integration/fetch_command_cache_miss_test.go`
  - Test scenario from `quickstart.md` Scenario 1
  - Mock Jira API responses for sprint ID + sprint issues
  - Run `hexa jira fetch in-progress`
  - Assert output contains ticket keys, summaries, assignees, priorities
  - Assert cache file created at `~/.hexa/cache/sprint_{id}.json`
  - Assert cache age shows "0s"
  - **Expected**: Test fails (command not implemented)

- [ ] **T008** [P] Integration test: Cache hit within TTL
  - File: `tests/integration/fetch_command_cache_hit_test.go`
  - Test scenario from `quickstart.md` Scenario 2
  - Pre-populate cache file with fresh timestamp
  - Run `hexa jira fetch in-progress` (no API call expected)
  - Assert output uses cached data
  - Assert cache age increases
  - **Expected**: Test fails (cache read logic not implemented)

- [ ] **T009** [P] Integration test: Cache expiry after 5 minutes
  - File: `tests/integration/fetch_command_cache_expiry_test.go`
  - Test scenario from `quickstart.md` Scenario 3
  - Pre-populate cache with timestamp >5 minutes old
  - Run `hexa jira fetch in-progress`
  - Assert "Récupération complète" message displayed
  - Assert cache refreshed (timestamp updated)
  - **Expected**: Test fails (TTL check not implemented)

- [ ] **T010** [P] Integration test: "me" filter with user profile fetch
  - File: `tests/integration/fetch_command_me_filter_test.go`
  - Test scenario from `quickstart.md` Scenario 4
  - Mock Jira `/rest/api/3/myself` response
  - Ensure `jira.userEmail` NOT in Viper config initially
  - Run `hexa jira fetch to-test --filter=me`
  - Assert "Fetching user profile" message
  - Assert user email persisted to config
  - Assert only tickets assigned to user email shown
  - **Expected**: Test fails (user profile fetch not implemented)

- [ ] **T011** [P] Integration test: "unassigned" filter
  - File: `tests/integration/fetch_command_unassigned_filter_test.go`
  - Test scenario from `quickstart.md` Scenario 6
  - Run `hexa jira fetch to-do --filter=unassigned`
  - Assert only tickets with null assignee shown
  - Assert "Non assigné" displayed for all results
  - **Expected**: Test fails (unassigned filter not implemented)

- [ ] **T012** [P] Integration test: --no-cache flag bypasses cache
  - File: `tests/integration/fetch_command_no_cache_test.go`
  - Test scenario from `quickstart.md` Scenario 7
  - Pre-populate fresh cache
  - Run `hexa jira fetch in-progress --no-cache`
  - Assert "Cache ignoré" message
  - Assert API call made (mock verifies)
  - Assert cache file updated with new timestamp
  - **Expected**: Test fails (--no-cache flag not implemented)

---

## Phase 3.3: Core Implementation (ONLY after tests are failing)

### Models & Structs

- [x] **T013** [P] Create Ticket, Fields, Status, Assignee, Priority structs
  - File: `internal/jira/ticket.go`
  - Define structs per `data-model.md` section 1
  - Add JSON tags for API unmarshaling
  - Use pointers for optional fields (`*Assignee`, `*Priority`)
  - **Validates**: T004 contract test (JSON unmarshaling)

- [x] **T014** [P] Create CacheEntry struct with TTL methods
  - File: `internal/cache/entry.go`
  - Define struct per `data-model.md` section 2
  - Implement `IsExpired() bool` method (5-minute TTL check)
  - Implement `Age() time.Duration` method
  - **Validates**: T008, T009 integration tests (cache logic)

- [x] **T015** [P] Create UserProfile struct
  - File: `internal/jira/user.go`
  - Define struct per `data-model.md` section 3
  - Add JSON tags for `/rest/api/3/myself` response
  - **Validates**: T005 contract test, T010 integration test

- [x] **T016** [P] Create status key mapping
  - File: `internal/jira/status_map.go`
  - Define `StatusMap` variable per `data-model.md` section 5
  - Implement `MapStatusKey(cliKey string) (string, error)` function
  - Implement `ValidStatusKeys() []string` function for help text
  - **Validates**: Error handling for invalid status keys

### Cache Manager

- [x] **T017** Cache file read/write operations
  - File: `internal/cache/manager.go`
  - Implement `ReadCache(sprintID int) (*CacheEntry, error)` function
  - Implement `WriteCache(entry *CacheEntry) error` function
  - Use `os.UserHomeDir()` + `~/.hexa/cache/sprint_{id}.json` path
  - Create cache directory with `os.MkdirAll` if missing
  - Handle corrupted cache files (return error, caller refreshes)
  - **Validates**: T007, T008, T009, T012 integration tests

- [x] **T018** Cache TTL validation logic
  - File: `internal/cache/manager.go` (add function)
  - Implement `ShouldRefresh(entry *CacheEntry, noCache bool) bool` function
  - Returns true if: cache expired OR noCache flag set OR cache file missing
  - **Validates**: T009 (expiry), T012 (no-cache flag)

### User Profile Fetch

- [x] **T019** Fetch current user profile from Jira API
  - File: `internal/jira/user.go` (add function)
  - Implement `FetchCurrentUser() (*UserProfile, error)` function
  - Call `/rest/api/3/myself` with Bearer token from `viper.GetString("jira.token")`
  - Unmarshal response into `UserProfile` struct
  - **Validates**: T005 contract test, T010 integration test

- [x] **T020** Persist user email to Viper config
  - File: `internal/jira/user.go` (add function)
  - Implement `SaveUserEmail(email string, level string) error` function
  - Use `viper.Set("jira.userEmail", email)` + `viper.WriteConfig()`
  - Support config levels: user, project, project-local (via Viper config path selection)
  - **Validates**: T010 integration test (config persistence)

### Ticket Fetch & Filter Logic

- [x] **T021** Fetch sprint tickets from Jira API with pagination
  - File: `internal/jira/tickets.go`
  - Implement `FetchSprintTickets(sprintID int) ([]Ticket, int, error)` function
  - Call `/rest/agile/1.0/sprint/{sprintId}/issue?startAt=0&maxResults=100`
  - Handle pagination (loop until `isLast=true` or all tickets fetched)
  - Return tickets array + total count
  - **Validates**: T004 contract test, T007 integration test

- [x] **T022** Filter tickets by status
  - File: `internal/jira/tickets.go` (add function)
  - Implement `FilterByStatus(tickets []Ticket, statusName string) []Ticket` function
  - Compare `ticket.Fields.Status.Name` with exact Jira status name
  - **Validates**: T007 integration test (status filtering)

- [x] **T023** Filter tickets by assignee (me, unassigned, all)
  - File: `internal/jira/tickets.go` (add function)
  - Implement `FilterByAssignee(tickets []Ticket, filter string, userEmail string) []Ticket` function
  - "me": match `ticket.Fields.Assignee.EmailAddress == userEmail`
  - "unassigned": `ticket.Fields.Assignee == nil`
  - "all": no filtering
  - **Validates**: T010 (me filter), T011 (unassigned filter)

### CLI Command

- [x] **T024** Create `hexa jira sprint fetch` Cobra subcommand
  - Files: `cmd/jira/sprint/sprint.go`, `cmd/jira/sprint/fetch.go`
  - Define Sprint subcommand root + fetch child command
  - Define Cobra command with:
    - Use: `fetch <status>`
    - Args: Exactly 1 arg (status key)
    - Flags: `--filter` (me|unassigned|all, default: all), `--no-cache` (bool, default: false)
  - Add sprint command to `cmd/jira/jira.go` via `init()` function
  - Implement RunE function:
    1. Validate status key via `MapStatusKey(args[0])`
    2. Get current sprint ID via existing `jira.GetCurrentSprintId()`
    3. Check cache via `cache.ReadCache(sprintID)` + `ShouldRefresh(...)`
    4. If refresh needed: Call `FetchSprintTickets(sprintID)` + `WriteCache(...)`
    5. Filter by status via `FilterByStatus(...)`
    6. If filter="me": Get user email (config or fetch) + `FilterByAssignee(...)`
    7. If filter="unassigned": `FilterByAssignee(..., "unassigned", "")`
    8. Format output (see T025)
  - **Validates**: All integration tests T007-T012

- [x] **T024b** Create `hexa jira sprint pulse` command for multi-status overview
  - File: `cmd/jira/sprint/pulse.go`
  - Single API call fetching all sprint tickets once
  - Filter in-memory for: "To Do" (me), "In Progress" (me), "DEPLOY IN UAT" (all), "Blocked" (all)
  - Display grouped output with counts per status
  - **Use case**: Replace 4 sequential bash script calls with 1 optimized command

---

## Phase 3.4: Integration

- [x] **T025** Output formatting with cache statistics
  - File: `cmd/jira/sprint/fetch.go` (add function)
  - Implement `FormatOutput(tickets []Ticket, cacheAge time.Duration, total int, filterStatus string, filterAssignee string)` function
  - Display format per `quickstart.md` examples:
    ```
    📋 Utilisation du cache (âge: {age})
    🔍 Recherche tickets: {status} (filtre: {filter})

    {KEY} - {Summary} [{Assignee or "Non assigné"}] ({Priority})
    ...

    📊 Total: {count} ticket(s) en status '{status}' (filtre: {filter})
    🔍 Cache: {total} tickets au total dans le sprint
    ```
  - **Validates**: All integration tests (output format)

- [x] **T026** Error handling and user-friendly messages
  - File: `cmd/jira/sprint/fetch.go` (enhance existing error handling)
  - Invalid status key: Show valid keys list (use `ValidStatusKeys()`)
  - Authentication failure (401): Show token config hint
  - Network error: Show URL being accessed, suggest config check
  - Corrupted cache: Warn and auto-refresh
  - No tickets found: Show "Aucun ticket trouvé" (not an error)
  - **Validates**: `quickstart.md` Scenarios 8, 9, 10, 11, 12

---

## Phase 3.5: Polish

- [ ] **T027** [P] Unit tests for cache TTL edge cases
  - File: `tests/unit/cache_test.go`
  - Test `IsExpired()` with: exactly 300s (edge), 299s (valid), 301s (expired)
  - Test `Age()` returns correct duration
  - Test `ShouldRefresh()` with noCache=true ignores TTL
  - **Coverage goal**: >80% for cache package

- [ ] **T028** [P] Unit tests for status mapping and validation
  - File: `tests/unit/status_map_test.go`
  - Test `MapStatusKey()` with all valid keys (11 total)
  - Test invalid key returns error with valid keys list
  - Test case-sensitivity ("In-Progress" vs "in-progress")
  - **Coverage goal**: 100% for status_map.go

- [ ] **T029** [P] Unit tests for ticket filtering logic
  - File: `tests/unit/tickets_test.go`
  - Test `FilterByStatus()` with mixed statuses
  - Test `FilterByAssignee()` with "me", "unassigned", "all" filters
  - Test null assignee handling
  - Test empty ticket array
  - **Coverage goal**: >80% for tickets.go

- [ ] **T030** Run manual quickstart scenarios
  - Execute all 12 scenarios from `specs/002-j-ai-besoin/quickstart.md`
  - Verify actual output matches expected output
  - Test on real Jira instance (not just mocks)
  - Document any deviations or edge cases discovered
  - **Deliverable**: Quickstart validation report (can be comment in PR)

- [ ] **T031** Performance validation
  - Measure cache hit performance: `time hexa jira fetch in-progress` (2nd run)
  - Target: <50ms for cache hits
  - Measure API call performance: `time hexa jira fetch in-progress --no-cache`
  - Target: <2s p95 (depends on network, use median for baseline)
  - **Deliverable**: Performance report in PR description

---

## Dependencies

**Strict TDD Order**:
1. **Setup (T001-T003)** → Must complete first
2. **All Tests (T004-T012)** → Must complete before ANY implementation
3. **Verify tests fail** → Manual check before proceeding
4. **Implementation (T013-T026)** → Only after tests are red
5. **Polish (T027-T031)** → After implementation passes tests

**Implementation Dependencies**:
- T013, T014, T015, T016 are parallel (different files)
- T017, T018 depend on T014 (CacheEntry struct)
- T019, T020 depend on T015 (UserProfile struct)
- T021, T022, T023 depend on T013 (Ticket struct)
- T024 depends on T016, T017, T018, T019, T020, T021, T022, T023 (uses all business logic)
- T025, T026 enhance T024 (same file, sequential)
- T027, T028, T029 are parallel (different test files)
- T030, T031 depend on T024-T026 (working command)

---

## Parallel Execution Examples

### Phase 3.2: All Contract Tests (Parallel)
```bash
# Launch T004-T006 together (3 agents in parallel):
Task 1: "Contract test for GET /rest/agile/1.0/sprint/{sprintId}/issue in tests/contract/jira_sprint_issues_test.go per data-model.md Ticket struct and contracts/jira-api.yaml SprintIssuesResponse schema"

Task 2: "Contract test for GET /rest/api/3/myself in tests/contract/jira_myself_test.go per data-model.md UserProfile struct and contracts/jira-api.yaml UserProfile schema"

Task 3: "Contract test for GET /rest/agile/1.0/board/{boardId}/sprint in tests/contract/jira_active_sprint_test.go per contracts/jira-api.yaml SprintListResponse schema"
```

### Phase 3.2: All Integration Tests (Parallel)
```bash
# Launch T007-T012 together (6 agents in parallel):
Task 1: "Integration test for fetch with cache miss in tests/integration/fetch_command_cache_miss_test.go per quickstart.md Scenario 1"

Task 2: "Integration test for cache hit in tests/integration/fetch_command_cache_hit_test.go per quickstart.md Scenario 2"

Task 3: "Integration test for cache expiry in tests/integration/fetch_command_cache_expiry_test.go per quickstart.md Scenario 3"

Task 4: "Integration test for me filter with user profile fetch in tests/integration/fetch_command_me_filter_test.go per quickstart.md Scenario 4"

Task 5: "Integration test for unassigned filter in tests/integration/fetch_command_unassigned_filter_test.go per quickstart.md Scenario 6"

Task 6: "Integration test for --no-cache flag in tests/integration/fetch_command_no_cache_test.go per quickstart.md Scenario 7"
```

### Phase 3.3: All Model Creation (Parallel)
```bash
# Launch T013-T016 together (4 agents in parallel):
Task 1: "Create Ticket, Fields, Status, Assignee, Priority structs in internal/jira/ticket.go per data-model.md section 1 with JSON tags"

Task 2: "Create CacheEntry struct with IsExpired() and Age() methods in internal/cache/entry.go per data-model.md section 2"

Task 3: "Create UserProfile struct in internal/jira/user.go per data-model.md section 3 with JSON tags"

Task 4: "Create StatusMap and MapStatusKey() function in internal/jira/status_map.go per data-model.md section 5"
```

### Phase 3.5: All Unit Tests (Parallel)
```bash
# Launch T027-T029 together (3 agents in parallel):
Task 1: "Unit tests for cache TTL edge cases in tests/unit/cache_test.go testing IsExpired(), Age(), ShouldRefresh()"

Task 2: "Unit tests for status mapping in tests/unit/status_map_test.go testing MapStatusKey() with all valid and invalid keys"

Task 3: "Unit tests for ticket filtering in tests/unit/tickets_test.go testing FilterByStatus() and FilterByAssignee()"
```

---

## Notes

- **[P] tasks**: Different files, no shared dependencies, safe for parallel execution
- **TDD enforcement**: ALL tests (T004-T012) must fail before starting implementation (T013+)
- **Commit strategy**: Commit after each task (or logical group of parallel tasks)
- **Go conventions**: Run `go fmt ./...`, `go vet ./...`, `golangci-lint run` before each commit
- **Cache directory**: Created at runtime by T017, no manual setup needed
- **User email**: Auto-fetched on first "me" filter (T019+T020), cached in Viper config
- **Error messages**: Must match `quickstart.md` examples for consistency

---

## Validation Checklist
*GATE: Verify before marking tasks.md complete*

- [x] All contracts (jira-api.yaml endpoints) have corresponding contract tests (T004-T006)
- [x] All entities (Ticket, CacheEntry, UserProfile, StatusMap) have model creation tasks (T013-T016)
- [x] All integration scenarios from quickstart.md have test tasks (T007-T012)
- [x] All tests (T004-T012) come before implementation tasks (T013-T026)
- [x] Parallel tasks ([P]) operate on different files (no conflicts)
- [x] Each task specifies exact file path (✅ all tasks include paths)
- [x] No [P] task modifies same file as another [P] task (✅ verified)
- [x] CLI command (T024) depends on all business logic (T013-T023)
- [x] Performance goals documented (T031: <50ms cache hit, <2s API call)

**Status**: ✅ 31 tasks generated, ready for execution
