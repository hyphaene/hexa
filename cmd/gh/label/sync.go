package label

import (
	"os"

	"github.com/hyphaene/hexa/internal/gh"
	"github.com/hyphaene/hexa/internal/git"
	"github.com/spf13/cobra"
)

func init() {
	LabelCommand.AddCommand(LabelSyncCommand)
}

var LabelSyncCommand = &cobra.Command{
	Use:   "sync",
	Short: "Sync GitHub labels from a configuration file",
	Long:  `Sync GitHub labels from a configuration file.`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := git.VerifyCurrentDirectoryIsGitRepo(); err != nil {
			os.Stderr.WriteString(err.Error() + "\n")
			os.Exit(1)
		}
		if err := gh.VerifyGhAuthenticated(); err != nil {
			os.Stderr.WriteString(err.Error() + "\n")
			os.Exit(1)
		}
		if err := gh.VerifyRemote(); err != nil {
			os.Stderr.WriteString(err.Error() + "\n")
			os.Exit(1)
		}

		// if err := verifyGhCLIInstalled(); err != nil {
		// 	os.Stderr.WriteString("GitHub CLI (gh) is not installed or not found in PATH\n")
		// 	os.Exit(1)
		// }
		// Placeholder for the actual label sync logic
		os.Stdout.WriteString("Label sync command executed\n")
	},
}
