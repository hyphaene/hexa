package sprint

import (
	"fmt"

	"github.com/hyphaene/hexa/internal/cache"
	"github.com/hyphaene/hexa/internal/jira"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var pulseCmd = &cobra.Command{
	Use:   "pulse",
	Short: "Sprint overview with key status categories",
	Long: `Display a comprehensive overview of the current sprint with tickets grouped by status.

Shows:
  - My TO DO tickets
  - My IN PROGRESS tickets
  - All DEPLOY IN UAT tickets
  - All BLOCKED tickets

This command fetches all sprint tickets once and filters in-memory for optimal performance.`,
	RunE: runPulse,
}

func init() {
	SprintCmd.AddCommand(pulseCmd)
}

func runPulse(cmd *cobra.Command, args []string) error {
	// Get current sprint ID
	sprintID, err := jira.GetCurrentSprintId()
	if err != nil {
		return fmt.Errorf("getting current sprint ID: %w", err)
	}

	// Check cache
	cachedEntry, err := cache.ReadCache(sprintID)
	if err != nil {
		cachedEntry = nil // Treat corrupted cache as cache miss
	}

	var tickets []jira.Ticket
	var total int

	// Determine if we need to refresh
	if cache.ShouldRefresh(cachedEntry, false) {
		_, _ = fmt.Fprintf(cmd.OutOrStdout(), "🔄 Récupération complète des tickets du sprint...\n")

		// Fetch from API
		fetchedTickets, fetchedTotal, err := jira.FetchSprintTickets(sprintID)
		if err != nil {
			return fmt.Errorf("fetching sprint tickets: %w", err)
		}

		// Write to cache
		if err := cache.WriteCache(sprintID, fetchedTickets, fetchedTotal); err != nil {
			_, _ = fmt.Fprintf(cmd.ErrOrStderr(), "Warning: failed to write cache: %v\n", err)
		}

		tickets = fetchedTickets
		total = fetchedTotal
	} else {
		// Use cached data
		tickets = cachedEntry.Issues
		total = cachedEntry.Total
		cacheAge := cachedEntry.Age()
		_, _ = fmt.Fprintf(cmd.OutOrStdout(), "📋 Utilisation du cache (âge: %s)\n", formatDuration(cacheAge))
	}

	// Get user email for "me" filters
	userEmail := viper.GetString("jira.userEmail")
	if userEmail == "" {
		_, _ = fmt.Fprintf(cmd.OutOrStdout(), "🔄 Fetching user profile from Jira API...\n")
		profile, err := jira.FetchCurrentUser()
		if err != nil {
			return fmt.Errorf("fetching user profile: %w", err)
		}

		userEmail = profile.EmailAddress

		// Save to config
		if err := jira.SaveUserEmail(userEmail); err != nil {
			_, _ = fmt.Fprintf(cmd.ErrOrStderr(), "Warning: failed to save user email to config: %v\n", err)
		}
	}

	// Filter by status and assignee in-memory
	myTodo := jira.FilterByStatus(tickets, "To Do")
	myTodo = jira.FilterByAssignee(myTodo, "me", userEmail)

	myInProgress := jira.FilterByStatus(tickets, "In Progress")
	myInProgress = jira.FilterByAssignee(myInProgress, "me", userEmail)

	allDeployUat := jira.FilterByStatus(tickets, "DEPLOY IN UAT")

	allBlocked := jira.FilterByStatus(tickets, "Blocked")

	// Display overview
	_, _ = fmt.Fprintf(cmd.OutOrStdout(), "\n📊 Sprint Pulse\n\n")

	_, _ = fmt.Fprintf(cmd.OutOrStdout(), "🔵 Mes TO DO: %d ticket(s)\n", len(myTodo))
	printTicketList(cmd, myTodo)

	_, _ = fmt.Fprintf(cmd.OutOrStdout(), "\n🟡 Mes IN PROGRESS: %d ticket(s)\n", len(myInProgress))
	printTicketList(cmd, myInProgress)

	_, _ = fmt.Fprintf(cmd.OutOrStdout(), "\n🟢 DEPLOY IN UAT: %d ticket(s)\n", len(allDeployUat))
	printTicketList(cmd, allDeployUat)

	_, _ = fmt.Fprintf(cmd.OutOrStdout(), "\n🔴 BLOCKED: %d ticket(s)\n", len(allBlocked))
	printTicketList(cmd, allBlocked)

	_, _ = fmt.Fprintf(cmd.OutOrStdout(), "\n🔍 Sprint total: %d tickets\n", total)

	return nil
}

func printTicketList(cmd *cobra.Command, tickets []jira.Ticket) {
	if len(tickets) == 0 {
		_, _ = fmt.Fprintf(cmd.OutOrStdout(), "  Aucun ticket.\n")
		return
	}

	for _, ticket := range tickets {
		assignee := "Non assigné"
		if ticket.Fields.Assignee != nil {
			assignee = ticket.Fields.Assignee.DisplayName
		}

		priority := "Medium"
		if ticket.Fields.Priority != nil {
			priority = ticket.Fields.Priority.Name
		}

		_, _ = fmt.Fprintf(cmd.OutOrStdout(), "  %s - %s [%s] (%s)\n",
			ticket.Key, ticket.Fields.Summary, assignee, priority)
	}
}
