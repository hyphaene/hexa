package git

import (
	"errors"
	"os"
	"os/exec"
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
