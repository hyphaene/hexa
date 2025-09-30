package jira

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/spf13/viper"
)

// BoardListResponse représente la réponse de l'API /rest/agile/1.0/board
type BoardListResponse struct {
	MaxResults int     `json:"maxResults"`
	StartAt    int     `json:"startAt"`
	Total      int     `json:"total"`
	IsLast     bool    `json:"isLast"`
	Values     []Board `json:"values"`
}

// Board représente un board Jira Agile
type Board struct {
	ID       int    `json:"id"`
	Self     string `json:"self"`
	Name     string `json:"name"`
	Type     string `json:"type"` // "scrum", "kanban", etc.
	Location struct {
		ProjectID   int    `json:"projectId"`
		DisplayName string `json:"displayName"`
		ProjectName string `json:"projectName"`
		ProjectKey  string `json:"projectKey"`
	} `json:"location"`
}

// GetBoardIdFromName récupère l'ID d'un board depuis son nom
func GetBoardIdFromName(boardName string) (int, error) {
	jiraToken := viper.GetString("jira.token")
	baseURL := viper.GetString("jira.url")

	// URL encode le nom du board
	encodedName := url.QueryEscape(boardName)
	apiURL := fmt.Sprintf("%s/rest/agile/1.0/board?name=%s", baseURL, encodedName)

	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return 0, fmt.Errorf("creating request: %w", err)
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", "Bearer "+jiraToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, fmt.Errorf("executing request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	var boardResp BoardListResponse
	if err := json.NewDecoder(resp.Body).Decode(&boardResp); err != nil {
		return 0, fmt.Errorf("decoding response: %w", err)
	}

	// Chercher le board avec un match exact du nom
	for _, board := range boardResp.Values {
		if board.Name == boardName {
			return board.ID, nil
		}
	}

	return 0, fmt.Errorf("no board found with exact name '%s'", boardName)
}