package gh

import (
	"os"

	"github.com/spf13/cobra"
)

func init() {
	GHCommand.AddCommand(LabelSyncCommand)
	if err := verifyGhCLIInstalled(); err != nil {
		os.Stderr.WriteString("GitHub CLI (gh) is not installed or not found in PATH\n")
		os.Exit(1)
	} else {
		os.Stdout.WriteString("GitHub CLI (gh) is installed\n")
	}
}

var LabelSyncCommand = &cobra.Command{
	Use:   "label-sync",
	Short: "Sync GitHub labels from a configuration file",
	Long:  `Sync GitHub labels from a configuration file.`,
	Run: func(cmd *cobra.Command, args []string) {

		// if err := verifyGhCLIInstalled(); err != nil {
		// 	os.Stderr.WriteString("GitHub CLI (gh) is not installed or not found in PATH\n")
		// 	os.Exit(1)
		// }
		// Placeholder for the actual label sync logic
		os.Stdout.WriteString("Label sync command executed\n")
	},
}

func verifyGhCLIInstalled() error {
	// Placeholder for actual verification logic

	return nil
}
