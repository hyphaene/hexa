package jira

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/spf13/viper"
)

// Test helpers - create fixtures for testing

// setupViperForTest configures viper for testing and returns cleanup function
func setupViperForTest(t *testing.T, jiraURL, jiraToken string) func() {
	t.Helper()

	// Save original values
	origURL := viper.Get("jira.url")
	origToken := viper.Get("jira.token")
	origMaxResults := viper.Get("jira.sprint.maxResults")
	origMaxRetries := viper.Get("jira.sprint.maxRetries")
	origRetryDelay := viper.Get("jira.sprint.retryDelay")

	// Set test values
	viper.Set("jira.url", jiraURL)
	viper.Set("jira.token", jiraToken)
	viper.Set("jira.sprint.maxResults", 50)
	viper.Set("jira.sprint.maxRetries", 1)
	viper.Set("jira.sprint.retryDelay", 100*time.Millisecond)
	viper.Set("jira.sprint.timeout", 30*time.Second)

	// Return cleanup function
	return func() {
		if origURL != nil {
			viper.Set("jira.url", origURL)
		}
		if origToken != nil {
			viper.Set("jira.token", origToken)
		}
		if origMaxResults != nil {
			viper.Set("jira.sprint.maxResults", origMaxResults)
		}
		if origMaxRetries != nil {
			viper.Set("jira.sprint.maxRetries", origMaxRetries)
		}
		if origRetryDelay != nil {
			viper.Set("jira.sprint.retryDelay", origRetryDelay)
		}
	}
}

func makeTicket(key string, status string, assigneeEmail string) Ticket {
	ticket := Ticket{
		Key: key,
		Fields: Fields{
			Summary: "Test ticket " + key,
			Status:  Status{Name: status},
		},
	}

	if assigneeEmail != "" {
		ticket.Fields.Assignee = &Assignee{
			DisplayName:  "Test User",
			EmailAddress: assigneeEmail,
		}
	}

	return ticket
}

// Unit tests for FilterByStatus

func TestFilterByStatus(t *testing.T) {
	tests := []struct {
		name       string
		tickets    []Ticket
		statusName string
		want       int // expected count
	}{
		{
			name: "filter finds matching status",
			tickets: []Ticket{
				makeTicket("PROJ-1", "To Do", ""),
				makeTicket("PROJ-2", "In Progress", ""),
				makeTicket("PROJ-3", "To Do", ""),
			},
			statusName: "To Do",
			want:       2,
		},
		{
			name: "filter returns empty when no matches",
			tickets: []Ticket{
				makeTicket("PROJ-1", "To Do", ""),
				makeTicket("PROJ-2", "In Progress", ""),
			},
			statusName: "Done",
			want:       0,
		},
		{
			name:       "filter on empty slice returns empty",
			tickets:    []Ticket{},
			statusName: "To Do",
			want:       0,
		},
		{
			name: "status matching is case sensitive",
			tickets: []Ticket{
				makeTicket("PROJ-1", "To Do", ""),
				makeTicket("PROJ-2", "to do", ""),
			},
			statusName: "To Do",
			want:       1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FilterByStatus(tt.tickets, tt.statusName)
			if len(got) != tt.want {
				t.Errorf("FilterByStatus() returned %d tickets, want %d", len(got), tt.want)
			}

			// Verify all returned tickets have correct status
			for _, ticket := range got {
				if ticket.Fields.Status.Name != tt.statusName {
					t.Errorf("ticket %s has status %q, want %q",
						ticket.Key, ticket.Fields.Status.Name, tt.statusName)
				}
			}
		})
	}
}

// Unit tests for FilterByAssignee

func TestFilterByAssignee(t *testing.T) {
	tests := []struct {
		name      string
		tickets   []Ticket
		filter    string
		userEmail string
		want      int
	}{
		{
			name: "filter 'me' returns only my tickets",
			tickets: []Ticket{
				makeTicket("PROJ-1", "To Do", "me@example.com"),
				makeTicket("PROJ-2", "To Do", "other@example.com"),
				makeTicket("PROJ-3", "To Do", "me@example.com"),
			},
			filter:    "me",
			userEmail: "me@example.com",
			want:      2,
		},
		{
			name: "filter 'unassigned' returns only unassigned tickets",
			tickets: []Ticket{
				makeTicket("PROJ-1", "To Do", ""),
				makeTicket("PROJ-2", "To Do", "user@example.com"),
				makeTicket("PROJ-3", "To Do", ""),
			},
			filter:    "unassigned",
			userEmail: "",
			want:      2,
		},
		{
			name: "filter 'all' returns all tickets",
			tickets: []Ticket{
				makeTicket("PROJ-1", "To Do", "me@example.com"),
				makeTicket("PROJ-2", "To Do", ""),
				makeTicket("PROJ-3", "To Do", "other@example.com"),
			},
			filter:    "all",
			userEmail: "me@example.com",
			want:      3,
		},
		{
			name: "filter 'me' with no matching tickets",
			tickets: []Ticket{
				makeTicket("PROJ-1", "To Do", "other@example.com"),
				makeTicket("PROJ-2", "To Do", ""),
			},
			filter:    "me",
			userEmail: "me@example.com",
			want:      0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FilterByAssignee(tt.tickets, tt.filter, tt.userEmail)
			if len(got) != tt.want {
				t.Errorf("FilterByAssignee() returned %d tickets, want %d", len(got), tt.want)
			}

			// Verify filtering logic
			for _, ticket := range got {
				switch tt.filter {
				case "me":
					if ticket.Fields.Assignee == nil || ticket.Fields.Assignee.EmailAddress != tt.userEmail {
						t.Errorf("ticket %s should be assigned to %s", ticket.Key, tt.userEmail)
					}
				case "unassigned":
					if ticket.Fields.Assignee != nil {
						t.Errorf("ticket %s should be unassigned", ticket.Key)
					}
				}
			}
		})
	}
}

// Unit tests for CalculatePageRequests

func TestCalculatePageRequests(t *testing.T) {
	tests := []struct {
		name       string
		totalCount int
		maxResults int
		wantPages  int
		checkPages []struct {
			pageNum    int
			startAt    int
			maxResults int
		}
	}{
		{
			name:       "exact multiple of maxResults",
			totalCount: 100,
			maxResults: 50,
			wantPages:  2,
			checkPages: []struct {
				pageNum    int
				startAt    int
				maxResults int
			}{
				{pageNum: 1, startAt: 0, maxResults: 50},
				{pageNum: 2, startAt: 50, maxResults: 50},
			},
		},
		{
			name:       "not exact multiple - rounds up",
			totalCount: 105,
			maxResults: 50,
			wantPages:  3,
			checkPages: []struct {
				pageNum    int
				startAt    int
				maxResults int
			}{
				{pageNum: 1, startAt: 0, maxResults: 50},
				{pageNum: 2, startAt: 50, maxResults: 50},
				{pageNum: 3, startAt: 100, maxResults: 50},
			},
		},
		{
			name:       "single page",
			totalCount: 30,
			maxResults: 50,
			wantPages:  1,
			checkPages: []struct {
				pageNum    int
				startAt    int
				maxResults int
			}{
				{pageNum: 1, startAt: 0, maxResults: 50},
			},
		},
		{
			name:       "small maxResults creates many pages",
			totalCount: 10,
			maxResults: 3,
			wantPages:  4,
			checkPages: []struct {
				pageNum    int
				startAt    int
				maxResults int
			}{
				{pageNum: 1, startAt: 0, maxResults: 3},
				{pageNum: 2, startAt: 3, maxResults: 3},
				{pageNum: 3, startAt: 6, maxResults: 3},
				{pageNum: 4, startAt: 9, maxResults: 3},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CalculatePageRequests(tt.totalCount, tt.maxResults)

			if len(got) != tt.wantPages {
				t.Errorf("CalculatePageRequests() returned %d pages, want %d", len(got), tt.wantPages)
			}

			// Verify specific page calculations
			for _, check := range tt.checkPages {
				idx := check.pageNum - 1
				if idx >= len(got) {
					t.Errorf("page %d not found in results", check.pageNum)
					continue
				}

				page := got[idx]
				if page.PageNum != check.pageNum {
					t.Errorf("page[%d].PageNum = %d, want %d", idx, page.PageNum, check.pageNum)
				}
				if page.StartAt != check.startAt {
					t.Errorf("page[%d].StartAt = %d, want %d", idx, page.StartAt, check.startAt)
				}
				if page.MaxResults != check.maxResults {
					t.Errorf("page[%d].MaxResults = %d, want %d", idx, page.MaxResults, check.maxResults)
				}
			}
		})
	}
}

// Unit tests for GetJiraFetchUrl

func TestGetJiraFetchUrl(t *testing.T) {
	tests := []struct {
		name       string
		sprintID   int
		startAt    int
		maxResults int
		want       string
	}{
		{
			name:       "basic url construction",
			sprintID:   123,
			startAt:    0,
			maxResults: 50,
			want:       "https://jira.example.com/rest/agile/1.0/sprint/123/issue?startAt=0&maxResults=50",
		},
		{
			name:       "with offset",
			sprintID:   456,
			startAt:    100,
			maxResults: 50,
			want:       "https://jira.example.com/rest/agile/1.0/sprint/456/issue?startAt=100&maxResults=50",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cleanup := setupViperForTest(t, "https://jira.example.com", "test-token")
			defer cleanup()

			got := GetJiraFetchUrl(tt.sprintID, tt.startAt, tt.maxResults)
			if got != tt.want {
				t.Errorf("GetJiraFetchUrl() = %q, want %q", got, tt.want)
			}
		})
	}
}

// Integration tests with mock HTTP server

func TestGetTickets(t *testing.T) {
	tests := []struct {
		name           string
		mockResponse   SprintIssuesResponse
		mockStatusCode int
		wantErr        bool
		wantTotal      int
		wantIssueCount int
	}{
		{
			name: "successful fetch with tickets",
			mockResponse: SprintIssuesResponse{
				MaxResults: 50,
				StartAt:    0,
				Total:      2,
				IsLast:     true,
				Issues: []Ticket{
					makeTicket("PROJ-1", "To Do", "user@example.com"),
					makeTicket("PROJ-2", "In Progress", ""),
				},
			},
			mockStatusCode: http.StatusOK,
			wantErr:        false,
			wantTotal:      2,
			wantIssueCount: 2,
		},
		{
			name: "empty response",
			mockResponse: SprintIssuesResponse{
				MaxResults: 50,
				StartAt:    0,
				Total:      0,
				IsLast:     true,
				Issues:     []Ticket{},
			},
			mockStatusCode: http.StatusOK,
			wantErr:        false,
			wantTotal:      0,
			wantIssueCount: 0,
		},
		{
			name:           "API returns 401 unauthorized",
			mockResponse:   SprintIssuesResponse{},
			mockStatusCode: http.StatusUnauthorized,
			wantErr:        true,
		},
		{
			name:           "API returns 500 server error",
			mockResponse:   SprintIssuesResponse{},
			mockStatusCode: http.StatusInternalServerError,
			wantErr:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock HTTP server
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// Verify request
				if r.Method != "GET" {
					t.Errorf("expected GET request, got %s", r.Method)
				}

				if r.Header.Get("Accept") != "application/json" {
					t.Errorf("expected Accept: application/json header")
				}

				if r.Header.Get("Authorization") == "" {
					t.Errorf("expected Authorization header")
				}

				// Send mock response
				w.WriteHeader(tt.mockStatusCode)
				if tt.mockStatusCode == http.StatusOK {
					json.NewEncoder(w).Encode(tt.mockResponse)
				}
			}))
			defer server.Close()

			// Setup viper config
			cleanup := setupViperForTest(t, server.URL, "test-token")
			defer cleanup()

			// Create HTTP client
			client := &http.Client{}

			// Execute test
			got, total, err := GetTickets(context.Background(), client, 123, 0, 50, false)

			// Check error expectation
			if (err != nil) != tt.wantErr {
				t.Errorf("GetTickets() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				return // Don't check response on error cases
			}

			// Verify response
			if total != tt.wantTotal {
				t.Errorf("GetTickets() total = %d, want %d", total, tt.wantTotal)
			}

			if len(got.Issues) != tt.wantIssueCount {
				t.Errorf("GetTickets() returned %d issues, want %d", len(got.Issues), tt.wantIssueCount)
			}
		})
	}
}

// Test that demonstrates the sorting behavior

func TestFetchSprintTickets_Sorting(t *testing.T) {
	// Create mock server that returns tickets in random order
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := SprintIssuesResponse{
			MaxResults: 50,
			StartAt:    0,
			Total:      3,
			IsLast:     true,
			Issues: []Ticket{
				makeTicket("PROJ-3", "To Do", ""),
				makeTicket("PROJ-1", "To Do", ""),
				makeTicket("PROJ-2", "To Do", ""),
			},
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	// Setup viper config
	cleanup := setupViperForTest(t, server.URL, "test-token")
	defer cleanup()

	// Execute
	tickets, _, err := FetchSprintTickets(123, false)
	if err != nil {
		t.Fatalf("FetchSprintTickets() error = %v", err)
	}

	// Verify sorting - tickets should be ordered by key
	expectedOrder := []string{"PROJ-1", "PROJ-2", "PROJ-3"}
	for i, ticket := range tickets {
		if ticket.Key != expectedOrder[i] {
			t.Errorf("tickets[%d].Key = %s, want %s", i, ticket.Key, expectedOrder[i])
		}
	}
}
