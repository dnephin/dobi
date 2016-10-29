package job

import (
	"fmt"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/dnephin/dobi/config"
	"github.com/dnephin/dobi/logging"
	"github.com/dnephin/dobi/tasks/common"
	"github.com/dnephin/dobi/tasks/context"
	"github.com/dnephin/dobi/tasks/iface"
)

// RemoveTask is a task which removes the container used by the run task and the
// artifact created by the run task.
type RemoveTask struct {
	name   common.TaskName
	config *config.JobConfig
}

func newRemoveTask(name common.TaskName, conf config.Resource) iface.Task {
	return &RemoveTask{name: name, config: conf.(*config.JobConfig)}
}

// Name returns the name of the task
func (t *RemoveTask) Name() common.TaskName {
	return t.name
}

func (t *RemoveTask) logger() *log.Entry {
	return logging.Log.WithFields(log.Fields{"task": t})
}

// Repr formats the task for logging
func (t *RemoveTask) Repr() string {
	return fmt.Sprintf("[job:rm %v] %v", t.name.Resource(), t.config.Artifact)
}

// Run creates the host path if it doesn't already exist
func (t *RemoveTask) Run(ctx *context.ExecuteContext, _ bool) (bool, error) {
	RemoveContainer(t.logger(), ctx.Client, ContainerName(ctx, t.name.Resource()), false)

	if t.config.Artifact != "" {
		if err := os.RemoveAll(t.config.Artifact); err != nil {
			t.logger().Warnf("failed to remove artifact %s: %s", t.config.Artifact, err)
		}
	}

	t.logger().Info("Removed")
	return true, nil
}

// Stop the task
func (t *RemoveTask) Stop(ctx *context.ExecuteContext) error {
	return nil
}
