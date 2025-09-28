package config

import (
	"github.com/hyphaene/hexa/cmd"
	"github.com/spf13/cobra"
)

func init() {
	cmd.RootCmd.AddCommand(ConfigCmd)
}

var ConfigCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage configuration for hexa CLI",
	Long:  `Manage configuration for hexa CLI.`,
}
