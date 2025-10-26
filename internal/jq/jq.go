package jq

import (
	"errors"
	"os/exec"
)

// VerifyJqInstalled checks if jq binary is available in PATH
func VerifyJqInstalled() error {
	cmd := exec.Command("jq", "--version")
	if _, err := cmd.CombinedOutput(); err != nil {
		return errors.New("jq is required for syntax highlighting\nInstall: brew install jq (macOS) | apt install jq (Linux)")
	}
	return nil
}