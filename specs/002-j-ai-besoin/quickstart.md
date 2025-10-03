# Quickstart: Jira Ticket Fetch Command

**Feature**: `002-j-ai-besoin` | **Date**: 2025-10-01

## Purpose
Manual test scenarios to validate `hexa jira fetch` command behavior. Execute these after implementation to confirm all acceptance criteria are met.

---

## Prerequisites

1. **Jira Configuration**:
   ```bash
   # Set required config values
   hexa config user set jira.url "https://your-domain.atlassian.net"
   hexa config local set jira.token "your-pat-token-here"
   hexa config user set jira.boardName "Your Board Name"
   ```

2. **Verify Active Sprint Exists**:
   ```bash
   # Should return current sprint ID
   hexa jira sprint
   ```
   Expected: Sprint ID (e.g., `123`) or error if no active sprint.

3. **Clear Cache** (for fresh test):
   ```bash
   rm -rf ~/.hexa/cache/sprint_*.json
   ```

---

## Scenario 1: Basic Ticket Fetch (All Tickets in Status)

### Test: Fetch all "In Progress" tickets

**Command**:
```bash
hexa jira fetch in-progress
```

**Expected Output**:
```
📋 Utilisation du cache (âge: 0s)
🔍 Recherche tickets: In Progress (filtre: all)

PROJ-123 - Implement user authentication [John Doe] (High)
PROJ-125 - Refactor API endpoints [Jane Smith] (Medium)
PROJ-127 - Update documentation [Non assigné] (Low)

📊 Total: 3 ticket(s) en status 'In Progress' (filtre: all)
🔍 Cache: 47 tickets au total dans le sprint
```

**Validation**:
- ✅ Status "In Progress" matches CLI key `in-progress`
- ✅ Displays: key, summary, assignee (or "Non assigné"), priority
- ✅ Shows cache age (0s on first run)
- ✅ Shows total count of filtered tickets
- ✅ Shows total tickets in sprint

---

## Scenario 2: Cache Hit (Within 5 Minutes)

### Test: Run same query twice within 5 minutes

**Commands**:
```bash
hexa jira fetch in-progress   # First call (cache miss)
sleep 30                       # Wait 30 seconds
hexa jira fetch in-progress   # Second call (cache hit)
```

**Expected Output (Second Call)**:
```
📋 Utilisation du cache (âge: 30s)
🔍 Recherche tickets: In Progress (filtre: all)

[... same ticket list as before ...]

📊 Total: 3 ticket(s) en status 'In Progress' (filtre: all)
🔍 Cache: 47 tickets au total dans le sprint
```

**Validation**:
- ✅ Cache age increases (30s instead of 0s)
- ✅ No API call made (verify via network monitor or logs)
- ✅ Results identical to first call

---

## Scenario 3: Cache Expiry (After 5 Minutes)

### Test: Verify cache auto-refresh after TTL

**Commands**:
```bash
hexa jira fetch in-progress   # Create cache
sleep 301                      # Wait 5 minutes + 1 second
hexa jira fetch in-progress   # Should refresh
```

**Expected Output (After Sleep)**:
```
🔄 Récupération complète des tickets du sprint...
📋 Utilisation du cache (âge: 0s)
🔍 Recherche tickets: In Progress (filtre: all)

[... updated ticket list ...]

📊 Total: 3 ticket(s) en status 'In Progress' (filtre: all)
🔍 Cache: 47 tickets au total dans le sprint
```

**Validation**:
- ✅ Shows "Récupération complète" message (API call made)
- ✅ Cache age resets to 0s
- ✅ Ticket data reflects latest Jira state

---

## Scenario 4: "Me" Filter (First Time)

### Test: Filter by current user (user email not in config)

**Setup**:
```bash
# Ensure user email NOT in config
hexa config user get jira.userEmail  # Should return empty or error
```

**Command**:
```bash
hexa jira fetch to-test --filter=me
```

**Expected Output**:
```
🔄 Fetching user profile from Jira API...
✅ User email saved to config: john.doe@example.com
📋 Utilisation du cache (âge: 0s)
🔍 Recherche tickets: To test (filtre: me)

PROJ-130 - Fix login bug [John Doe] (High)
PROJ-132 - Update API docs [John Doe] (Medium)

📊 Total: 2 ticket(s) en status 'To test' (filtre: me)
🔍 Cache: 47 tickets au total dans le sprint
```

**Validation**:
- ✅ Shows "Fetching user profile" message
- ✅ Persists email to config (verify: `hexa config user get jira.userEmail`)
- ✅ Filters tickets where assignee email matches
- ✅ No API call on subsequent "me" filter runs

---

## Scenario 5: "Me" Filter (Subsequent Runs)

### Test: Reuse cached user email

**Command**:
```bash
hexa jira fetch blocked --filter=me
```

**Expected Output**:
```
📋 Utilisation du cache (âge: 15s)
🔍 Recherche tickets: Blocked (filtre: me)

PROJ-140 - Waiting for API access [John Doe] (High)

📊 Total: 1 ticket(s) en status 'Blocked' (filtre: me)
🔍 Cache: 47 tickets au total dans le sprint
```

**Validation**:
- ✅ No "Fetching user profile" message (uses config value)
- ✅ Correctly filters by user email from config
- ✅ Works across different status queries

---

## Scenario 6: Unassigned Filter

### Test: Show only tickets with no assignee

**Command**:
```bash
hexa jira fetch to-do --filter=unassigned
```

**Expected Output**:
```
📋 Utilisation du cache (âge: 45s)
🔍 Recherche tickets: To Do (filtre: unassigned)

PROJ-150 - Research new framework [Non assigné] (Medium)
PROJ-151 - Update changelog [Non assigné] (Low)

📊 Total: 2 ticket(s) en status 'To Do' (filtre: unassigned)
🔍 Cache: 47 tickets au total dans le sprint
```

**Validation**:
- ✅ Shows only tickets where assignee is null
- ✅ Displays "Non assigné" for all results
- ✅ Correct count of unassigned tickets

---

## Scenario 7: Cache Bypass Flag

### Test: Force fresh data fetch ignoring cache

**Command**:
```bash
hexa jira fetch in-progress              # Use cache
hexa jira fetch in-progress --no-cache   # Force refresh
```

**Expected Output (Second Command)**:
```
🔄 Récupération complète des tickets du sprint...
📋 Cache ignoré (--no-cache activé)
🔍 Recherche tickets: In Progress (filtre: all)

[... fresh ticket list ...]

📊 Total: 3 ticket(s) en status 'In Progress' (filtre: all)
🔍 Cache: 47 tickets au total dans le sprint
```

**Validation**:
- ✅ Shows "Cache ignoré" message
- ✅ API call made even if cache is fresh
- ✅ Cache file updated with new timestamp

---

## Scenario 8: Invalid Status Key

### Test: Error handling for invalid status

**Command**:
```bash
hexa jira fetch invalid-status
```

**Expected Output (STDERR)**:
```
Error: invalid status key 'invalid-status'

Valid status keys:
  - to-do
  - in-progress
  - to-test
  - uat
  - deploy-uat
  - to-deploy
  - blocked
  - prep
  - new
  - closed
  - archived

Usage: hexa jira fetch <status> [--filter=<me|unassigned|all>] [--no-cache]
```

**Validation**:
- ✅ Error message displayed to STDERR
- ✅ Lists all valid status keys
- ✅ Shows usage hint
- ✅ Exit code non-zero

---

## Scenario 9: No Tickets Found

### Test: Handle empty result set

**Command**:
```bash
# Assuming no tickets in "Archived" status
hexa jira fetch archived
```

**Expected Output**:
```
📋 Utilisation du cache (âge: 10s)
🔍 Recherche tickets: Archived (filtre: all)

Aucun ticket trouvé.

📊 Total: 0 ticket(s) en status 'Archived' (filtre: all)
🔍 Cache: 47 tickets au total dans le sprint
```

**Validation**:
- ✅ Shows "Aucun ticket trouvé" message
- ✅ Total count is 0
- ✅ No error (empty result is valid)

---

## Scenario 10: Jira API Authentication Failure

### Test: Error handling for invalid/missing token

**Setup**:
```bash
# Temporarily remove token
hexa config local set jira.token "invalid-token-12345"
```

**Command**:
```bash
hexa jira fetch in-progress
```

**Expected Output (STDERR)**:
```
Error: Jira API authentication failed (401 Unauthorized)

Please verify your Jira token:
  hexa config local get jira.token

To update your token:
  hexa config local set jira.token "your-valid-pat-here"
```

**Validation**:
- ✅ Clear error message about authentication
- ✅ Provides remediation steps
- ✅ Exit code non-zero
- ✅ No crash or panic

**Cleanup**:
```bash
# Restore valid token
hexa config local set jira.token "your-valid-pat-here"
```

---

## Scenario 11: Jira API Unreachable

### Test: Network error handling

**Setup**:
```bash
# Temporarily set invalid Jira URL
hexa config user set jira.url "https://nonexistent.atlassian.net"
```

**Command**:
```bash
hexa jira fetch in-progress
```

**Expected Output (STDERR)**:
```
Error: Failed to connect to Jira API
  URL: https://nonexistent.atlassian.net
  Reason: dial tcp: lookup nonexistent.atlassian.net: no such host

Please verify your Jira URL:
  hexa config user get jira.url
```

**Validation**:
- ✅ Network error displayed clearly
- ✅ Shows attempted URL
- ✅ Provides configuration check hint
- ✅ Exit code non-zero

**Cleanup**:
```bash
# Restore valid URL
hexa config user set jira.url "https://your-domain.atlassian.net"
```

---

## Scenario 12: Corrupted Cache File

### Test: Auto-recovery from corrupted cache

**Setup**:
```bash
# Find cache file
CACHE_FILE=$(ls ~/.hexa/cache/sprint_*.json | head -1)
# Corrupt it
echo "invalid json {{{" > $CACHE_FILE
```

**Command**:
```bash
hexa jira fetch in-progress
```

**Expected Output**:
```
⚠️  Cache file corrupted, refreshing...
🔄 Récupération complète des tickets du sprint...
📋 Utilisation du cache (âge: 0s)
🔍 Recherche tickets: In Progress (filtre: all)

[... fresh ticket list ...]

📊 Total: 3 ticket(s) en status 'In Progress' (filtre: all)
🔍 Cache: 47 tickets au total dans le sprint
```

**Validation**:
- ✅ Detects corrupted cache (JSON parse error)
- ✅ Shows warning message
- ✅ Fetches fresh data from API
- ✅ Overwrites corrupted cache with valid data
- ✅ No crash or data loss

---

## Success Criteria

All scenarios must pass with expected outputs. The feature is ready for release when:

- ✅ All 12 scenarios execute without crashes
- ✅ Cache TTL logic works (5-minute window)
- ✅ User email auto-fetch and persistence works
- ✅ All filters (me, unassigned, all) produce correct results
- ✅ Error messages are clear and actionable
- ✅ Output format matches bash script behavior (ticket key, summary, assignee, priority)
- ✅ Cache statistics displayed correctly

---

## Performance Benchmarks

Run these commands to validate performance goals (from Technical Context in plan.md):

```bash
# Cache hit performance (<50ms)
time hexa jira fetch in-progress  # Second run, cache hit
# Expected: real < 0.050s

# API call performance (<2s p95)
time hexa jira fetch in-progress --no-cache
# Expected: real < 2.000s (p95)

# Startup overhead (<100ms)
time hexa --version
# Expected: real < 0.100s
```

**Pass Criteria**: 95% of runs meet performance goals.
