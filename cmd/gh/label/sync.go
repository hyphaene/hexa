package label

import (
	"encoding/json"
	"os"
	"os/exec"

	"github.com/charmbracelet/huh"
	"github.com/hyphaene/hexa/internal/config"
	"github.com/hyphaene/hexa/internal/gh"
	"github.com/hyphaene/hexa/internal/git"
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
	}

	for _, check := range checks {
		if err := check(); err != nil {
			return err
		}
	}

	// fmt.Fprintf(os.Stdout, "json ok\n%+v\n", labels)

	return writeLabelsInProject()
}

func writeLabelsInProject() error {
	// projectConfFilePath, err := config.GetProjectConfigPath()
	// if err != nil {
	// 	return err
	// }

	labels, err := gh.ReadGithubLabelsFromConfig()
	if err != nil {
		return err
	}

	// labelsField, err := config.ReadYAMLField(projectConfFilePath, "github.labels")
	// if err != nil {
	// 	return err
	// }

	if labels != nil {
		huh.NewConfirm().
			Title("Data already here, want to sync again?").
			Affirmative("Yes!").
			Negative("No.").
			Value(&confirm).Run()
	} else {
		huh.NewConfirm().
			Title("No data found, fetch from github ?").
			Affirmative("Yes!").
			Negative("No.").
			Value(&confirm).Run()
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

		jq := exec.Command("jq", ".")
		jq.Stdout = os.Stdout
		jq.Stderr = os.Stderr

		stdin, err := jq.StdinPipe()
		if err != nil {
			return err
		}

		if err := jq.Start(); err != nil {
			return err
		}

		if _, err := stdin.Write(jsonOutput); err != nil {
			return err
		}
		if err := stdin.Close(); err != nil {
			return err
		}

		if err := jq.Wait(); err != nil {
			return err
		}

	}

	return nil
}
