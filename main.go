package main

import (
	"os"
	"os/exec"
	"strings"

	"github.com/bitrise-io/go-utils/v2/log"
)

type envvar struct {
	name         string
	required     bool
	defaultValue *string
}

func newEnvvar(name string, required bool, defaultValue *string) envvar {
	return envvar{name, required, defaultValue}
}

func main() {
	logger := log.NewLogger()
	buildProduct := newEnvvar("build_product", true, nil)
	setEnv(getEnv(logger, "turbo_project", "turbo_environment", "turbo_library_version", "turbo_version", buildProduct), logger)
	os.Exit(0)
}

func getEnv(logger log.Logger, envvars ...any) map[string]string {
	result := map[string]string{}
	for _, v := range envvars {
		var ev envvar
		switch v := v.(type) {
		case envvar:
			ev = v
		case string:
			ev = newEnvvar(v, false, nil)
		default:
			logger.Warnf("Invalid type passed to getEnv. Skipping.")
			continue
		}
		varvalue := strings.TrimSpace(os.Getenv(ev.name))
		if len(varvalue) > 0 {
			result[ev.name] = varvalue
		} else if ev.defaultValue != nil {
			result[ev.name] = *ev.defaultValue
		} else if ev.required {
			logger.Errorf("No value for the required variable $%s was found.", ev.name)
			os.Exit(1)
		}
		logger.Infof("$%s = %s", ev.name, varvalue)
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
