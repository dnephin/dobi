package mount

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/dnephin/dobi/config"
	"github.com/dnephin/dobi/logging"
	"github.com/dnephin/dobi/tasks/context"
	"github.com/dnephin/dobi/tasks/task"
)

// Task is a mount task
type Task struct {
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
	return fmt.Sprintf("[mount:%s %s] %s:%s",
		t.name.Action(), t.name.Resource(), t.config.Bind, t.config.Path)
}

// Run performs the task action
func (t *Task) Run(ctx *context.ExecuteContext, _ bool) (bool, error) {
	return t.run(t, ctx)
}

// Stop implements the types.Task interface
func (t *Task) Stop(*context.ExecuteContext) error {
	return nil
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

	if err := t.create(ctx); err != nil {
		return false, err
	}
	logger.Info("Created")
	return true, nil
}

func (t *createAction) create(ctx *context.ExecuteContext) error {
	path := AbsBindPath(t.task.config, ctx.WorkingDir)
	mode := os.FileMode(t.task.config.Mode)

	switch t.task.config.File {
	case true:
		return ioutil.WriteFile(path, []byte{}, mode)
	default:
		return os.MkdirAll(path, mode)
	}
}

func (t *createAction) exists(ctx *context.ExecuteContext) bool {
	_, err := os.Stat(AbsBindPath(t.task.config, ctx.WorkingDir))
	if err != nil {
		return false
	}

	return true
}
