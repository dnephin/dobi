package run

import (
	"fmt"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/dnephin/dobi/config"
	"github.com/dnephin/dobi/logging"
	"github.com/dnephin/dobi/tasks/common"
	"github.com/dnephin/dobi/tasks/context"
)

// RemoveTask is a task which removes the container used by the run task and the
// artifact created by the run task.
type RemoveTask struct {
	name   string
	config *config.JobConfig
}

// NewRemoveTask creates a new RemoveTask object
func NewRemoveTask(name string, conf *config.JobConfig) *RemoveTask {
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
	return fmt.Sprintf("[run:rm %v] %v", t.name, t.config.Artifact)
}

// Run creates the host path if it doesn't already exist
func (t *RemoveTask) Run(ctx *context.ExecuteContext) error {
	t.logger().Debug("Run")
	RemoveContainer(t.logger(), ctx.Client, ContainerName(ctx, t.name))

	if t.config.Artifact != "" {
		if err := os.RemoveAll(t.config.Artifact); err != nil {
			t.logger().Warnf("failed to remove artifact %s: %s", t.config.Artifact, err)
		}
	}

	t.logger().Info("Removed")
	return nil
}

// Dependencies returns the list of dependencies. The remove task doesn't depend
// on anything.
func (t *RemoveTask) Dependencies() []string {
	return []string{}
}

// Stop the task
func (t *RemoveTask) Stop(ctx *context.ExecuteContext) error {
	return nil
}
