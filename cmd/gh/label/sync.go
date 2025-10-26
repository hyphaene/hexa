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
		execute()
		os.Stdout.WriteString("Label sync command executed\n")
	},
}

func execute() {

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

	if path, err := git.GetRepoRootPath(); err != nil {
		os.Stderr.WriteString(err.Error() + "\n")
		os.Exit(1)
	} else {
		os.Stdout.WriteString("here is the repo/worktree path\n\n")
		os.Stdout.WriteString(path)
		os.Stdout.WriteString("\n\n")
	}
}
