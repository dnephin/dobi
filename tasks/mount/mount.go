package mount

import (
	"fmt"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/dnephin/dobi/config"
	"github.com/dnephin/dobi/logging"
	"github.com/dnephin/dobi/tasks/context"
)

// CreateTask is a task which creates a directory on the host
type CreateTask struct {
	name   string
	config *config.MountConfig
}

// NewCreateTask creates a new CreateTask object
func NewCreateTask(name string, conf *config.MountConfig) *CreateTask {
	return &CreateTask{name: name, config: conf}
}

// Name returns the name of the task
func (t *CreateTask) Name() string {
	return t.name
}

func (t *CreateTask) logger() *log.Entry {
	return logging.Log.WithFields(log.Fields{"task": t})
}

// Repr formats the task for logging
func (t *CreateTask) Repr() string {
	return fmt.Sprintf("[mount:create %s] %s:%s", t.name, t.config.Bind, t.config.Path)
}

// Run creates the host path if it doesn't already exist
func (t *CreateTask) Run(ctx *context.ExecuteContext) error {
	t.logger().Debug("Run")

	if t.exists(ctx) {
		t.logger().Debug("is fresh")
		return nil
	}

	err := os.MkdirAll(AbsBindPath(t.config, ctx.WorkingDir), 0777)
	if err != nil {
		return err
	}
	ctx.SetModified(t.name)
	t.logger().Info("Created")
	return nil
}

func (t *CreateTask) exists(ctx *context.ExecuteContext) bool {
	_, err := os.Stat(AbsBindPath(t.config, ctx.WorkingDir))
	if err != nil {
		return false
	}

	return true
}

// Stop the task
func (t *CreateTask) Stop(ctx *context.ExecuteContext) error {
	return nil
}
