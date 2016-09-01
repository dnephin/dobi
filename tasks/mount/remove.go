package mount

import (
	"fmt"

	log "github.com/Sirupsen/logrus"
	"github.com/dnephin/dobi/config"
	"github.com/dnephin/dobi/logging"
	"github.com/dnephin/dobi/tasks/common"
	"github.com/dnephin/dobi/tasks/context"
)

// RemoveTask is a noop task. It exists because all tasks must have an rm action
type RemoveTask struct {
	name   string
	config *config.MountConfig
}

// NewRemoveTask removes a new RemoveTask object
func NewRemoveTask(name string, conf *config.MountConfig) *RemoveTask {
	return &RemoveTask{name: name, config: conf}
}

// Name returns the name of the task
func (t *RemoveTask) Name() common.TaskName {
	return common.NewTaskName(t.name, "rm")
}

func (t *RemoveTask) logger() *log.Entry {
	return logging.Log.WithFields(log.Fields{"task": t})
}

// Repr formats the task for logging
func (t *RemoveTask) Repr() string {
	return fmt.Sprintf("[mount:rm %s] %s:%s", t.name, t.config.Bind, t.config.Path)
}

// Run does nothing
func (t *RemoveTask) Run(ctx *context.ExecuteContext) error {
	t.logger().Debug("Run")
	t.logger().Warn("Bind mounts are not removable")
	return nil
}

// Dependencies returns the list of dependencies
func (t *RemoveTask) Dependencies() []string {
	return t.config.Dependencies()
}

// Stop the task
func (t *RemoveTask) Stop(ctx *context.ExecuteContext) error {
	return nil
}
