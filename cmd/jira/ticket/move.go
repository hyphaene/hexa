package ticket

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	TicketCmd.AddCommand(moveTicketCmd)
}

var moveTicketCmd = &cobra.Command{
	Use:   "move",
	Short: "Move a Jira ticket",
	Long:  `Move a Jira ticket to a different status.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Moving Jira ticket...")
		// Add move logic here
	},
}
