package alias

import (
	"fmt"

	"github.com/dnephin/dobi/config"
	"github.com/dnephin/dobi/logging"
	"github.com/dnephin/dobi/tasks/context"
	"github.com/dnephin/dobi/tasks/task"
)

// Task is a task which creates a directory on the host
type Task struct {
	name   task.Name
	config *config.AliasConfig
}

// Name returns the name of the task
func (t *Task) Name() task.Name {
	return t.name
}

// Repr formats the task for logging
func (t *Task) Repr() string {
	return fmt.Sprintf("[alias:%s %s]", t.name.Action(), t.name.Resource())
}

// Run does nothing. Dependencies were already run.
func (t *Task) Run(ctx *context.ExecuteContext, depsModified bool) (bool, error) {
	logging.ForTask(t).Info("Done")
	return depsModified, nil
}

// Stop the task
func (t *Task) Stop(ctx *context.ExecuteContext) error {
	return nil
}
