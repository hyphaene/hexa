package jira

import (
	"fmt"

	internalJira "github.com/hyphaene/hexa/internal/jira"
	"github.com/spf13/cobra"
)

var getCurrentSprintIdCmd = &cobra.Command{
	Use:   "get-current-sprint-id",
	Short: "Get the ID of the current active sprint",
	Long: `Retrieves the ID of the currently active sprint for the configured Jira board.

Prerequisites:
  1. Run 'hexa jira init --board-name "YOUR_BOARD" --config-path .hexa.local.yml' first
  2. This caches the board ID and avoids repeated API calls

Example:
  hexa jira get-current-sprint-id`,
	RunE: runGetCurrentSprintId,
}

func init() {
	JiraCmd.AddCommand(getCurrentSprintIdCmd)
}

func runGetCurrentSprintId(cmd *cobra.Command, args []string) error {
	sprintID, err := internalJira.GetCurrentSprintId()
	if err != nil {
		return fmt.Errorf("failed to get current sprint ID: %w", err)
	}

	fmt.Printf("🎯 Current active sprint ID: %d\n", sprintID)
	return nil
}