package tasks

import (
	"github.com/dnephin/buildpipe/config"
)

// CommandTask
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

// Run creates the host path if it doesn't already exist
func (t *CommandTask) Run() error {
	return nil
}
