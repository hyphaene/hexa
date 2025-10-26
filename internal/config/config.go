package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"gopkg.in/yaml.v3"

	"github.com/hyphaene/hexa/internal/env"
	"github.com/hyphaene/hexa/internal/git"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

// Initialize loads root and project configurations into the global Viper instance
func Initialize() {
	err := godotenv.Load()
	if err != nil && env.Debug {
		fmt.Println("No .env file found (this is ok)")
	}
	rootConfig := getRootConfig()
	projectConfig := GetProjectConfig()
	secretProjectConfig := getSecretProjectConfig()

	// Clear any existing config
	viper.Reset()

	// Set defaults before loading configs
	setDefaults()

	viper.SetEnvPrefix("HEXA")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	if rootConfig != nil {
		if err := viper.MergeConfigMap(rootConfig); err != nil && env.Debug {
			fmt.Println("Error merging root config:", err)
		}
	}
	if projectConfig != nil {
		if err := viper.MergeConfigMap(projectConfig); err != nil && env.Debug {
			fmt.Println("Error merging project config:", err)
		}
	}
	if secretProjectConfig != nil {
		if err := viper.MergeConfigMap(secretProjectConfig); err != nil && env.Debug {
			fmt.Println("Error merging secret project config:", err)
		}
	}

	// Validate configuration values after all configs are merged
	if err := ValidateConfig(); err != nil {
		fmt.Fprintf(os.Stderr, "Configuration validation failed: %v\n", err)
		os.Exit(1)
	}

	if env.Debug {
		fmt.Println("Viper configuration initialized successfully")
	}
}

// GetMergedConfig returns the complete merged configuration for debugging
func GetMergedConfig() map[string]any {
	return viper.AllSettings()
}

func getRootConfig() map[string]any {
	homeDir, _ := os.UserHomeDir()
	configPath := filepath.Join(homeDir, ".hexa.yml")

	if env.Debug {
		fmt.Println("Attempting to read root config from:", configPath)
	}

	return getConfig(configPath)
}

// GetProjectConfigPath returns the absolute path to .hexa.yml in the git repository root.
// Returns error if not in a git repository or if repo root cannot be determined.
func GetProjectConfigPath() (string, error) {
	repoRoot, err := git.GetRepoRootPath()
	if err != nil {
		return "", fmt.Errorf("failed to find git repo root: %w", err)
	}
	return filepath.Join(repoRoot, ".hexa.yml"), nil
}

// EnsureProjectConfigExists creates .hexa.yml at repo root with minimal template if it doesn't exist.
// Creates the file automatically without prompting - safe to call on fresh repos.
// Returns error only on real failures (permissions, IO errors).
func EnsureProjectConfigExists() error {
	path, err := GetProjectConfigPath()
	if err != nil {
		return fmt.Errorf("getting project config path: %w", err)
	}

	if _, err := os.Stat(path); err == nil {
		// File already exists
		return nil
	} else if !os.IsNotExist(err) {
		// Real error (permission, etc)
		return fmt.Errorf("checking config file: %w", err)
	}

	// File doesn't exist, create minimal template
	template := `# Hexa Project Configuration

github:
  labels: []
`

	if err := os.WriteFile(path, []byte(template), 0644); err != nil {
		return fmt.Errorf("creating config file: %w", err)
	}

	if env.Debug {
		fmt.Printf("Created project config at: %s\n", path)
	}

	return nil
}

// GetProjectConfig reads and returns the project configuration from .hexa.yml at repo root.
// Returns nil if file doesn't exist or cannot be read (not considered an error).
func GetProjectConfig() map[string]any {

	configPath, _ := GetProjectConfigPath()
	if env.Debug {
		fmt.Println("Attempting to read project config from:", configPath)
	}

	return getConfig(configPath)
}

func getSecretProjectConfig() map[string]any {
	workingDir, err := git.GetRepoRootPath()
	if err != nil {
		if env.Debug {
			fmt.Println("Error getting working directory:", err)
		}
		return nil
	}

	configPath := filepath.Join(workingDir, ".hexa.local.yml")
	if env.Debug {
		fmt.Println("Attempting to read project config from:", configPath)
	}

	return getConfig(configPath)
}

func getConfig(configPath string) map[string]any {
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

// setDefaults sets default values for all configuration keys
func setDefaults() {
	// Jira sprint configuration
	viper.SetDefault("jira.sprint.maxResults", 25)
	viper.SetDefault("jira.sprint.maxRetries", 3)
	viper.SetDefault("jira.sprint.retryDelay", 1*time.Second)
	viper.SetDefault("jira.sprint.timeout", 30*time.Second)
}

// ValidateConfig validates configuration values
func ValidateConfig() error {
	if viper.GetInt("jira.sprint.maxResults") <= 0 {
		return fmt.Errorf("jira.sprint.maxResults must be > 0, got %d", viper.GetInt("jira.sprint.maxResults"))
	}
	if viper.GetInt("jira.sprint.maxRetries") < 1 {
		return fmt.Errorf("jira.sprint.maxRetries must be >= 1, got %d", viper.GetInt("jira.sprint.maxRetries"))
	}
	if viper.GetDuration("jira.sprint.retryDelay") < 0 {
		return fmt.Errorf("jira.sprint.retryDelay must be >= 0, got %s", viper.GetDuration("jira.sprint.retryDelay"))
	}
	if viper.GetDuration("jira.sprint.timeout") <= 0 {
		return fmt.Errorf("jira.sprint.timeout must be > 0, got %s", viper.GetDuration("jira.sprint.timeout"))
	}
	return nil
}
