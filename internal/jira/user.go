package jira

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/spf13/viper"
)

// UserProfile represents the authenticated Jira user
type UserProfile struct {
	AccountID    string `json:"accountId"`    // Unique Jira account ID
	EmailAddress string `json:"emailAddress"` // User's email (used for "me" filter)
	DisplayName  string `json:"displayName"`  // Full name for display
}

// FetchCurrentUser fetches the authenticated user's profile from Jira API
func FetchCurrentUser() (*UserProfile, error) {
	jiraToken := viper.GetString("jira.token")
	jiraURL := viper.GetString("jira.url")

	if jiraToken == "" || jiraURL == "" {
		return nil, fmt.Errorf("jira.token and jira.url must be configured")
	}

	url := fmt.Sprintf("%s/rest/api/3/myself", jiraURL)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", "Bearer "+jiraToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("calling Jira API: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("jira API returned status %d", resp.StatusCode)
	}

	var profile UserProfile
	if err := json.NewDecoder(resp.Body).Decode(&profile); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}

	return &profile, nil
}

// SaveUserEmail persists the user's email to Viper config
func SaveUserEmail(email string) error {
	viper.Set("jira.userEmail", email)

	// Write to config file if one is being used
	if viper.ConfigFileUsed() != "" {
		if err := viper.WriteConfig(); err != nil {
			return fmt.Errorf("saving config: %w", err)
		}
	}

	return nil
}
