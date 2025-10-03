package jira

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/schollz/progressbar/v3"
	"github.com/spf13/viper"
)

// SprintIssuesResponse represents the API response for sprint tickets
type SprintIssuesResponse struct {
	MaxResults int      `json:"maxResults"`
	StartAt    int      `json:"startAt"`
	Total      int      `json:"total"`
	IsLast     bool     `json:"isLast"`
	Issues     []Ticket `json:"issues"`
}

// FetchSprintTickets fetches all tickets from a sprint with pagination
func FetchSprintTickets(sprintID int) ([]Ticket, int, error) {
	jiraToken := viper.GetString("jira.token")
	jiraURL := viper.GetString("jira.url")

	if jiraToken == "" || jiraURL == "" {
		return nil, 0, fmt.Errorf("jira.token and jira.url must be configured")
	}

	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	var allTickets []Ticket
	startAt := 0
	maxResults := 100
	totalCount := 0

	var bar *progressbar.ProgressBar

	for page := 1; ; page++ {
		url := fmt.Sprintf("%s/rest/agile/1.0/sprint/%d/issue?startAt=%d&maxResults=%d",
			jiraURL, sprintID, startAt, maxResults)

		// Show API call info
		fmt.Fprintf(os.Stderr, "🌐 [API Call %d] GET %s\n", page, url)

		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return nil, 0, fmt.Errorf("creating request: %w", err)
		}

		req.Header.Add("Accept", "application/json")
		req.Header.Add("Authorization", "Bearer "+jiraToken)

		resp, err := client.Do(req)
		if err != nil {
			return nil, 0, fmt.Errorf("calling Jira API: %w", err)
		}

		if resp.StatusCode != http.StatusOK {
			_ = resp.Body.Close()
			return nil, 0, fmt.Errorf("jira API returned status %d", resp.StatusCode)
		}

		var sprintResp SprintIssuesResponse
		if err := json.NewDecoder(resp.Body).Decode(&sprintResp); err != nil {
			_ = resp.Body.Close()
			return nil, 0, fmt.Errorf("decoding response: %w", err)
		}
		_ = resp.Body.Close()

		// Initialize progress bar after first response
		if bar == nil && sprintResp.Total > 0 {
			bar = progressbar.NewOptions(sprintResp.Total,
				progressbar.OptionSetWriter(os.Stderr),
				progressbar.OptionSetDescription("📥 Fetching tickets"),
				progressbar.OptionShowCount(),
				progressbar.OptionSetWidth(40),
				progressbar.OptionThrottle(65*time.Millisecond),
				progressbar.OptionShowIts(),
				progressbar.OptionSetItsString("tickets"),
				progressbar.OptionOnCompletion(func() {
					fmt.Fprint(os.Stderr, "\n")
				}),
			)
		}

		// Update progress
		if bar != nil {
			_ = bar.Add(len(sprintResp.Issues))
		}

		allTickets = append(allTickets, sprintResp.Issues...)
		totalCount = sprintResp.Total

		fmt.Fprintf(os.Stderr, "✅ [Page %d] Received %d tickets (total: %d/%d)\n",
			page, len(sprintResp.Issues), len(allTickets), totalCount)

		// Stop if no more tickets or explicitly marked as last page
		if sprintResp.IsLast || len(sprintResp.Issues) == 0 {
			break
		}

		startAt += maxResults
	}

	// Ensure progress bar completes
	if bar != nil {
		_ = bar.Finish()
	}

	return allTickets, totalCount, nil
}

// FilterByStatus filters tickets by exact status name
func FilterByStatus(tickets []Ticket, statusName string) []Ticket {
	var filtered []Ticket
	for _, ticket := range tickets {
		if ticket.Fields.Status.Name == statusName {
			filtered = append(filtered, ticket)
		}
	}
	return filtered
}

// FilterByAssignee filters tickets by assignee
// filter can be: "me" (match userEmail), "unassigned" (nil assignee), "all" (no filtering)
func FilterByAssignee(tickets []Ticket, filter string, userEmail string) []Ticket {
	if filter == "all" {
		return tickets
	}

	var filtered []Ticket
	for _, ticket := range tickets {
		switch filter {
		case "me":
			if ticket.Fields.Assignee != nil && ticket.Fields.Assignee.EmailAddress == userEmail {
				filtered = append(filtered, ticket)
			}
		case "unassigned":
			if ticket.Fields.Assignee == nil {
				filtered = append(filtered, ticket)
			}
		}
	}
	return filtered
}
