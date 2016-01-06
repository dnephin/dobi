package tasks

import (
	"fmt"
	//"os"

	log "github.com/Sirupsen/logrus"
	"github.com/dnephin/buildpipe/config"
	//"github.com/fsouza/go-dockerclient"
)

// CommandTask is a task which runs a command in a container to produce a
// file or set of files.
type CommandTask struct {
	baseTask
	config *config.CommandConfig
}

// NewCommandTask creates a new CommandTask object
func NewCommandTask(options taskOptions, conf *config.CommandConfig) *CommandTask {
	return &CommandTask{
		baseTask: baseTask{
			name:   options.name,
			client: options.client,
		},
		config: conf,
	}
}

func (t *CommandTask) String() string {
	return fmt.Sprintf("CommandTask(name=%s, config=%s)", t.name, t.config)
}

func (t *CommandTask) logger() *log.Entry {
	return log.WithFields(log.Fields{
		"task":     "Command",
		"name":     t.name,
		"use":      t.config.Use,
		"command":  t.config.Command,
		"artifact": t.config.Artifact,
	})
}

// Run creates the host path if it doesn't already exist
func (t *CommandTask) Run(ctx *ExecuteContext) error {
	return nil
}

func (t *CommandTask) isStale(ctx *ExecuteContext) bool {
	if ctx.isModified(t.config.Dependencies()...) {
		return true
	}

	if t.config.Artifact == "" {
		return true
	}

	return false
}
