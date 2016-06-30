package tasks

import (
	"fmt"

	log "github.com/Sirupsen/logrus"
	"github.com/dnephin/dobi/config"
)

// AliasTask is a task which creates a directory on the host
type AliasTask struct {
	baseTask
	config *config.AliasConfig
}

// NewAliasTask creates a new AliasTask object
func NewAliasTask(options taskOptions, conf *config.AliasConfig) *AliasTask {
	return &AliasTask{
		baseTask: baseTask{name: options.name},
		config:   conf,
	}
}

func (t *AliasTask) String() string {
	return fmt.Sprintf("AliasTask(name=%s, config=%s)", t.name, t.config)
}

func (t *AliasTask) logger() *log.Entry {
	return log.WithFields(log.Fields{"task": t})
}

// Repr formats the task for logging
func (t *AliasTask) Repr() string {
	return fmt.Sprintf("[alias %s]", t.name)
}

// Run does nothing. Dependencies were already run.
func (t *AliasTask) Run(ctx *ExecuteContext) error {
	t.logger().Debug("Run")

	if ctx.isModified(t.config.Tasks...) {
		ctx.setModified(t.name)
	}
	t.logger().Info("Done")
	return nil
}

// Prepare the task
func (t *AliasTask) Prepare(ctx *ExecuteContext) error {
	return nil
}

// Stop the task
func (t *AliasTask) Stop(ctx *ExecuteContext) error {
	return nil
}
