package tasks

import (
	"github.com/dnephin/buildpipe/config"
)

// ImageTask
type ImageTask struct {
	baseTask
	config *config.ImageConfig
}

// NewImageTask creates a new ImageTask object
func NewImageTask(options taskOptions, conf *config.ImageConfig) *ImageTask {
	return &ImageTask{
		baseTask: baseTask{
			name:   options.name,
			client: options.client,
		},
		config: conf,
	}
}

// Run builds or pulls an image if it is out of date
func (t *ImageTask) Run() error {
	return nil
}
