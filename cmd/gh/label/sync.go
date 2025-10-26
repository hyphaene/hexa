package label

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"

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

var confirm bool

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

	labels, err := gh.ReadGithubLabelsFromConfig()
	if err != nil {
		return err
	}

	var promptErr error
	if len(labels) > 0 {
		promptErr = huh.NewConfirm().
			Title("Data already here, want to sync again?").
			Affirmative("Yes!").
			Negative("No.").
			Value(&confirm).Run()
	} else {
		promptErr = huh.NewConfirm().
			Title("No data found, fetch from github ?").
			Affirmative("Yes!").
			Negative("No.").
			Value(&confirm).Run()
	}

	if promptErr != nil {
		return fmt.Errorf("interactive prompt failed (TTY required): %w", promptErr)
	}

	if confirm {
		os.Stdout.WriteString("c'est partiii\n")

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

		jsonOutput, err := json.Marshal(fetchedLabels)
		if err != nil {
			return err
		}

		jqCmd := exec.Command("jq", ".")
		jqCmd.Stdout = os.Stdout
		jqCmd.Stderr = os.Stderr

		stdin, err := jqCmd.StdinPipe()
		if err != nil {
			return err
		}

		if err := jqCmd.Start(); err != nil {
			return err
		}

		if _, err := stdin.Write(jsonOutput); err != nil {
			return err
		}
		if err := stdin.Close(); err != nil {
			return err
		}

		if err := jqCmd.Wait(); err != nil {
			return err
		}

	}

	return nil
}
