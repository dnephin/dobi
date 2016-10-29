package job

import (
	"fmt"

	"github.com/dnephin/dobi/config"
	"github.com/dnephin/dobi/tasks/common"
	"github.com/dnephin/dobi/tasks/types"
)

// GetTaskConfig returns a new task for the action
func GetTaskConfig(name, action string, conf *config.JobConfig) (types.TaskConfig, error) {
	switch action {
	case "", "run":
		return types.NewTaskConfig(
			common.NewDefaultTaskName(name, action),
			conf,
			deps(conf),
			newRunTask), nil
	case "remove", "rm":
		return types.NewTaskConfig(
			common.NewTaskName(name, action),
			conf,
			common.NoDependencies,
			newRemoveTask), nil
	default:
		return nil, fmt.Errorf("Invalid run action %q for task %q", action, name)
	}
}

func deps(conf *config.JobConfig) func() []string {
	return func() []string {
		return conf.Dependencies()
	}
}
