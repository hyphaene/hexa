package self

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/hyphaene/hexa/cmd"
	"github.com/spf13/cobra"
)

func init() {
	cmd.RootCmd.AddCommand(completionCmd)
	completionCmd.AddCommand(completionInstallCmd)

	// Add standard Cobra completion commands
	completionCmd.AddCommand(&cobra.Command{
		Use:   "bash",
		Short: "Generate bash completion",
		Long:  "Generate bash completion script for hexa CLI.",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Root().GenBashCompletion(os.Stdout)
		},
	})

	completionCmd.AddCommand(&cobra.Command{
		Use:   "zsh",
		Short: "Generate zsh completion",
		Long:  "Generate zsh completion script for hexa CLI.",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Root().GenZshCompletion(os.Stdout)
		},
	})

	completionCmd.AddCommand(&cobra.Command{
		Use:   "fish",
		Short: "Generate fish completion",
		Long:  "Generate fish completion script for hexa CLI.",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Root().GenFishCompletion(os.Stdout, true)
		},
	})
}

var completionCmd = &cobra.Command{
	Use:   "completion",
	Short: "Generate and manage shell completions",
	Long:  `Generate and manage shell completions for hexa CLI.`,
}

var completionInstallCmd = &cobra.Command{
	Use:   "install",
	Short: "Install zsh completion for hexa CLI",
	Long: `Install zsh completion for hexa CLI.
This command will:
1. Generate completion file
2. Create completion directory if needed
3. Install completion file
4. Provide instructions for activation`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Installing zsh completion for hexa CLI...")

		// Generate completion file in current directory
		tempCompletionFile := "_hexa"
		fmt.Println("→ Generating completion file...")

		// Generate completion and save to file
		output, err := exec.Command("hexa", "completion", "zsh").Output()
		if err != nil {
			fmt.Printf("Error generating completion: %v\n", err)
			return
		}

		err = os.WriteFile(tempCompletionFile, output, 0644)
		if err != nil {
			fmt.Printf("Error writing completion file: %v\n", err)
			return
		}

		// Create completion directory
		completionDir := "/usr/local/share/zsh/site-functions"
		fmt.Printf("→ Creating completion directory: %s\n", completionDir)
		execCommand("sudo", "mkdir", "-p", completionDir)

		// Install completion file
		completionPath := filepath.Join(completionDir, "_hexa")
		fmt.Printf("→ Installing completion file to: %s\n", completionPath)
		execCommand("sudo", "cp", tempCompletionFile, completionPath)

		// Clean up temp file
		os.Remove(tempCompletionFile)

		fmt.Println("\n✅ Zsh completion installed successfully!")
		fmt.Println("\nTo activate completion, restart your shell or run:")
		fmt.Println("  source ~/.zshrc")
		fmt.Println("\nThen you can use tab completion with:")
		fmt.Println("  hexa <TAB>")
		fmt.Println("  hw <TAB>")
	},
}