package cmd

import (
	"github.com/spf13/cobra"
)

var (
	// Used for flags.
	cfgFile     string
	userLicense string

	// Version info injected at build time
	appVersion string
	buildCommit string
	buildDate string

	rootCmd = &cobra.Command{
		Use:   "hexa",
		Short: "Hexactitude CLI - Unified automation and scripting toolkit",
		Long: `Hexa is a unified CLI for automation and scripting tasks.
It replaces 22+ bash scripts with a single, distributable Go binary
organized around functional domains (JIRA, GIT, SETUP, AI).`,
	}
)

// SetVersionInfo sets the version information injected by the build system
func SetVersionInfo(version, commit, date string) {
	appVersion = version
	buildCommit = commit
	buildDate = date
}

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}
