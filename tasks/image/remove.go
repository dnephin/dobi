package image

import (
	"fmt"

	log "github.com/Sirupsen/logrus"
	"github.com/dnephin/dobi/config"
	"github.com/dnephin/dobi/logging"
	"github.com/dnephin/dobi/tasks/context"
)

// RemoveTask creates a Docker image
type RemoveTask struct {
	name   string
	config *config.ImageConfig
}

// NewRemoveTask creates a new RemoveTask object
func NewRemoveTask(name string, conf *config.ImageConfig) *RemoveTask {
	return &RemoveTask{name: name, config: conf}
}

// Name returns the name of the task
func (t *RemoveTask) Name() string {
	return t.name
}

func (t *RemoveTask) String() string {
	return fmt.Sprintf("image.RemoveTask(name=%s, config=%s)", t.name, t.config)
}

func (t *RemoveTask) logger() *log.Entry {
	return logging.Log.WithFields(log.Fields{"task": t})
}

// Repr formats the task for logging
func (t *RemoveTask) Repr() string {
	return fmt.Sprintf("[image:remove %s] %s", t.name, t.config.Image)
}

// Run builds or pulls an image if it is out of date
func (t *RemoveTask) Run(ctx *context.ExecuteContext) error {
	t.logger().Debug("Run")

	t.logger().Info("Removing")
	if err := t.removeTags(ctx); err != nil {
		return err
	}
	ctx.SetModified(t.name)
	t.logger().Info("Removed")
	return nil
}

func (t *RemoveTask) removeTags(ctx *context.ExecuteContext) error {
	for _, tag := range t.config.Tags {
		tag = ctx.Env.GetVar(tag)
		if err := ctx.Client.RemoveImage(t.config.Image + ":" + tag); err != nil {
			t.logger().Warnf("failed to remove %q: %s", tag, err)
		}
	}
	return nil
}

// Prepare the task
// TODO: move this to common function or maybe prepare should be based on the
// resource?
func (t *RemoveTask) Prepare(ctx *context.ExecuteContext) error {
	for _, tag := range t.config.Tags {
		if _, err := ctx.Env.Resolve(tag); err != nil {
			return err
		}
	}
	return nil
}

// Stop the task
func (t *RemoveTask) Stop(ctx *context.ExecuteContext) error {
	return nil
}
