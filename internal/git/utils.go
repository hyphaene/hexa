package git

import (
	"errors"
	"fmt"
	"os/exec"
	"strings"
)

func VerifyCurrentDirectoryIsGitRepo() error {
	cmd := exec.Command("git", "rev-parse", "--git-dir")
	if _, err := cmd.CombinedOutput(); err != nil {
		return errors.New("pas dans un repository git")
	}
	return nil
}

func GetRepoRootPath() (string, error) {
	cmd := exec.Command("git", "rev-parse", "--show-toplevel")
	stdoutOutput, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to get repo root: %s (%w)", strings.TrimSpace(string(stdoutOutput)), err)
	}
	return strings.TrimSpace(string(stdoutOutput)), nil
}
