package service

import (
	"fmt"
	"github.com/dnephin/dobi/config"
	"github.com/dnephin/dobi/tasks/task"
	"github.com/dnephin/dobi/tasks/types"
)

// GetTaskConfig returns a new task for the action
func GetTaskConfig(name, action string, conf *config.ServiceConfig) (types.TaskConfig, error) {
	switch action {
	case "", "run":
		return types.NewTaskConfig(
			task.NewDefaultName(name, action),
			conf,
			deps(conf),
			newServeTask), nil
	case "rm", "remove":
		return types.NewTaskConfig(
			task.NewDefaultName(name, action),
			conf,
			deps(conf),
			newRemoveTask), nil
	}

	return nil, fmt.Errorf("invalid run action %q for task %q", action, name)
}

func deps(conf *config.ServiceConfig) func() []string {
	return func() []string {
		return conf.Dependencies()
	}
}
