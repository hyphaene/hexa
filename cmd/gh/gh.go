package gh

import (
	"github.com/hyphaene/hexa/cmd"
	"github.com/spf13/cobra"
)

func init() {
	cmd.RootCmd.AddCommand(GHCommand)
}

var GHCommand = &cobra.Command{
	Use:   "gh",
	Short: "GitHub CLI",
	Long:  `gh is the GitHub command line tool.`,
}
