package gh

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
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

func VerifyRemote() error {
	cmd := exec.Command("gh", "repo", "view")
	if _, err := cmd.CombinedOutput(); err != nil {
		return errors.New("no remote detected")
	} else {
		os.Stdout.WriteString("Remote existant :)\n")
	}
	return nil
}

func FetchLabels() error {

	/*
		labels_json=$(gh label list --json name,description,color --jq '.')

	*/

	cmd := exec.Command("gh", "label", "list", "--json", "name,description,color", "--jq", "'.'")

	if jsonResponse, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("git rev-parse failed: %s (%w)", strings.TrimSpace(string(jsonResponse)), err)
	}
	return nil
}
