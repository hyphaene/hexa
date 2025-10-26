package jq

import (
	"encoding/json"
	"errors"
	"os"
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

// PrettyPrint marshals data to JSON and pipes it through jq for colored output
func PrettyPrint(data any) error {
	jsonOutput, err := json.Marshal(data)
	if err != nil {
		return err
	}

	jqCmd := exec.Command("jq", ".")
	jqCmd.Stdout = os.Stdout
	jqCmd.Stderr = os.Stderr

	stdin, err := jqCmd.StdinPipe()
	if err != nil {
		return err
	}

	if err := jqCmd.Start(); err != nil {
		return err
	}

	if _, err := stdin.Write(jsonOutput); err != nil {
		return err
	}

	if err := stdin.Close(); err != nil {
		return err
	}

	return jqCmd.Wait()
}