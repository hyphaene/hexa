package ticket

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/hyphaene/hexa/internal/jira"
)

func init() {
	TicketCmd.AddCommand(getTicketCmd)
}

var getTicketCmd = &cobra.Command{
	Use:   "get",
	Short: "Get details of a Jira ticket",
	Long:  `Fetch and display details of a specific Jira ticket.`,
	Run: func(cmd *cobra.Command, args []string) {
		sprintId, err := jira.GetCurrentSprintId()
		if err != nil {
			fmt.Println("Error fetching current sprint ID:", err)
			return
		}
		fmt.Println("Current Sprint ID:", sprintId)
		// Add logic to fetch and display ticket details here
	},
}
