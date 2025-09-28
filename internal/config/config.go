package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"

	"github.com/hyphaene/hexa/internal/env"
	"github.com/spf13/viper"
)

// Initialize loads root and project configurations into the global Viper instance
func Initialize() {
	rootConfig := getRootConfig()
	projectConfig := getProjectConfig()

	// Clear any existing config
	viper.Reset()

	// Merge configurations with project config taking precedence
	viper.MergeConfigMap(rootConfig)
	viper.MergeConfigMap(projectConfig) // project override root

	if env.Debug {
		fmt.Println("Viper configuration initialized successfully")
	}
}

// GetMergedConfig returns the complete merged configuration for debugging
func GetMergedConfig() map[string]interface{} {
	return viper.AllSettings()
}

func getRootConfig() map[string]interface{} {
	homeDir, _ := os.UserHomeDir()
	configPath := filepath.Join(homeDir, ".hexa.yml")

	if env.Debug {
		fmt.Println("Attempting to read root config from:", configPath)
	}

	return getConfig(configPath)
}

func getProjectConfig() map[string]interface{} {
	workingDir, err := os.Getwd()
	if err != nil {
		if env.Debug {
			fmt.Println("Error getting working directory:", err)
		}
		return nil
	}

	configPath := filepath.Join(workingDir, ".hexa.yml")
	if env.Debug {
		fmt.Println("Attempting to read project config from:", configPath)
	}

	return getConfig(configPath)
}

func getConfig(configPath string) map[string]interface{} {
	// Create a new Viper instance for this config file
	v := viper.New()
	v.SetConfigFile(configPath)

	if err := v.ReadInConfig(); err != nil {
		if env.Debug {
			fmt.Printf("Config file not found or error reading %s: %v\n", configPath, err)
		}
		return nil
	}

	if env.Debug {
		fmt.Println("Successfully loaded config from:", v.ConfigFileUsed())
		yamlBytes, err := yaml.Marshal(v.AllSettings())
		if err == nil {
			fmt.Printf("Config settings:\n%s\n", yamlBytes)
		}
	}

	return v.AllSettings()
}