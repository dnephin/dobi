package mount

import (
	"fmt"

	"github.com/dnephin/dobi/config"
	"github.com/dnephin/dobi/tasks/context"
	"github.com/dnephin/dobi/tasks/task"
	"github.com/dnephin/dobi/tasks/types"
)

// GetTaskConfig returns a new task for the action
func GetTaskConfig(name, action string, conf *config.MountConfig) (types.TaskConfig, error) {

	newTaskConfig := func(name task.Name, builder types.TaskBuilder) (types.TaskConfig, error) {
		return types.NewTaskConfig(name, conf, task.NoDependencies, builder), nil
	}
	switch action {
	case "", "create":
		return newTaskConfig(task.NewDefaultName(name, action), NewTask(runCreate))
	case "remove", "rm":
		return newTaskConfig(task.NewName(name, action), NewTask(remove))
	default:
		return nil, fmt.Errorf("Invalid mount action %q for task %q", action, name)
	}
}

// NewTask creates a new Task object
func NewTask(
	runFunc func(task *Task, ctx *context.ExecuteContext) (bool, error)) types.TaskBuilder {
	return func(name task.Name, conf config.Resource) types.Task {
		return &Task{name: name, config: conf.(*config.MountConfig), run: runFunc}
	}
}

func remove(task *Task, ctx *context.ExecuteContext) (bool, error) {
	task.logger().Warn("Bind mounts are not removable")
	return false, nil
}
