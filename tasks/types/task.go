package types

import (
	"github.com/dnephin/dobi/config"
	"github.com/dnephin/dobi/logging"
	"github.com/dnephin/dobi/tasks/context"
)

// Task interface performs some operation with a resource config
type Task interface {
	logging.LogRepresenter
	Name() Name
	Run(*context.ExecuteContext, bool) (bool, error)
	Stop(*context.ExecuteContext) error
}

// RunFunc is a function which performs the task. It receives a context and a
// bool indicating if any dependencies were modified. It should return true if
// the resource was modified, otherwise false.
type RunFunc func(*context.ExecuteContext, bool) (bool, error)

// TaskConfig is a data object which stores the full configuration of a Task
type TaskConfig interface {
	Name() Name
	Resource() config.Resource
	Dependencies() []string
	Task(config.Resource) Task
}

type taskConfig struct {
	name      Name
	resource  config.Resource
	deps      func() []string
	buildTask func(Name, config.Resource) Task
}

func (t taskConfig) Name() Name {
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
type TaskBuilder func(Name, config.Resource) Task

// NewTaskConfig returns a TaskConfig from components
func NewTaskConfig(
	name Name,
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

// NoStop implements the Stop() method from the types.Task interface. It can be
// used by tasks that don't do anything during the `stop` phase of execution.
type NoStop struct{}

// Stop does nothing
func (t *NoStop) Stop(_ *context.ExecuteContext) error {
	return nil
}
