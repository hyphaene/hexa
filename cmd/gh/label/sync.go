package label

import (
	"fmt"

	"github.com/charmbracelet/huh"
	"github.com/hyphaene/hexa/internal/config"
	"github.com/hyphaene/hexa/internal/gh"
	"github.com/hyphaene/hexa/internal/git"
	"github.com/hyphaene/hexa/internal/jq"
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
		jq.VerifyJqInstalled,
		config.EnsureProjectConfigExists,
	}

	for _, check := range checks {
		if err := check(); err != nil {
			return err
		}
	}

	return writeLabelsInProject()
}

func writeLabelsInProject() error {
	confirmed, err := promptUserForSync()
	if err != nil || !confirmed {
		return err
	}
	return syncLabelsToConfig()
}

func promptUserForSync() (bool, error) {
	confirm := true
	err := huh.NewConfirm().
		Title("Fetch labels from GitHub and save to .hexa.yml?").
		Description("This will overwrite any existing labels configuration").
		Affirmative("Yes!").
		Negative("No.").
		Value(&confirm).Run()

	if err != nil {
		return false, fmt.Errorf("interactive prompt failed (TTY required): %w", err)
	}
	return confirm, nil
}

func syncLabelsToConfig() error {
	fetchedLabels, err := gh.FetchLabels()
	if err != nil {
		return err
	}

	projectConfFilePath, err := config.GetProjectConfigPath()
	if err != nil {
		return err
	}

	if err := config.UpdateYAMLField(projectConfFilePath, "github.labels", fetchedLabels); err != nil {
		return err
	}

	return jq.PrettyPrint(fetchedLabels)
}
