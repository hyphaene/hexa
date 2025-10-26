package git

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func VerifyCurrentDirectoryIsGitRepo() error {
	cmd := exec.Command("git", "rev-parse", "--git-dir")
	if _, err := cmd.CombinedOutput(); err != nil {
		return errors.New("pas dans un repository git")
	} else {
		os.Stdout.WriteString("Dans un repository git\n")
	}
	return nil
}

func GetRepoRootPath() (string, error) {
	cmd := exec.Command("git", "rev-parse", "--show-toplevel")
	if sdtinOuput, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("git rev-parse failed: %s (%w)", strings.TrimSpace(string(sdtinOuput)), err)
	} else {
		return strings.TrimSpace(string(sdtinOuput)), nil
	}
}
