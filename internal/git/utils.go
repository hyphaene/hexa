package git

import (
	"errors"
	"fmt"
	"os/exec"
	"strings"
)

// VerifyCurrentDirectoryIsGitRepo checks if the current directory is inside a git repository.
// Returns error if not in a git repo or if git command fails.
func VerifyCurrentDirectoryIsGitRepo() error {
	cmd := exec.Command("git", "rev-parse", "--git-dir")
	if _, err := cmd.CombinedOutput(); err != nil {
		return errors.New("pas dans un repository git")
	}
	return nil
}

// GetRepoRootPath returns the absolute path to the root of the git repository.
// Uses git rev-parse --show-toplevel to find the repo root.
func GetRepoRootPath() (string, error) {
	cmd := exec.Command("git", "rev-parse", "--show-toplevel")
	stdoutOutput, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to get repo root: %s (%w)", strings.TrimSpace(string(stdoutOutput)), err)
	}
	return strings.TrimSpace(string(stdoutOutput)), nil
}
