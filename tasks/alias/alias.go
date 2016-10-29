package alias

import (
	"github.com/dnephin/dobi/config"
	"github.com/dnephin/dobi/logging"
	"github.com/dnephin/dobi/tasks/context"
	"github.com/dnephin/dobi/tasks/task"
	"github.com/dnephin/dobi/tasks/types"
)

// Task is an alias task
type Task struct {
	types.NoStop
	name   task.Name
	config *config.AliasConfig
}

// Name returns the name of the task
func (t *Task) Name() task.Name {
	return t.name
}

// Repr formats the task for logging
func (t *Task) Repr() string {
	return t.name.Format("alias")
}

// Run does nothing. Dependencies were already run.
func (t *Task) Run(ctx *context.ExecuteContext, depsModified bool) (bool, error) {
	logging.ForTask(t).Info("Done")
	return depsModified, nil
}
