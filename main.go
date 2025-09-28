package main

import (
	"github.com/hyphaene/hexa/cmd"
	"github.com/hyphaene/hexa/internal/config"
	_ "github.com/hyphaene/hexa/internal/env"

	// Import commands to trigger their init() functions
	_ "github.com/hyphaene/hexa/cmd/config"
	_ "github.com/hyphaene/hexa/cmd/jira"
	_ "github.com/hyphaene/hexa/cmd/jira/ticket"
	_ "github.com/hyphaene/hexa/cmd/self"
)

// Variables injected by GoReleaser at build time
var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	// Initialize configuration before anything else
	config.Initialize()
	cmd.SetVersionInfo(version, commit, date)
	cmd.Execute()
}
