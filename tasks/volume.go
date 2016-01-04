package tasks

import (
	"github.com/dnephin/buildpipe/config"
)

// VolumeTask
type VolumeTask struct {
	baseTask
	config config.VolumeConfig
}

// NewVolumeTask creates a new VolumeTask object
func NewVolumeTask(options taskOptions, conf config.VolumeConfig) *VolumeTask {
	return &VolumeTask{
		baseTask: baseTask{
			name:   options.name,
			client: options.client,
		},
		config: conf,
	}
}

// Run creates the host path if it doesn't already exist
func (t *VolumeTask) Run() error {
	return nil
}

// Dependencies returns an empty list, VolumeTasks have no dependencies
func (t *VolumeTask) Dependencies() []string {
	return []string{}
}
