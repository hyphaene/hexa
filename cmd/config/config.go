package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"

	"github.com/hyphaene/hexa/cmd"
	"github.com/hyphaene/hexa/cmd/env"
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
		rootConfig := getRootConfig()
		projectConfig := getProjectConfig()
		rootConfigYAML, _ := yaml.Marshal(rootConfig)
		projectConfigYAML, _ := yaml.Marshal(projectConfig)
		fmt.Println("---")
		fmt.Println("Root Config:")
		fmt.Println(string(rootConfigYAML))
		fmt.Println("---")
		fmt.Println("Project Config:")
		fmt.Println(string(projectConfigYAML))
	},
}

func getRootConfig() map[string]interface{} {
	// Construction du path complet et configuration de Viper
	homeDir, _ := os.UserHomeDir()
	configPath := filepath.Join(homeDir, ".hexa.yml")

	if env.Debug {
		fmt.Println("Attempting to read config from:", configPath)
	}

	config := getConfig(configPath)
	return config
}

func getProjectConfig() map[string]interface{} {
	// Construction du path complet et configuration de Viper
	workingDir, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting working directory:", err)
		return nil
	}
	configPath := filepath.Join(workingDir, ".hexa.yml")
	if env.Debug {
		fmt.Println("Attempting to read config from:", configPath)
	}

	config := getConfig(configPath)
	return config
}

func getConfig(configPath string) map[string]interface{} {
	// Construction du path complet et configuration de Viper

	viper.SetConfigFile(configPath)

	// Lecture du fichier de configuration
	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("Error reading config file:", err)
		return nil
	}

	// Confirmation du fichier utilisé
	if env.Debug {
		fmt.Println("Successfully loaded config from:", viper.ConfigFileUsed())
	}
	viperConfig := viper.AllSettings()
	yamlBytes, err := yaml.Marshal(viperConfig)
	if err != nil {
		if env.Debug {
			fmt.Println("Error marshalling config to YAML:", err)
		}
	} else {
		if env.Debug {
			fmt.Printf("Config settings:\n%s\n", yamlBytes)
		}
	}

	return viperConfig
}
