package ticket

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	TicketCmd.AddCommand(commentTicketCmd)
}

var commentTicketCmd = &cobra.Command{
	Use:   "comment",
	Short: "Add a comment to a Jira ticket",
	Long:  `Add a comment to a Jira ticket.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Adding comment to Jira ticket...")
		// Add comment logic here
	},
}
