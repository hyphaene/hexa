package jira

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/hyphaene/hexa/internal/config"
	internalJira "github.com/hyphaene/hexa/internal/jira"
	"github.com/spf13/cobra"
)

var (
	boardName  string
	configPath string
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize Jira configuration with board ID resolution",
	Long: `Resolves a Jira board ID from its name and stores it in the specified config file.
This avoids repeated API calls to resolve the board ID on every command execution.

Example:
  hexa jira init --board-name "SEE x SOP" --config-path .hexa.local.yml`,
	RunE: runInit,
}

func init() {
	initCmd.Flags().StringVar(&boardName, "board-name", "", "Name of the Jira board to resolve (required)")
	initCmd.Flags().StringVar(&configPath, "config-path", "", "Path to the config file to update (required)")
	initCmd.MarkFlagRequired("board-name")
	initCmd.MarkFlagRequired("config-path")

	JiraCmd.AddCommand(initCmd)
}

func runInit(cmd *cobra.Command, args []string) error {
	fmt.Printf("🔍 Resolving board ID for '%s'...\n", boardName)

	// Résoudre le Board ID via API
	boardID, err := internalJira.GetBoardIdFromName(boardName)
	if err != nil {
		return fmt.Errorf("failed to resolve board ID: %w", err)
	}

	fmt.Printf("✅ Board found: '%s' (ID: %d)\n", boardName, boardID)

	// Résoudre le chemin absolu du fichier config (expand ~ si nécessaire)
	expandedPath := configPath
	if len(configPath) > 0 && configPath[0] == '~' {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("getting home directory: %w", err)
		}
		expandedPath = filepath.Join(homeDir, configPath[1:])
	}

	absPath, err := filepath.Abs(expandedPath)
	if err != nil {
		return fmt.Errorf("resolving config path: %w", err)
	}

	// Vérifier si le fichier existe, sinon créer
	if _, err := os.Stat(absPath); os.IsNotExist(err) {
		fmt.Printf("📝 Creating config file: %s\n", absPath)
		if err := os.WriteFile(absPath, []byte{}, 0644); err != nil {
			return fmt.Errorf("creating config file: %w", err)
		}
	}

	// Écrire jira.boardId avec notation pointée (préserve les autres champs de jira)
	if err := config.UpdateYAMLField(absPath, "jira.boardId", boardID); err != nil {
		return fmt.Errorf("updating config file: %w", err)
	}

	fmt.Printf("✅ Configuration saved to: %s\n", absPath)
	fmt.Printf("   jira:\n")
	fmt.Printf("     boardId: %d\n", boardID)
	fmt.Println()
	fmt.Println("💡 Tip: Run 'hexa jira refresh' if the board ID becomes stale.")

	return nil
}
