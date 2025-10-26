package gh

import (
	"fmt"
	"os/exec"
	"strings"
)

func VerifyGhAuthenticated() error {
	cmd := exec.Command("gh", "auth", "status")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("gh authentication failed: %s", strings.TrimSpace(string(output)))
	}
	return nil
}

func VerifyRemote() error {
	cmd := exec.Command("gh", "repo", "view")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("no GitHub remote detected: %s", strings.TrimSpace(string(output)))
	}
	return nil
}
