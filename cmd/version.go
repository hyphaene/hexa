package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of hexa",
	Long:  `All software has versions. This is hexa's`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Hexa CLI v0.1.0 -- HEAD")
	},
}
