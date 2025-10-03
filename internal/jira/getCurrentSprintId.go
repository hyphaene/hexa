package jira

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/spf13/viper"
)

// Réponse complète de l'API
type SprintListResponse struct {
	MaxResults int      `json:"maxResults"`
	StartAt    int      `json:"startAt"`
	IsLast     bool     `json:"isLast"`
	Values     []Sprint `json:"values"`
}

// Un sprint individuel
type Sprint struct {
	ID            int       `json:"id"`
	Self          string    `json:"self"`
	State         string    `json:"state"` // "active", "future", "closed"
	Name          string    `json:"name"`
	StartDate     time.Time `json:"startDate"`
	EndDate       time.Time `json:"endDate"`
	ActivatedDate time.Time `json:"activatedDate"`
	OriginBoardID int       `json:"originBoardId"`
	Goal          string    `json:"goal"`
	Synced        bool      `json:"synced"`
	AutoStartStop bool      `json:"autoStartStop"`
}

// GetSprintIdFromNumber resolves sprint ID from a sprint number (eg. 35 -> "Sprint SEE x SOP 35")
func GetSprintIdFromNumber(sprintNumber int) (int, error) {
	boardName := viper.GetString("jira.boardName")
	if boardName == "" {
		return 0, fmt.Errorf("jira.boardName not configured")
	}

	// Get board ID
	var boardID int
	if viper.IsSet("jira.boardId") {
		boardID = viper.GetInt("jira.boardId")
	} else {
		var err error
		boardID, err = GetBoardIdFromName(boardName)
		if err != nil {
			return 0, fmt.Errorf("resolving board ID: %w", err)
		}
	}

	// Search for sprint matching "Sprint {boardName} {number}"
	sprintName := fmt.Sprintf("Sprint %s %d", boardName, sprintNumber)
	url := fmt.Sprintf("%s/rest/agile/1.0/board/%d/sprint?maxResults=500", viper.GetString("jira.url"), boardID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return 0, fmt.Errorf("creating request: %w", err)
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", "Bearer "+viper.GetString("jira.token"))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, fmt.Errorf("calling API: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	var sprintResp SprintListResponse
	if err := json.NewDecoder(resp.Body).Decode(&sprintResp); err != nil {
		return 0, fmt.Errorf("decoding response: %w", err)
	}

	// Find matching sprint
	for _, sprint := range sprintResp.Values {
		if sprint.Name == sprintName {
			return sprint.ID, nil
		}
	}

	return 0, fmt.Errorf("sprint '%s' not found", sprintName)
}

func GetCurrentSprintId() (int, error) {
	jiraToken := viper.GetString("jira.token")

	// Priorité 1: utiliser jira.boardId si présent (évite appel API)
	var boardID int
	if viper.IsSet("jira.boardId") {
		boardID = viper.GetInt("jira.boardId")
	} else {
		// Priorité 2: résoudre via jira.boardName (fallback)
		boardName := viper.GetString("jira.boardName")
		if boardName == "" {
			return 0, fmt.Errorf("neither jira.boardId nor jira.boardName is configured. Run 'hexa jira init --board-name \"YOUR_BOARD\" --config-path .hexa.local.yml' to initialize")
		}
		var err error
		boardID, err = GetBoardIdFromName(boardName)
		if err != nil {
			return 0, fmt.Errorf("resolving board ID from name '%s': %w. Consider running 'hexa jira init' to cache the board ID", boardName, err)
		}
	}

	url := fmt.Sprintf("%s/rest/agile/1.0/board/%d/sprint?state=active", viper.GetString("jira.url"), boardID)

	resp, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatalln(err)
	}

	resp.Header.Add("Accept", "application/json")
	resp.Header.Add("Authorization", "Bearer "+jiraToken)

	response, err := http.DefaultClient.Do(resp)

	if err != nil {
		log.Fatalln(err)
	}
	defer func() {
		if cerr := response.Body.Close(); cerr != nil {
			log.Printf("closing response body: %v", cerr)
		}
	}()

	if response.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("failed to fetch sprints, status code: %d", response.StatusCode)

	}

	var sprintResp SprintListResponse
	if err := json.NewDecoder(response.Body).Decode(&sprintResp); err != nil {
		return 0, fmt.Errorf("decoding response: %w", err)
	}

	return sprintResp.Values[0].ID, nil
}
