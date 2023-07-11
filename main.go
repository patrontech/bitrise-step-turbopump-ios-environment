package main

import (
	"os"
	"os/exec"
	"strings"

	"github.com/bitrise-io/go-utils/v2/log"
)

func main() {
	logger := log.NewLogger()
	setEnv(getEnv(logger, "turbo_project", "turbo_environment", "turbo_library_version"), logger)
	os.Exit(0)
}

// getEnv retrieves the given variables from the environment as a map.
// Note that it exits the program early using os.Exit(1) if the variable
// is not set. All variables are assumed to be required.
func getEnv(logger log.Logger, varnames ...string) map[string]string {
	result := map[string]string{}
	for _, varname := range varnames {
		varvalue := os.Getenv(varname)
		if len(strings.TrimSpace(varvalue)) == 0 {
			logger.Errorf("No value for the required variable $%s was found.", varname)
			os.Exit(1)
		}
		result[varname] = varvalue
	}
	return result
}

func setEnv(vars map[string]string, logger log.Logger) {
	for key, value := range vars {
		// You can find more usage examples on envman's GitHub page
		//  at: https://github.com/bitrise-io/envman
		cmdLog, err := exec.Command("bitrise", "envman", "add", "--key", key, "--value", value).CombinedOutput()
		if err != nil {
			logger.Errorf("Failed to expose output with envman, error: %#v | output: %s", err, cmdLog)
			os.Exit(1)
		}
	}
}
