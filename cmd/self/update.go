package self

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
	"strings"

	"github.com/hyphaene/hexa/cmd"
	"github.com/spf13/cobra"
)

const runFormat = "Running: %s %s\n"

func execCommand(command string, args ...string) {

	argsStr := strings.Join(args, " ")
	fmt.Printf(runFormat, command, argsStr)
	cmd := exec.Command(command, args...)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		log.Printf("stderr: %s\n", stderr.String())
		log.Fatalf("error while doing %s: %v", command, err)
	}
	fmt.Println(out.String())
	fmt.Printf("Command %s %s completed successfully\n", command, argsStr)
}

// joinArgs joins command arguments into a single string separated by spaces.

func init() {
	cmd.RootCmd.AddCommand(updateCmd)
}

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update the hexa CLI",
	Long:  `Update the hexa CLI to the latest version.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Updating hexa CLI...")
		setHomebrewTap()
		updateBrew()
		upgradeBin()
		printVersion()
		// Add update logic here
	},
}

func setHomebrewTap() {
	execCommand("brew", "tap", "hyphaene/hexa")
}

func updateBrew() {
	execCommand("brew", "update")
}

func upgradeBin() {
	execCommand("brew", "upgrade", "hexa")

}

func printVersion() {
	execCommand("hexa", "version")
}
