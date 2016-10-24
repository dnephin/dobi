package mount

import (
	"fmt"

	"github.com/dnephin/dobi/config"
	"github.com/dnephin/dobi/tasks/common"
	"github.com/dnephin/dobi/tasks/context"
	"github.com/dnephin/dobi/tasks/iface"
)

// GetTask returns a new task for the action
func GetTask(name, action string, conf *config.MountConfig) (iface.TaskConfig, error) {

	newTaskConfig := func(name common.TaskName, builder iface.TaskBuilder) (iface.TaskConfig, error) {
		return iface.NewTaskConfig(name, conf, noDependencies, builder), nil
	}
	switch action {
	case "", "create":
		return newTaskConfig(common.NewDefaultTaskName(name, action), NewTask(runCreate))
	case "remove", "rm":
		return newTaskConfig(common.NewTaskName(name, action), NewTask(remove))
	default:
		return nil, fmt.Errorf("Invalid mount action %q for task %q", action, name)
	}
}

// NewTask creates a new Task object
func NewTask(
	runFunc func(task *Task, ctx *context.ExecuteContext) (bool, error)) iface.TaskBuilder {
	return func(name common.TaskName, conf config.Resource) iface.Task {
		return &Task{name: name, config: conf.(*config.MountConfig), run: runFunc}
	}
}

func remove(task *Task, ctx *context.ExecuteContext) (bool, error) {
	task.logger().Warn("Bind mounts are not removable")
	return false, nil
}

func noDependencies() []string {
	return []string{}
}
