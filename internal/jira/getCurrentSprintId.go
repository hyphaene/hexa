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

func GetCurrentSprintId() (int, error) {
	jiraToken := viper.GetString("jira.token")

	const SEE_SOP_BOARD_ID = "14242"

	url := viper.GetString("jira.url") + "/rest/agile/1.0/board/" + SEE_SOP_BOARD_ID + "/sprint?state=active"

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

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("failed to fetch sprints, status code: %d", response.StatusCode)

	}

	var sprintResp SprintListResponse
	if err := json.NewDecoder(response.Body).Decode(&sprintResp); err != nil {
		return 0, fmt.Errorf("decoding response: %w", err)
	}

	return sprintResp.Values[0].ID, nil
}
