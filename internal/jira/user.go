package jira

// UserProfile represents the authenticated Jira user
type UserProfile struct {
	AccountID    string `json:"accountId"`    // Unique Jira account ID
	EmailAddress string `json:"emailAddress"` // User's email (used for "me" filter)
	DisplayName  string `json:"displayName"`  // Full name for display
}
