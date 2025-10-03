package jira

import (
	"github.com/hyphaene/hexa/cmd"
	"github.com/hyphaene/hexa/cmd/jira/sprint"

	"github.com/spf13/cobra"
)

var JiraCmd = &cobra.Command{
	Use:   "jira",
	Short: "Jira related commands",
	Long:  `Commands to interact with Jira tickets and projects.`,
}

func init() {
	cmd.RootCmd.AddCommand(JiraCmd)
	JiraCmd.AddCommand(sprint.SprintCmd)
}
