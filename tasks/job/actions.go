package job

import (
	"fmt"

	"github.com/dnephin/dobi/config"
	"github.com/dnephin/dobi/tasks/common"
	"github.com/dnephin/dobi/tasks/iface"
)

// GetTask returns a new task for the action
func GetTask(name, action string, conf *config.JobConfig) (iface.TaskConfig, error) {
	switch action {
	case "", "run":
		return iface.NewTaskConfig(
			common.NewDefaultTaskName(name, action),
			conf,
			deps(conf),
			newRunTask), nil
	case "remove", "rm":
		return iface.NewTaskConfig(
			common.NewTaskName(name, action),
			conf,
			noDependencies,
			newRemoveTask), nil
	default:
		return nil, fmt.Errorf("Invalid run action %q for task %q", name, action)
	}
}

func deps(conf *config.JobConfig) func() []string {
	return func() []string {
		return conf.Dependencies()
	}
}

func noDependencies() []string {
	return []string{}
}
