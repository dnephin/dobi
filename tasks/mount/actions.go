package mount

import (
	"fmt"

	"github.com/dnephin/dobi/config"
	"github.com/dnephin/dobi/tasks/context"
	"github.com/dnephin/dobi/tasks/task"
	"github.com/dnephin/dobi/tasks/types"
)

// GetTaskConfig returns a new task for the action
func GetTaskConfig(name task.Name, conf *config.MountConfig) (types.TaskConfig, error) {

	newTaskConfig := func(name task.Name, builder types.TaskBuilder) (types.TaskConfig, error) {
		return types.NewTaskConfig(name, conf, task.NoDependencies(), builder), nil
	}
	switch name.Action() {
	case task.Create:
		return newTaskConfig(name, NewTask(runCreate))
	case task.Remove:
		return newTaskConfig(name, NewTask(remove))
	default:
		return nil, fmt.Errorf("invalid mount action %q for task %q", name.Action(), name.Resource())
	}
}

// NewTask creates a new Task object
func NewTask(
	runFunc func(task *Task, ctx *context.ExecuteContext) (bool, error)) types.TaskBuilder {
	return func(name task.Name, conf config.Resource) types.Task {
		return &Task{name: name, config: conf.(*config.MountConfig), run: runFunc}
	}
}
