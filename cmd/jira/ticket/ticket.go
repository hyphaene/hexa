package ticket

import (
	"github.com/hyphaene/hexa/cmd/jira"

	"github.com/spf13/cobra"
)

func init() {
	jira.JiraCmd.AddCommand(TicketCmd)
}

var (
	TicketCmd = &cobra.Command{
		Use:   "ticket",
		Short: "Jira ticket related commands",
		Long:  `Commands to interact with Jira tickets.`,
	}
)
