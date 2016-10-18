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
	name   common.TaskName
	config *config.AliasConfig
}

// Name returns the name of the task
func (t *Task) Name() common.TaskName {
	return t.name
}

func (t *Task) logger() *log.Entry {
	return logging.Log.WithFields(log.Fields{"task": t})
}

// Repr formats the task for logging
func (t *Task) Repr() string {
	return fmt.Sprintf("[alias:%s %s]", t.name.Action(), t.name.Resource())
}

// Run does nothing. Dependencies were already run.
func (t *Task) Run(ctx *context.ExecuteContext, depsModified bool) (bool, error) {
	t.logger().Info("Done")
	return depsModified, nil
}

// Stop the task
func (t *Task) Stop(ctx *context.ExecuteContext) error {
	return nil
}
