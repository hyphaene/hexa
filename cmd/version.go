package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of hexa",
	Long:  `Display the version number of hexa.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Format version string
		version := appVersion
		if version == "" || version == "dev" {
			version = "dev"
		} else if !strings.HasPrefix(version, "v") {
			version = "v" + version
		}

		fmt.Println(version)
	},
}
