# Feature Specification: Jira Ticket Fetch Command

**Feature Branch**: `002-j-ai-besoin`
**Created**: 2025-10-01
**Status**: Draft
**Input**: User description: "j'ai besoin de transférer le script'/Users/maximilien/Hexactitude/claude/scripts/jira/jira_fetch.sh' en tant que commande de mon cli hexa"

## Execution Flow (main)
```
1. Parse user description from Input
   → Migrate bash script functionality to native Go CLI command
2. Extract key concepts from description
   → Actors: CLI users, Jira API
   → Actions: fetch tickets by status, filter by assignee, cache results
   → Data: Jira tickets, sprint information, user assignments
   → Constraints: API authentication, caching strategy (5min TTL)
3. For each unclear aspect:
   → All clarified (see Clarifications section)
4. Fill User Scenarios & Testing section
   → See below
5. Generate Functional Requirements
   → See below
6. Identify Key Entities
   → See below
7. Run Review Checklist
   → SUCCESS (all clarifications resolved)
8. Return: SUCCESS (spec ready for planning)
```

---

## ⚡ Quick Guidelines
- ✅ Focus on WHAT users need and WHY
- ❌ Avoid HOW to implement (no tech stack, APIs, code structure)
- 👥 Written for business stakeholders, not developers

---

## Clarifications

### Session 2025-10-01
- Q: Sprint identification - How should the system determine which sprint is "current"? → A: Use existing internal function from `/Users/maximilien/Code/hexa/internal/jira/getCurrentSprintId.go`
- Q: Authentication configuration - Where should Jira credentials (PAT token) be sourced from? → A: Use Viper config with key `HEXA_JIRA_TOKEN`
- Q: Cache behavior - Should the system cache full ticket list to avoid repeated API calls? → A: Yes, keep cache with addition of bypass flag for force refresh
- Q: Cache storage location - Where should the cached sprint ticket data be stored? → A: User config directory (`~/.hexa/cache/sprint_{id}.json`)
- Q: User identification for "me" filter - How should the system identify the current user's email for filtering "me" tickets? → A: Derive from Jira API using PAT token, then persist to Viper config file (user/project/project.local level) for long-term cache

---

## User Scenarios & Testing *(mandatory)*

### Primary User Story
As a developer working with Jira tickets, I need to quickly see tickets from the current sprint filtered by status and assignee without leaving my terminal, so I can prioritize my work and understand team progress.

### Acceptance Scenarios
1. **Given** I run the fetch command with status "in-progress", **When** the command executes, **Then** I see all tickets in "In Progress" status from the current sprint with their key, summary, assignee, and priority
2. **Given** I run the fetch command with status "to-test" and filter "me", **When** the command executes, **Then** I see only tickets assigned to me in "To test" status
2a. **Given** I run the fetch command with filter "me" for the first time, **When** user email is not in config, **Then** the system fetches authenticated user profile from Jira API and persists email to Viper config for future use
3. **Given** I run the fetch command twice within 5 minutes, **When** the second execution happens, **Then** the results come from cache instead of making a new API call
4. **Given** I run the fetch command with filter "unassigned", **When** the command executes, **Then** I see only tickets with no assignee in the specified status
5. **Given** cache is older than 5 minutes, **When** I run the fetch command, **Then** the system automatically refreshes data from Jira API
6. **Given** cache exists and is fresh, **When** I run the fetch command with bypass flag, **Then** the system ignores cache and fetches fresh data from Jira API

### Edge Cases
- What happens when Jira API authentication fails? System MUST display clear error message about missing or invalid credentials
- What happens when the sprint has no tickets matching the filter? System MUST display "No tickets found" message with filter context
- What happens when Jira API is unreachable? System MUST display error message and not crash
- What happens when an invalid status key is provided? System MUST list available status options
- What happens when cache file is corrupted? System MUST refresh from API automatically

## Requirements *(mandatory)*

### Functional Requirements
- **FR-001**: System MUST fetch tickets from the current active sprint only
- **FR-002**: System MUST support filtering by status using these status keys: to-do, in-progress, to-test, uat, deploy-uat, to-deploy, blocked, prep, new, closed, archived
- **FR-003**: System MUST support filtering by assignee with options: me (current user), unassigned, all (default)
- **FR-003a**: System MUST fetch authenticated user email from Jira API when "me" filter is used and email is not in config
- **FR-003b**: System MUST persist fetched user email to Viper config file at user, project, or project.local level for long-term reuse
- **FR-004**: System MUST display for each ticket: key, summary, assignee name (or "Non assigné"), and priority
- **FR-005**: System MUST cache sprint ticket data with 5-minute time-to-live
- **FR-006**: System MUST automatically refresh cache when expired or missing
- **FR-007**: System MUST provide a flag to bypass cache and force fresh data fetch from Jira API
- **FR-008**: System MUST display total count of tickets matching the filters
- **FR-009**: System MUST display cache statistics (total tickets in sprint)
- **FR-010**: System MUST authenticate to Jira using PAT token from Viper config key `HEXA_JIRA_TOKEN`
- **FR-011**: System MUST identify current sprint using existing internal getCurrentSprintId functionality
- **FR-012**: System MUST validate status key before making API call and display available options on invalid input
- **FR-013**: Users MUST be able to see cache age in output to understand data freshness
- **FR-014**: System MUST handle pagination when sprint contains more tickets than API page limit

### Key Entities *(include if feature involves data)*
- **Sprint**: Represents a Jira sprint with ID, contains multiple tickets
- **Ticket**: Jira issue with key, summary, status, assignee (optional), priority, belongs to one sprint
- **Status**: Enumeration of allowed workflow states (To Do, In Progress, To test, etc.)
- **Assignee**: User assigned to ticket with display name and email address
- **Cache**: Temporary storage of sprint tickets with timestamp, stored in user config directory at `~/.hexa/cache/sprint_{id}.json` with 5-minute TTL
- **User Profile**: Authenticated Jira user with email address, fetched from API on first "me" filter use, persisted to Viper config for long-term reuse

---

## Review & Acceptance Checklist
*GATE: Automated checks run during main() execution*

### Content Quality
- [x] No implementation details (languages, frameworks, APIs)
- [x] Focused on user value and business needs
- [x] Written for non-technical stakeholders
- [x] All mandatory sections completed

### Requirement Completeness
- [x] No [NEEDS CLARIFICATION] markers remain
- [x] Requirements are testable and unambiguous
- [x] Success criteria are measurable
- [x] Scope is clearly bounded
- [x] Dependencies and assumptions identified

---

## Execution Status
*Updated by main() during processing*

- [x] User description parsed
- [x] Key concepts extracted
- [x] Ambiguities marked
- [x] User scenarios defined
- [x] Requirements generated
- [x] Entities identified
- [x] Review checklist passed

---
