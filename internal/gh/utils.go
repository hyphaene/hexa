package gh

import (
	"errors"
	"os/exec"
)

func VerifyGhAuthenticated() error {
	cmd := exec.Command("gh", "auth", "status")
	if _, err := cmd.CombinedOutput(); err != nil {
		return errors.New("not authenticated with gh cli")
	}
	return nil
}

func VerifyRemote() error {
	cmd := exec.Command("gh", "repo", "view")
	if _, err := cmd.CombinedOutput(); err != nil {
		return errors.New("no remote detected")
	}
	return nil
}
