package types

import (
	"github.com/dnephin/dobi/config"
	"github.com/dnephin/dobi/logging"
	"github.com/dnephin/dobi/tasks/common"
	"github.com/dnephin/dobi/tasks/context"
)

// Task interface performs some operation with a resource config
type Task interface {
	logging.LogRepresenter
	Name() common.TaskName
	Run(*context.ExecuteContext, bool) (bool, error)
	Stop(*context.ExecuteContext) error
}

// RunFunc is a function which performs the task. It received a context and a
// bool indicating if any dependencies were modified. It should return true if
// the resource was modified, otherwise false.
type RunFunc func(*context.ExecuteContext, bool) (bool, error)

// TaskConfig is a data object which stores the full configuration of a Task
type TaskConfig interface {
	Name() common.TaskName
	Resource() config.Resource
	Dependencies() []string
	Task(config.Resource) Task
}

type taskConfig struct {
	name      common.TaskName
	resource  config.Resource
	deps      func() []string
	buildTask func(common.TaskName, config.Resource) Task
}

func (t taskConfig) Name() common.TaskName {
	return t.name
}

func (t taskConfig) Resource() config.Resource {
	return t.resource
}

func (t taskConfig) Dependencies() []string {
	return t.deps()
}

func (t taskConfig) Task(res config.Resource) Task {
	return t.buildTask(t.name, res)
}

// TaskBuilder is a function which creates a new Task from a name and config
type TaskBuilder func(common.TaskName, config.Resource) Task

// NewTaskConfig returns a TaskConfig from components
func NewTaskConfig(
	name common.TaskName,
	resource config.Resource,
	deps func() []string,
	buildTask TaskBuilder,
) TaskConfig {
	return &taskConfig{
		name:      name,
		resource:  resource,
		deps:      deps,
		buildTask: buildTask,
	}
}
