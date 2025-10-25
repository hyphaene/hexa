package jira

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sort"
	"sync"
	"time"

	"github.com/schollz/progressbar/v3"
	"github.com/spf13/viper"
	"golang.org/x/sync/errgroup"
)

// SprintIssuesResponse represents the API response for sprint tickets
type SprintIssuesResponse struct {
	MaxResults int      `json:"maxResults"`
	StartAt    int      `json:"startAt"`
	Total      int      `json:"total"`
	IsLast     bool     `json:"isLast"`
	Issues     []Ticket `json:"issues"`
}

// HTTPError represents an HTTP error with status code
type HTTPError struct {
	StatusCode int
	Message    string
}

func (e *HTTPError) Error() string {
	return e.Message
}

// isRetriable checks if an error is worth retrying
func isRetriable(err error) bool {
	if err == nil {
		return false
	}

	// Check if it's an HTTPError
	var httpErr *HTTPError
	if errors, ok := err.(*HTTPError); ok {
		httpErr = errors
	} else {
		// For wrapped errors, try to extract HTTPError
		for e := err; e != nil; {
			if he, ok := e.(*HTTPError); ok {
				httpErr = he
				break
			}
			// Try to unwrap
			if unwrapper, ok := e.(interface{ Unwrap() error }); ok {
				e = unwrapper.Unwrap()
			} else {
				break
			}
		}
	}

	if httpErr != nil {
		// Non-retriable status codes
		switch httpErr.StatusCode {
		case 400, 401, 403, 404, 405, 406, 410: // Client errors
			return false
		case 429: // Rate limit - retriable
			return true
		case 500, 502, 503, 504: // Server errors - retriable
			return true
		default:
			// Other 4xx are non-retriable, 5xx are retriable
			if httpErr.StatusCode >= 400 && httpErr.StatusCode < 500 {
				return false
			}
			return true
		}
	}

	// Network errors, timeouts, etc. are retriable
	return true
}

// FetchSprintTickets fetches all tickets from a sprint with pagination
func FetchSprintTickets(sprintID int, verbose bool) ([]Ticket, int, error) {
	jiraToken := viper.GetString("jira.token")
	jiraURL := viper.GetString("jira.url")

	if jiraToken == "" || jiraURL == "" {
		return nil, 0, fmt.Errorf("jira.token and jira.url must be configured")
	}

	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: viper.GetDuration("jira.sprint.timeout"),
	}

	var allTickets []Ticket

	maxResults := viper.GetInt("jira.sprint.maxResults")
	maxRetries := viper.GetInt("jira.sprint.maxRetries")
	retryDelay := viper.GetDuration("jira.sprint.retryDelay")

	var totalCount int
	var pageRequests []PageRequest

	// Optimize first call: fetch full batch if maxResults <= 50
	if maxResults <= 50 {
		// Fetch first batch and use it
		firstBatch, total, err := GetTickets(context.Background(), client, sprintID, 0, maxResults, verbose)
		if err != nil {
			return nil, 0, fmt.Errorf("fetching tickets: %w", err)
		}
		allTickets = append(allTickets, firstBatch.Issues...)
		totalCount = total

		// Calculate remaining pages (skip page 1)
		allPageRequests := CalculatePageRequests(totalCount, maxResults)
		if len(allPageRequests) > 1 {
			pageRequests = allPageRequests[1:]
		}
	} else {
		// Fallback: fetch 1 ticket just for total count
		_, total, err := GetTickets(context.Background(), client, sprintID, 0, 1, verbose)
		if err != nil {
			return nil, 0, fmt.Errorf("fetching tickets: %w", err)
		}
		totalCount = total
		pageRequests = CalculatePageRequests(totalCount, maxResults)
	}

	// Initialize progress bar or verbose logs
	var bar *progressbar.ProgressBar
	if !verbose {
		bar = progressbar.NewOptions(totalCount,
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
		// Ensure progress bar is finalized even on error
		defer func() {
			if bar != nil {
				_ = bar.Finish()
			}
		}()
		// Update progress bar with already fetched tickets
		if len(allTickets) > 0 {
			_ = bar.Add(len(allTickets))
		}
	} else {
		fmt.Fprint(os.Stderr, "🚀 Starting to fetch tickets...\n")
		fmt.Fprintf(os.Stderr, "ℹ️  Total tickets to fetch: %d\n", totalCount)
		if len(allTickets) > 0 {
			fmt.Fprintf(os.Stderr, "✅ [API Call 1] Received %d tickets (total: %d/%d)\n",
				len(allTickets), len(allTickets), totalCount)
		}
	}

	var mu sync.Mutex
	g, ctx := errgroup.WithContext(context.Background())

	for _, req := range pageRequests {
		g.Go(func() error {
			lastResults := min(req.StartAt+req.MaxResults, totalCount)

			var sprintResp SprintIssuesResponse
			var err error

			// Retry loop with linear backoff
			for attempt := 0; attempt < maxRetries; attempt++ {
				// Check if context cancelled (another goroutine failed)
				if ctx.Err() != nil {
					return ctx.Err()
				}

				if verbose && attempt > 0 {
					fmt.Fprintf(os.Stderr, "🔄 [API Call %d] Retry %d/%d...\n", req.PageNum, attempt, maxRetries-1)
				}

				if verbose && attempt == 0 {
					fmt.Fprintf(os.Stderr, "🌐 [API Call %d] Fetching tickets %d to %d...\n", req.PageNum, req.StartAt+1, lastResults)
				}

				sprintResp, _, err = GetTickets(ctx, client, sprintID, req.StartAt, req.MaxResults, verbose)
				if err == nil {
					// Success
					break
				}

				// Check if error is retriable
				if !isRetriable(err) {
					fmt.Fprintf(os.Stderr, "❌ [API Call %d] Non-retriable error: %v\n", req.PageNum, err)
					return fmt.Errorf("page %d (tickets %d-%d) non-retriable error: %w",
						req.PageNum, req.StartAt+1, lastResults, err)
				}

				// If not last attempt, wait before retry
				if attempt < maxRetries-1 {
					if verbose {
						fmt.Fprintf(os.Stderr, "⚠️  [API Call %d] Attempt %d failed: %v. Retrying in %s...\n",
							req.PageNum, attempt+1, err, retryDelay)
					}
					time.Sleep(retryDelay)
				}
			}

			// If still error after all retries, fail fast
			if err != nil {
				fmt.Fprintf(os.Stderr, "❌ [API Call %d] All retry attempts exhausted\n", req.PageNum)
				return fmt.Errorf("page %d (tickets %d-%d) failed after %d attempts: %w",
					req.PageNum, req.StartAt+1, lastResults, maxRetries, err)
			}

			// Append results and update progress in a thread-safe way
			mu.Lock()
			allTickets = append(allTickets, sprintResp.Issues...)
			if bar != nil {
				_ = bar.Add(len(sprintResp.Issues))
			} else if verbose {
				fmt.Fprintf(os.Stderr, "✅ [API Call %d] Received %d tickets (total: %d/%d)\n",
					req.PageNum, len(sprintResp.Issues), len(allTickets), totalCount)
			}
			mu.Unlock()

			return nil
		})
	}

	// Wait for all goroutines, fail-fast on first error
	if err := g.Wait(); err != nil {
		return nil, 0, err
	}

	// Sort tickets by key for deterministic order
	sort.Slice(allTickets, func(i, j int) bool {
		return allTickets[i].Key < allTickets[j].Key
	})

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

func GetTickets(ctx context.Context, client *http.Client, sprintID int, startAt int, maxResults int, verbose bool) (SprintIssuesResponse, int, error) {
	jiraToken := viper.GetString("jira.token")
	noData := SprintIssuesResponse{}

	if jiraToken == "" {
		return noData, 0, fmt.Errorf("jira.token must be configured")
	}
	url := GetJiraFetchUrl(sprintID, startAt, maxResults)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
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
		return noData, 0, &HTTPError{
			StatusCode: resp.StatusCode,
			Message:    fmt.Sprintf("jira API returned status %d", resp.StatusCode),
		}
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
