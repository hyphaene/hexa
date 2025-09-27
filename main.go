package main

import "github.com/hyphaene/hexa/cmd"

// Variables injected by GoReleaser at build time
var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	// Pass version info to cmd package
	cmd.SetVersionInfo(version, commit, date)
	cmd.Execute()
}
