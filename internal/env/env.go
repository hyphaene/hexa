package env

import (
	"os"
	"strings"
)

var (
	Debug bool
)

func init() {
	Debug = readBool("DEBUG")
}

func readBool(envVar string) bool {
	val, exists := os.LookupEnv(envVar)

	if exists && strings.ToLower(val) == "true" {
		return true
	}
	return false
}
