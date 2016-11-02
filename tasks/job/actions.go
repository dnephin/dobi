package job

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/dnephin/dobi/config"
	"github.com/dnephin/dobi/tasks/task"
	"github.com/dnephin/dobi/tasks/types"
)

// GetTaskConfig returns a new task for the action
func GetTaskConfig(name, action string, conf *config.JobConfig) (types.TaskConfig, error) {
	switch action {
	case "", "run":
		return types.NewTaskConfig(
			task.NewDefaultName(name, action),
			conf,
			deps(conf),
			newRunTask), nil
	case "remove", "rm":
		return types.NewTaskConfig(
			task.NewName(name, action),
			conf,
			task.NoDependencies,
			newRemoveTask), nil
	}
	if strings.HasPrefix(action, "capture") {
		variable, err := parseCapture(action)
		if err != nil {
			return nil, err
		}
		return types.NewTaskConfig(
			task.NewName(name, action),
			conf,
			deps(conf),
			newCaptureTask(variable)), nil
	}
	return nil, fmt.Errorf("Invalid run action %q for task %q", action, name)
}

func deps(conf *config.JobConfig) func() []string {
	return func() []string {
		return conf.Dependencies()
	}
}

var (
	captureRegex = regexp.MustCompile(`^capture\((\w+)\)$`)
)

func parseCapture(action string) (string, error) {
	matches := captureRegex.FindStringSubmatch(action)
	if len(matches) > 1 {
		return matches[1], nil
	}
	return "", fmt.Errorf("invalid capture format %q", action)
}
