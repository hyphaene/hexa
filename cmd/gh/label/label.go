package label

import (
	"github.com/hyphaene/hexa/cmd/gh"
	"github.com/spf13/cobra"
)

var LabelCommand = &cobra.Command{
	Use:   "label",
	Short: "Manage GitHub labels",
	Long:  `Commands to manage GitHub labels.`,
}

func init() {
	gh.GHCommand.AddCommand(LabelCommand)
}
