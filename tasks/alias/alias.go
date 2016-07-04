package alias

import (
	"fmt"

	log "github.com/Sirupsen/logrus"
	"github.com/dnephin/dobi/config"
	"github.com/dnephin/dobi/logging"
	"github.com/dnephin/dobi/tasks/context"
)

// Task is a task which creates a directory on the host
type Task struct {
	name   string
	config *config.AliasConfig
}

// NewTask creates a new Task object
func NewTask(name string, conf *config.AliasConfig) *Task {
	return &Task{name: name, config: conf}
}

// Name returns the name of the task
func (t *Task) Name() string {
	return t.name
}

func (t *Task) String() string {
	return fmt.Sprintf("Task(name=%s, config=%s)", t.name, t.config)
}

func (t *Task) logger() *log.Entry {
	return logging.Log.WithFields(log.Fields{"task": t})
}

// Repr formats the task for logging
func (t *Task) Repr() string {
	return fmt.Sprintf("[alias %s]", t.name)
}

// Run does nothing. Dependencies were already run.
func (t *Task) Run(ctx *context.ExecuteContext) error {
	t.logger().Debug("Run")

	if ctx.IsModified(t.config.Tasks...) {
		ctx.SetModified(t.name)
	}
	t.logger().Info("Done")
	return nil
}

// Prepare the task
func (t *Task) Prepare(ctx *context.ExecuteContext) error {
	return nil
}

// Stop the task
func (t *Task) Stop(ctx *context.ExecuteContext) error {
	return nil
}
