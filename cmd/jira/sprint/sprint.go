package sprint

import (
	"github.com/spf13/cobra"
)

var SprintCmd = &cobra.Command{
	Use:   "sprint",
	Short: "Sprint-related commands",
	Long:  `Commands to interact with Jira sprint data (fetch tickets, pulse overview, etc.)`,
}
