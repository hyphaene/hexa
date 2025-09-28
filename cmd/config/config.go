package config

import (
	"fmt"

	"gopkg.in/yaml.v3"

	"github.com/hyphaene/hexa/cmd"
	"github.com/hyphaene/hexa/internal/config"
	"github.com/hyphaene/hexa/internal/env"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	cmd.RootCmd.AddCommand(ConfigCmd)
}

var ConfigCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage configuration for hexa CLI",
	Long:  `Manage configuration for hexa CLI.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Configuration is already loaded in the global Viper instance via main()
		mergedConfig := config.GetMergedConfig()
		mergedConfigYAML, _ := yaml.Marshal(mergedConfig)

		fmt.Println("Debug Mode:", env.Debug)
		fmt.Println("User Name:", viper.GetString("user.me"))
		fmt.Println("---")
		fmt.Println("Active Configuration:")
		fmt.Println(string(mergedConfigYAML))

	},
}
