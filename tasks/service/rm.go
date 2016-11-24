package service

import (
	"bytes"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/dnephin/dobi/config"
	"github.com/dnephin/dobi/logging"
	"github.com/dnephin/dobi/tasks/context"
	"github.com/dnephin/dobi/tasks/task"
	"github.com/dnephin/dobi/tasks/types"
	"github.com/docker/docker/api/types/swarm"
	docker "github.com/fsouza/go-dockerclient"
	"io"
)

func newRemoveTask(name task.Name, conf config.Resource) types.Task {
	return &RemoveTask{name: name, config: conf.(*config.ServiceConfig)}
}

// RemoveTask is a task which remove a docker service
type RemoveTask struct {
	types.NoStop
	name      task.Name
	config    *config.ServiceConfig
	outStream io.Writer
}

// Name returns the name of the task
func (t *RemoveTask) Name() task.Name {
	return t.name
}

func (t *RemoveTask) logger() *log.Entry {
	return logging.ForTask(t)
}

// Repr formats the task for logging
func (t *RemoveTask) Repr() string {
	buff := &bytes.Buffer{}
	return fmt.Sprintf("%s%v", t.name.Format("job"), buff.String())
}

// Run the job command in a container
func (t *RemoveTask) removeOptions(ctx *context.ExecuteContext, service swarm.Service) docker.RemoveServiceOptions {
	return docker.RemoveServiceOptions{ID: service.ID}
}

// Run containers as a docker service
func (t *RemoveTask) Run(ctx *context.ExecuteContext, depsModified bool) (bool, error) {
	service, running, err := isServiceRunning(ctx, t.name.Resource())
	if err != nil {
		return false, err

	}
	if running {
		err := ctx.Client.RemoveService(t.removeOptions(ctx, service))
		if err != nil {
			return false, err
		}
		t.logger().Info("removed")
		return true, nil
	}

	t.logger().Info("not running, wont remove")
	return true, nil
}
