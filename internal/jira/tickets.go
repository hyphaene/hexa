package jira

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"

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

	_, total, err := GetTickets(client, sprintID, 0, 1, 0)

	if err != nil {
		return nil, 0, fmt.Errorf("fetching tickets: %w", err)
	}
	totalCount := total

	fmt.Print("🚀 Starting to fetch tickets...\n")
	fmt.Printf("ℹ️  Total tickets to fetch: %d\n", totalCount)

	const maxResults = 25 // Adjust based on Jira API limits and performance

	pageRequests := CalculatePageRequests(totalCount, maxResults)

	mu := sync.Mutex{}
	var wg sync.WaitGroup

	for _, req := range pageRequests {
		wg.Add(1)
		go func(req PageRequest) {
			defer wg.Done()
			lastResults := min(req.StartAt+req.MaxResults, totalCount)
			fmt.Printf("🌐 [API Call %d] Fetching tickets %d to %d...\n", req.PageNum, req.StartAt+1, lastResults)
			sprintResp, _, err := GetTickets(client, sprintID, req.StartAt, req.MaxResults, req.PageNum)
			if err != nil {
				fmt.Fprintf(os.Stderr, "❌ [API Call %d] Error fetching tickets: %v\n", req.PageNum, err)
				return
			}

			// Append results in a thread-safe way
			mu.Lock()
			allTickets = append(allTickets, sprintResp.Issues...)
			mu.Unlock()

			fmt.Fprintf(os.Stderr, "✅ [API Call %d] Received %d tickets (total: %d/%d)\n",
				req.PageNum, len(sprintResp.Issues), len(allTickets), totalCount)

		}(req)
	}
	wg.Wait()

	// var bar *progressbar.ProgressBar

	// for page := 1; ; page++ {
	// 	sprintResp, total, err := GetTickets(client, sprintID, startAt, maxResults, page)
	// 	if err != nil {
	// 		return nil, 0, fmt.Errorf("fetching tickets: %w", err)
	// 	}
	// 	totalCount = total

	// 	// Initialize progress bar after first response
	// 	if bar == nil && sprintResp.Total > 0 {
	// 		bar = progressbar.NewOptions(sprintResp.Total,
	// 			progressbar.OptionSetWriter(os.Stderr),
	// 			progressbar.OptionSetDescription("📥 Fetching tickets"),
	// 			progressbar.OptionShowCount(),
	// 			progressbar.OptionSetWidth(40),
	// 			progressbar.OptionThrottle(65*time.Millisecond),
	// 			progressbar.OptionShowIts(),
	// 			progressbar.OptionSetItsString("tickets"),
	// 			progressbar.OptionOnCompletion(func() {
	// 				fmt.Fprint(os.Stderr, "\n")
	// 			}),
	// 		)
	// 	}

	// 	// Update progress
	// 	if bar != nil {
	// 		_ = bar.Add(len(sprintResp.Issues))
	// 	}

	// 	allTickets = append(allTickets, sprintResp.Issues...)
	// 	totalCount = sprintResp.Total

	// 	fmt.Fprintf(os.Stderr, "✅ [Page %d] Received %d tickets (total: %d/%d)\n",
	// 		page, len(sprintResp.Issues), len(allTickets), totalCount)

	// 	// Stop if no more tickets or explicitly marked as last page
	// 	if sprintResp.IsLast || len(sprintResp.Issues) == 0 {
	// 		break
	// 	}

	// 	startAt += maxResults
	// }

	// // Ensure progress bar completes
	// if bar != nil {
	// 	_ = bar.Finish()
	// }

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

func GetJiraFetchUrl(sprintID int, startAt int, maxResults int) string {
	jiraURL := viper.GetString("jira.url")
	return fmt.Sprintf("%s/rest/agile/1.0/sprint/%d/issue?startAt=%d&maxResults=%d",
		jiraURL, sprintID, startAt, maxResults)
}

// func getTicketsCount(sprintID int) (int, error) {
// 	url := jira.GetJiraFetchUrl(sprintID, 0, 1)

// }

func GetTickets(client *http.Client, sprintID int, startAt int, maxResults int, page int) (SprintIssuesResponse, int, error) {
	jiraToken := viper.GetString("jira.token")
	noData := SprintIssuesResponse{}

	if jiraToken == "" {
		return noData, 0, fmt.Errorf("jira.token must be configured")
	}
	url := GetJiraFetchUrl(sprintID, startAt, maxResults)

	// Show API call info
	fmt.Fprintf(os.Stderr, "🌐 [API Call %d] GET %s\n", page, url)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return noData, 0, fmt.Errorf("creating request: %w", err)
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", "Bearer "+jiraToken)

	resp, err := client.Do(req)
	if err != nil {
		return SprintIssuesResponse{}, 0, fmt.Errorf("calling Jira API: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		_ = resp.Body.Close()
		return noData, 0, fmt.Errorf("jira API returned status %d", resp.StatusCode)
	}

	var sprintResp SprintIssuesResponse
	if err := json.NewDecoder(resp.Body).Decode(&sprintResp); err != nil {
		_ = resp.Body.Close()
		return noData, 0, fmt.Errorf("decoding response: %w", err)
	}
	_ = resp.Body.Close()
	return sprintResp, sprintResp.Total, nil
}

type PageRequest struct {
	PageNum    int
	StartAt    int
	MaxResults int
}

func CalculatePageRequests(totalCount int, maxResults int) []PageRequest {
	// round up
	totalPages := (totalCount + maxResults - 1) / maxResults
	requests := make([]PageRequest, totalPages)
	for i := range totalPages {
		requests[i] = PageRequest{
			PageNum:    i + 1,
			StartAt:    i * maxResults,
			MaxResults: maxResults,
		}
	}

	return requests
}
