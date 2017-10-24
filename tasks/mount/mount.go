package mount

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/dnephin/dobi/config"
	"github.com/dnephin/dobi/logging"
	"github.com/dnephin/dobi/tasks/context"
	"github.com/dnephin/dobi/tasks/task"
	"github.com/dnephin/dobi/tasks/types"
	docker "github.com/fsouza/go-dockerclient"
)

// Task is a mount task
type Task struct {
	types.NoStop
	name   task.Name
	config *config.MountConfig
	run    func(*Task, *context.ExecuteContext) (bool, error)
}

// Name returns the name of the task
func (t *Task) Name() task.Name {
	return t.name
}

// Repr formats the task for logging
func (t *Task) Repr() string {
	return fmt.Sprintf("%s %s:%s", t.name.Format("mount"), t.config.Bind, t.config.Path)
}

// Run performs the task action
func (t *Task) Run(ctx *context.ExecuteContext, _ bool) (bool, error) {
	return t.run(t, ctx)
}

type createAction struct {
	task *Task
}

// Run creates the host path if it doesn't already exist
func runCreate(task *Task, ctx *context.ExecuteContext) (bool, error) {
	c := createAction{task: task}
	return c.run(ctx)
}

func (t *createAction) run(ctx *context.ExecuteContext) (bool, error) {
	logger := logging.ForTask(t.task)

	if t.exists(ctx) {
		logger.Debug("is fresh")
		return false, nil
	}

	var err error
	switch {
	case t.task.config.IsBind():
		err = t.createBind(ctx)
	default:
		err = t.createNamed(ctx)
	}
	if err != nil {
		return false, err
	}
	logger.Info("Created")
	return true, nil
}

func (t *createAction) createBind(ctx *context.ExecuteContext) error {
	path := AbsBindPath(t.task.config, ctx.WorkingDir)
	mode := os.FileMode(t.task.config.Mode)

	switch t.task.config.File {
	case true:
		return ioutil.WriteFile(path, []byte{}, mode)
	default:
		return os.MkdirAll(path, mode)
	}
}

func (t *createAction) createNamed(ctx *context.ExecuteContext) error {
	_, err := ctx.Client.CreateVolume(docker.CreateVolumeOptions{
		Name: t.task.config.Name,
	})
	return err
}

func (t *createAction) exists(ctx *context.ExecuteContext) bool {
	_, err := os.Stat(AbsBindPath(t.task.config, ctx.WorkingDir))
	return err == nil
}

func remove(task *Task, ctx *context.ExecuteContext) (bool, error) {
	if task.config.Name == "" {
		logging.ForTask(task).Warn("Bind mounts are not removable")
		return false, nil
	}

	return true, ctx.Client.RemoveVolume(task.config.Name)
}
