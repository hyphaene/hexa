package self

import (
	"fmt"
	"strings"

	"github.com/hyphaene/hexa/cmd"
	"github.com/spf13/cobra"
)

func init() {
	cmd.RootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of hexa",
	Long:  `Display the version number of hexa.`,
	Run: func(command *cobra.Command, args []string) {
		// Format version string
		version := cmd.AppVersion
		if version == "" || version == "dev" {
			version = "dev"
		} else if !strings.HasPrefix(version, "v") {
			version = "v" + version
		}

		fmt.Println(version)
	},
}
