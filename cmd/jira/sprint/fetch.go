package sprint

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/hyphaene/hexa/internal/cache"
	"github.com/hyphaene/hexa/internal/jira"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	filterFlag       string
	noCacheFlag      bool
	jsonFlag         bool
	verboseFlag      bool
	sprintNumberFlag int
	outputFlag       string
)

var fetchCmd = &cobra.Command{
	Use:   "fetch [status]",
	Short: "Fetch tickets from current sprint filtered by status",
	Long: `Fetch and display tickets from the current sprint filtered by status.

If no status is provided, fetches all tickets.

Status keys (CLI-friendly):
  to-do, in-progress, to-test, uat, deploy-uat, to-deploy,
  blocked, prep, new, closed, archived

Filter options:
  --filter=all        Show all tickets (default)
  --filter=me         Show only tickets assigned to you
  --filter=unassigned Show only unassigned tickets

Cache behavior:
  By default, ticket data is cached for 5 minutes.
  Use --no-cache to force a fresh fetch from Jira API.`,
	Args: cobra.MaximumNArgs(1),
	RunE: runFetch,
}

func init() {
	SprintCmd.AddCommand(fetchCmd)
	fetchCmd.Flags().StringVar(&filterFlag, "filter", "all", "Filter by assignee: me|unassigned|all")
	fetchCmd.Flags().BoolVar(&noCacheFlag, "no-cache", false, "Bypass cache and fetch fresh data")
	fetchCmd.Flags().BoolVar(&jsonFlag, "json", false, "Output results in JSON format")
	fetchCmd.Flags().BoolVarP(&verboseFlag, "verbose", "v", false, "Show detailed progress information")
	fetchCmd.Flags().IntVar(&sprintNumberFlag, "sprint-number", 0, "Fetch specific sprint by number (e.g., 35)")
	fetchCmd.Flags().StringVarP(&outputFlag, "output", "o", "", "Write output to file (markdown by default, JSON if --json)")
}

func runFetch(cmd *cobra.Command, args []string) error {
	// Verbose logging
	if verboseFlag && !jsonFlag {
		_, _ = fmt.Fprintf(cmd.ErrOrStderr(), "🔍 [DEBUG] Starting fetch command\n")
		_, _ = fmt.Fprintf(cmd.ErrOrStderr(), "🔍 [DEBUG] Jira URL: %s\n", viper.GetString("jira.url"))
		_, _ = fmt.Fprintf(cmd.ErrOrStderr(), "🔍 [DEBUG] Board ID: %d\n", viper.GetInt("jira.boardId"))
		if viper.GetString("jira.token") == "" {
			_, _ = fmt.Fprintf(cmd.ErrOrStderr(), "⚠️  [WARN] jira.token is not configured!\n")
		}
	}

	var statusName string
	filterByStatus := true

	// If no status arg provided, fetch all tickets
	if len(args) == 0 {
		filterByStatus = false
	} else {
		// Validate and map status key to Jira status name
		statusKey := args[0]
		var err error
		statusName, err = jira.MapStatusKey(statusKey)
		if err != nil {
			_, _ = fmt.Fprintf(cmd.ErrOrStderr(), "Error: %v\n\n", err)
			_, _ = fmt.Fprintf(cmd.ErrOrStderr(), "Valid status keys:\n")
			for _, key := range jira.ValidStatusKeys() {
				_, _ = fmt.Fprintf(cmd.ErrOrStderr(), "  - %s\n", key)
			}
			_, _ = fmt.Fprintf(cmd.ErrOrStderr(), "\nUsage: hexa jira sprint fetch [status] [--filter=<me|unassigned|all>] [--no-cache]\n")
			return fmt.Errorf("invalid status key")
		}
		if verboseFlag && !jsonFlag {
			_, _ = fmt.Fprintf(cmd.ErrOrStderr(), "🔍 [DEBUG] Status mapped: %s -> %s\n", statusKey, statusName)
		}
	}

	// Get sprint ID (current or specific number)
	var sprintID int
	if sprintNumberFlag > 0 {
		if verboseFlag && !jsonFlag {
			_, _ = fmt.Fprintf(cmd.ErrOrStderr(), "🔍 [DEBUG] Resolving sprint number %d...\n", sprintNumberFlag)
		}
		var err error
		sprintID, err = jira.GetSprintIdFromNumber(sprintNumberFlag)
		if err != nil {
			return fmt.Errorf("resolving sprint number %d: %w", sprintNumberFlag, err)
		}
		if verboseFlag && !jsonFlag {
			_, _ = fmt.Fprintf(cmd.ErrOrStderr(), "🔍 [DEBUG] Sprint ID: %d\n", sprintID)
		}
	} else {
		if verboseFlag && !jsonFlag {
			_, _ = fmt.Fprintf(cmd.ErrOrStderr(), "🔍 [DEBUG] Fetching current sprint ID...\n")
		}
		var err error
		sprintID, err = jira.GetCurrentSprintId()
		if err != nil {
			return fmt.Errorf("getting current sprint ID: %w", err)
		}
		if verboseFlag && !jsonFlag {
			_, _ = fmt.Fprintf(cmd.ErrOrStderr(), "🔍 [DEBUG] Sprint ID: %d\n", sprintID)
		}
	}

	// Check cache
	if verboseFlag && !jsonFlag {
		_, _ = fmt.Fprintf(cmd.ErrOrStderr(), "🔍 [DEBUG] Checking cache for sprint %d...\n", sprintID)
	}
	cachedEntry, err := cache.ReadCache(sprintID)
	if err != nil {
		if !jsonFlag {
			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "⚠️  Cache file corrupted, refreshing...\n")
		}
		cachedEntry = nil // Treat corrupted cache as cache miss
	} else if cachedEntry != nil && verboseFlag && !jsonFlag {
		_, _ = fmt.Fprintf(cmd.ErrOrStderr(), "🔍 [DEBUG] Cache found (age: %s)\n", formatDuration(cachedEntry.Age()))
	}

	var tickets []jira.Ticket
	var total int
	var cacheAge time.Duration

	// Determine if we need to refresh
	if cache.ShouldRefresh(cachedEntry, noCacheFlag) {
		if !jsonFlag {
			if noCacheFlag {
				_, _ = fmt.Fprintf(cmd.OutOrStdout(), "🔄 Récupération complète des tickets du sprint...\n")
			} else if cachedEntry == nil {
				if verboseFlag {
					_, _ = fmt.Fprintf(cmd.ErrOrStderr(), "🔍 [DEBUG] No cache found, fetching from API...\n")
				}
			} else {
				_, _ = fmt.Fprintf(cmd.OutOrStdout(), "🔄 Récupération complète des tickets du sprint...\n")
			}
		}

		// Fetch from API
		if verboseFlag && !jsonFlag {
			_, _ = fmt.Fprintf(cmd.ErrOrStderr(), "🔍 [DEBUG] Calling Jira API /rest/agile/1.0/sprint/%d/issue...\n", sprintID)
		}
		fetchedTickets, fetchedTotal, err := jira.FetchSprintTickets(sprintID)
		if err != nil {
			return handleAPIError(err, cmd)
		}
		if verboseFlag && !jsonFlag {
			_, _ = fmt.Fprintf(cmd.ErrOrStderr(), "🔍 [DEBUG] Received %d tickets from API\n", fetchedTotal)
		}

		// Write to cache
		if err := cache.WriteCache(sprintID, fetchedTickets, fetchedTotal); err != nil {
			// Non-fatal: log warning but continue
			if !jsonFlag {
				_, _ = fmt.Fprintf(cmd.ErrOrStderr(), "Warning: failed to write cache: %v\n", err)
			}
		} else if verboseFlag && !jsonFlag {
			_, _ = fmt.Fprintf(cmd.ErrOrStderr(), "🔍 [DEBUG] Cache written successfully\n")
		}

		tickets = fetchedTickets
		total = fetchedTotal
		cacheAge = 0
	} else {
		// Use cached data
		if verboseFlag && !jsonFlag {
			_, _ = fmt.Fprintf(cmd.ErrOrStderr(), "🔍 [DEBUG] Using cached data\n")
		}
		tickets = cachedEntry.Issues
		total = cachedEntry.Total
		cacheAge = cachedEntry.Age()
	}

	// Filter by status (only if status arg provided)
	if filterByStatus {
		tickets = jira.FilterByStatus(tickets, statusName)
	}

	// Filter by assignee
	switch filterFlag {
	case "me":
		userEmail := viper.GetString("jira.userEmail")
		if userEmail == "" {
			// Fetch user profile
			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "🔄 Fetching user profile from Jira API...\n")
			profile, err := jira.FetchCurrentUser()
			if err != nil {
				return fmt.Errorf("fetching user profile: %w", err)
			}

			userEmail = profile.EmailAddress

			// Save to config
			if err := jira.SaveUserEmail(userEmail); err != nil {
				// Non-fatal: log warning
				_, _ = fmt.Fprintf(cmd.ErrOrStderr(), "Warning: failed to save user email to config: %v\n", err)
			} else {
				_, _ = fmt.Fprintf(cmd.OutOrStdout(), "✅ User email saved to config: %s\n", userEmail)
			}
		}

		tickets = jira.FilterByAssignee(tickets, "me", userEmail)
	case "unassigned":
		tickets = jira.FilterByAssignee(tickets, "unassigned", "")
	}

	// Display output
	displayStatus := statusName
	if !filterByStatus {
		displayStatus = "tous statuts"
	}

	// Handle output flag
	if outputFlag != "" {
		return writeToFile(outputFlag, tickets, cacheAge, total, displayStatus, filterFlag, noCacheFlag, sprintID, jsonFlag)
	}

	if jsonFlag {
		return outputJSON(cmd, tickets, cacheAge, total, displayStatus, filterFlag, noCacheFlag, sprintID)
	}

	formatOutput(cmd, tickets, cacheAge, total, displayStatus, filterFlag, noCacheFlag)

	return nil
}

func formatOutput(cmd *cobra.Command, tickets []jira.Ticket, cacheAge time.Duration, total int, statusName string, filter string, noCache bool) {
	// Cache status
	if noCache {
		_, _ = fmt.Fprintf(cmd.OutOrStdout(), "📋 Cache ignoré (--no-cache activé)\n")
	} else {
		_, _ = fmt.Fprintf(cmd.OutOrStdout(), "📋 Utilisation du cache (âge: %s)\n", formatDuration(cacheAge))
	}

	_, _ = fmt.Fprintf(cmd.OutOrStdout(), "🔍 Recherche tickets: %s (filtre: %s)\n\n", statusName, filter)

	// Display tickets
	if len(tickets) == 0 {
		_, _ = fmt.Fprintf(cmd.OutOrStdout(), "Aucun ticket trouvé.\n\n")
	} else {
		for _, ticket := range tickets {
			assignee := "Non assigné"
			if ticket.Fields.Assignee != nil {
				assignee = ticket.Fields.Assignee.DisplayName
			}

			priority := "Medium"
			if ticket.Fields.Priority != nil {
				priority = ticket.Fields.Priority.Name
			}

			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "%s - %s [%s] (%s)\n",
				ticket.Key, ticket.Fields.Summary, assignee, priority)
		}
		_, _ = fmt.Fprintf(cmd.OutOrStdout(), "\n")
	}

	// Summary statistics
	_, _ = fmt.Fprintf(cmd.OutOrStdout(), "📊 Total: %d ticket(s) en status '%s' (filtre: %s)\n", len(tickets), statusName, filter)
	_, _ = fmt.Fprintf(cmd.OutOrStdout(), "🔍 Cache: %d tickets au total dans le sprint\n", total)
}

func formatDuration(d time.Duration) string {
	if d < time.Second {
		return "0s"
	}
	if d < time.Minute {
		return fmt.Sprintf("%.0fs", d.Seconds())
	}
	if d < time.Hour {
		minutes := int(d.Minutes())
		seconds := int(d.Seconds()) % 60
		if seconds == 0 {
			return fmt.Sprintf("%dm", minutes)
		}
		return fmt.Sprintf("%dm%ds", minutes, seconds)
	}
	return fmt.Sprintf("%.1fh", d.Hours())
}

func handleAPIError(err error, cmd *cobra.Command) error {
	errMsg := err.Error()

	if strings.Contains(errMsg, "status 401") {
		_, _ = fmt.Fprintf(cmd.ErrOrStderr(), "Error: Jira API authentication failed (401 Unauthorized)\n\n")
		_, _ = fmt.Fprintf(cmd.ErrOrStderr(), "Please verify your Jira token:\n")
		_, _ = fmt.Fprintf(cmd.ErrOrStderr(), "  hexa config local get jira.token\n\n")
		_, _ = fmt.Fprintf(cmd.ErrOrStderr(), "To update your token:\n")
		_, _ = fmt.Fprintf(cmd.ErrOrStderr(), "  hexa config local set jira.token \"your-valid-pat-here\"\n")
		return fmt.Errorf("authentication failed")
	}

	if strings.Contains(errMsg, "dial tcp") || strings.Contains(errMsg, "no such host") {
		_, _ = fmt.Fprintf(cmd.ErrOrStderr(), "Error: Failed to connect to Jira API\n")
		_, _ = fmt.Fprintf(cmd.ErrOrStderr(), "  Reason: %v\n\n", err)
		_, _ = fmt.Fprintf(cmd.ErrOrStderr(), "Please verify your Jira URL:\n")
		_, _ = fmt.Fprintf(cmd.ErrOrStderr(), "  hexa config user get jira.url\n")
		return fmt.Errorf("connection failed")
	}

	return fmt.Errorf("jira API error: %w", err)
}

// JSONOutput represents the JSON structure for output
type JSONOutput struct {
	Sprint struct {
		ID    int    `json:"id"`
		Total int    `json:"total"`
		Cache struct {
			Age     string `json:"age"`
			Expired bool   `json:"expired"`
		} `json:"cache"`
	} `json:"sprint"`
	Filter struct {
		Status   string `json:"status"`
		Assignee string `json:"assignee"`
	} `json:"filter"`
	Tickets []jira.Ticket `json:"tickets"`
	Summary struct {
		Count      int    `json:"count"`
		TotalCache int    `json:"total_cache"`
		CacheUsed  bool   `json:"cache_used"`
	} `json:"summary"`
}

func outputJSON(cmd *cobra.Command, tickets []jira.Ticket, cacheAge time.Duration, total int, statusName string, filter string, noCache bool, sprintID int) error {
	output := JSONOutput{}
	output.Sprint.ID = sprintID
	output.Sprint.Total = total
	output.Sprint.Cache.Age = formatDuration(cacheAge)
	output.Sprint.Cache.Expired = cacheAge > 5*time.Minute
	output.Filter.Status = statusName
	output.Filter.Assignee = filter
	output.Tickets = tickets
	output.Summary.Count = len(tickets)
	output.Summary.TotalCache = total
	output.Summary.CacheUsed = !noCache

	data, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		return fmt.Errorf("marshaling JSON: %w", err)
	}

	_, _ = fmt.Fprintf(cmd.OutOrStdout(), "%s\n", data)
	return nil
}

// writeToFile writes output to a file (markdown or JSON based on jsonFlag)
func writeToFile(filepath string, tickets []jira.Ticket, cacheAge time.Duration, total int, statusName string, filter string, noCache bool, sprintID int, asJSON bool) error {
	var content []byte
	var err error

	if asJSON {
		// JSON format
		output := JSONOutput{}
		output.Sprint.ID = sprintID
		output.Sprint.Total = total
		output.Sprint.Cache.Age = formatDuration(cacheAge)
		output.Sprint.Cache.Expired = cacheAge > 5*time.Minute
		output.Filter.Status = statusName
		output.Filter.Assignee = filter
		output.Tickets = tickets
		output.Summary.Count = len(tickets)
		output.Summary.TotalCache = total
		output.Summary.CacheUsed = !noCache

		content, err = json.MarshalIndent(output, "", "  ")
		if err != nil {
			return fmt.Errorf("marshaling JSON: %w", err)
		}
	} else {
		// Markdown format
		var md strings.Builder

		md.WriteString(fmt.Sprintf("# Sprint Report - %s\n\n", statusName))
		md.WriteString(fmt.Sprintf("**Sprint ID**: %d\n", sprintID))
		md.WriteString(fmt.Sprintf("**Filter**: %s\n", filter))
		md.WriteString(fmt.Sprintf("**Cache**: %s\n", formatDuration(cacheAge)))
		if noCache {
			md.WriteString("**Cache Status**: Bypassed (--no-cache)\n")
		}
		md.WriteString(fmt.Sprintf("\n## Summary\n\n"))
		md.WriteString(fmt.Sprintf("- **Total tickets in sprint**: %d\n", total))
		md.WriteString(fmt.Sprintf("- **Filtered tickets**: %d\n\n", len(tickets)))

		md.WriteString("## Tickets\n\n")

		if len(tickets) == 0 {
			md.WriteString("_No tickets found._\n")
		} else {
			for _, ticket := range tickets {
				assignee := "Non assigné"
				if ticket.Fields.Assignee != nil {
					assignee = ticket.Fields.Assignee.DisplayName
				}

				priority := "Medium"
				if ticket.Fields.Priority != nil {
					priority = ticket.Fields.Priority.Name
				}

				md.WriteString(fmt.Sprintf("### %s\n\n", ticket.Key))
				md.WriteString(fmt.Sprintf("**Summary**: %s\n\n", ticket.Fields.Summary))
				md.WriteString(fmt.Sprintf("- **Assignee**: %s\n", assignee))
				md.WriteString(fmt.Sprintf("- **Priority**: %s\n", priority))
				md.WriteString(fmt.Sprintf("- **Status**: %s\n\n", ticket.Fields.Status.Name))
			}
		}

		content = []byte(md.String())
	}

	// Write to file
	if err := os.WriteFile(filepath, content, 0644); err != nil {
		return fmt.Errorf("writing to file %s: %w", filepath, err)
	}

	fmt.Fprintf(os.Stderr, "✅ Output written to: %s\n", filepath)
	return nil
}
