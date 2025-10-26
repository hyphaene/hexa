package gh

import (
	"fmt"
	"os/exec"
	"strings"
)

// VerifyGhAuthenticated checks if gh CLI is authenticated with GitHub.
// Returns error with gh CLI output if authentication fails.
func VerifyGhAuthenticated() error {
	cmd := exec.Command("gh", "auth", "status")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("gh authentication failed: %s", strings.TrimSpace(string(output)))
	}
	return nil
}

// VerifyRemote checks if the current directory has a GitHub remote configured.
// Returns error with gh CLI output if no remote is detected.
func VerifyRemote() error {
	cmd := exec.Command("gh", "repo", "view")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("no GitHub remote detected: %s", strings.TrimSpace(string(output)))
	}
	return nil
}
