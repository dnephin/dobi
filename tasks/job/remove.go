package job

import (
	"fmt"
	"os"

	"github.com/dnephin/dobi/config"
	"github.com/dnephin/dobi/logging"
	"github.com/dnephin/dobi/tasks/context"
	"github.com/dnephin/dobi/tasks/task"
	"github.com/dnephin/dobi/tasks/types"
)

// RemoveTask is a task which removes the container used by the run task and the
// artifact created by the run task.
type RemoveTask struct {
	name   task.Name
	config *config.JobConfig
}

func newRemoveTask(name task.Name, conf config.Resource) types.Task {
	return &RemoveTask{name: name, config: conf.(*config.JobConfig)}
}

// Name returns the name of the task
func (t *RemoveTask) Name() task.Name {
	return t.name
}

// Repr formats the task for logging
func (t *RemoveTask) Repr() string {
	return fmt.Sprintf("[job:rm %v] %v", t.name.Resource(), t.config.Artifact)
}

// Run creates the host path if it doesn't already exist
func (t *RemoveTask) Run(ctx *context.ExecuteContext, _ bool) (bool, error) {
	logger := logging.ForTask(t)

	RemoveContainer(logger, ctx.Client, ContainerName(ctx, t.name.Resource()), false)

	if t.config.Artifact != "" {
		if err := os.RemoveAll(t.config.Artifact); err != nil {
			logger.Warnf("failed to remove artifact %s: %s", t.config.Artifact, err)
		}
	}

	logger.Info("Removed")
	return true, nil
}

// Stop the task
func (t *RemoveTask) Stop(ctx *context.ExecuteContext) error {
	return nil
}
