package jira

import (
	"fmt"
	"sort"
)

// StatusMap provides CLI key → Jira status name mapping
var StatusMap = map[string]string{
	"to-do":       "To Do",
	"in-progress": "In Progress",
	"to-test":     "To test",
	"uat":         "UAT",
	"deploy-uat":  "DEPLOY IN UAT",
	"to-deploy":   "To deploy",
	"blocked":     "Blocked",
	"prep":        "Prep",
	"new":         "New",
	"closed":      "Closed",
	"archived":    "Archived",
}

// ValidStatusKeys returns all valid CLI status keys for help text
func ValidStatusKeys() []string {
	keys := make([]string, 0, len(StatusMap))
	for k := range StatusMap {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

// MapStatusKey converts CLI key to Jira status name
func MapStatusKey(cliKey string) (string, error) {
	if jiraName, ok := StatusMap[cliKey]; ok {
		return jiraName, nil
	}
	return "", fmt.Errorf("invalid status key '%s', valid keys: %v", cliKey, ValidStatusKeys())
}
