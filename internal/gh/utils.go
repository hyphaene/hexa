package gh

import (
	"errors"
	"os"
	"os/exec"
)

func VerifyGhAuthenticated() error {
	cmd := exec.Command("gh", "auth", "status")
	if _, err := cmd.CombinedOutput(); err != nil {
		return errors.New("not authenticated with gh cli")
	} else {
		os.Stdout.WriteString("Authentifié :)\n")
	}
	return nil
}
