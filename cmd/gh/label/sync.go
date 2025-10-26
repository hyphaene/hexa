package label

import (
	"fmt"
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
	RunE: func(cmd *cobra.Command, args []string) error {
		return execute()
	},
}

func execute() error {
	checks := []func() error{
		git.VerifyCurrentDirectoryIsGitRepo,
		gh.VerifyGhAuthenticated,
		gh.VerifyRemote,
	}

	for _, check := range checks {
		if err := check(); err != nil {
			return err
		}
	}

	path, err := git.GetRepoRootPath()
	if err != nil {
		return err
	}

	fmt.Fprintf(os.Stdout, "here is the repo/worktree path\n\n%s\n\n", path)
	return nil
}
