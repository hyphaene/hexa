package jira

// Ticket represents a single Jira issue with relevant fields for display and filtering
type Ticket struct {
	Key    string `json:"key"` // e.g., "PROJ-123"
	Fields Fields `json:"fields"`
}

// Fields contains the nested field structure from Jira API response
type Fields struct {
	Summary  string    `json:"summary"`
	Status   Status    `json:"status"`
	Assignee *Assignee `json:"assignee"` // Pointer: null when unassigned
	Priority *Priority `json:"priority"` // Pointer: null when no priority set
}

// Status represents the workflow status of a ticket
type Status struct {
	Name string `json:"name"` // e.g., "In Progress", "To test"
}

// Assignee represents the user assigned to a ticket
type Assignee struct {
	DisplayName  string `json:"displayName"`  // e.g., "John Doe"
	EmailAddress string `json:"emailAddress"` // e.g., "john.doe@example.com"
}

// Priority represents the priority level of a ticket
type Priority struct {
	Name string `json:"name"` // e.g., "High", "Medium", "Low"
}
