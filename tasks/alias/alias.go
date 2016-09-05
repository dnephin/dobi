package alias

import (
	"fmt"

	log "github.com/Sirupsen/logrus"
	"github.com/dnephin/dobi/config"
	"github.com/dnephin/dobi/logging"
	"github.com/dnephin/dobi/tasks/common"
	"github.com/dnephin/dobi/tasks/context"
)

// Task is a task which creates a directory on the host
type Task struct {
	name   string
	config *config.AliasConfig
	action action
}

type action struct {
	name         string
	Dependencies func(*Task) []string
}

// NewTask creates a new Task object
func NewTask(name string, conf *config.AliasConfig, act action) *Task {
	return &Task{name: name, config: conf, action: act}
}

// Name returns the name of the task
func (t *Task) Name() common.TaskName {
	return common.NewTaskName(t.name, t.action.name)
}

func (t *Task) logger() *log.Entry {
	return logging.Log.WithFields(log.Fields{"task": t})
}

// Repr formats the task for logging
func (t *Task) Repr() string {
	return fmt.Sprintf("[alias:%s %s]", t.action.name, t.name)
}

// Run does nothing. Dependencies were already run.
func (t *Task) Run(ctx *context.ExecuteContext) error {
	if ctx.IsModified(t.config.Tasks...) {
		ctx.SetModified(t.name)
	}
	t.logger().Info("Done")
	return nil
}

// Dependencies returns the dependencies for the task
func (t *Task) Dependencies() []string {
	return t.action.Dependencies(t)
}

// Stop the task
func (t *Task) Stop(ctx *context.ExecuteContext) error {
	return nil
}
