package job

import (
	"fmt"
	"github.com/dnephin/dobi/config"
	"github.com/dnephin/dobi/logging"
	"github.com/dnephin/dobi/tasks/context"
	"github.com/dnephin/dobi/tasks/task"
	"github.com/dnephin/dobi/tasks/types"
	docker "github.com/fsouza/go-dockerclient"
	"github.com/pkg/errors"
	"path/filepath"
)

// StopTask is a task which stops the container if running
type StopTask struct {
	types.NoStop
	name   task.Name
	config *config.JobConfig
}

func newStopTask(name task.Name, conf config.Resource) types.Task {
	return &StopTask{name: name, config: conf.(*config.JobConfig)}
}

// Name returns the name of the task
func (t *StopTask) Name() task.Name {
	return t.name
}

// Repr formats the task for logging
func (t *StopTask) Repr() string {
	return fmt.Sprintf("%s %v", t.name.Format("job"), t.config.Artifact)
}

// Run creates the host path if it doesn't already exist
func (t *StopTask) Run(ctx *context.ExecuteContext, _ bool) (bool, error) {
	logger := logging.ForTask(t)
	name := ContainerName(ctx, t.name.Resource())
	id, ok, err := t.isServing(ctx)
	if err != nil {
		return false, errors.Wrap(err, "checking if container is running")
	}
	if ok {
		logger.Infof("Stopping container %s", name)
		ctx.Client.StopContainer(id, uint(0))
		return true, nil
	}
	logger.Infof("Container %s is not running", name)
	return true, nil
}

func (t *StopTask) isServing(ctx *context.ExecuteContext) (string, bool, error) {
	name := ContainerName(ctx, t.name.Resource())
	containers, err := ctx.Client.ListContainers(docker.ListContainersOptions{All: true})
	if err != nil {
		return "", false, err
	}
	for _, container := range containers {
		for _, n := range container.Names {
			if name == filepath.Base(n) {
				return container.ID, true, nil
			}
		}
	}
	return "", false, nil
}
