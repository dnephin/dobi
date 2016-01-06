package tasks

import (
	"fmt"

	log "github.com/Sirupsen/logrus"
	"github.com/dnephin/buildpipe/config"
)

// ImageTask creates a Docker image
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

func (t *ImageTask) String() string {
	return fmt.Sprintf("ImageTask(name=%s, config=%s)", t.name, t.config)
}

func (t *ImageTask) logger() *log.Entry {
	return log.WithFields(log.Fields{
		"task":       "Image",
		"name":       t.name,
		"image":      t.config.Image,
		"dockerfile": t.config.Dockerfile,
		"context":    t.config.Context,
	})
}

// Run builds or pulls an image if it is out of date
func (t *ImageTask) Run() error {
	t.logger().Info("run")

	if !t.isStale() {
		return nil
	}

	return t.build()
}

func (t *ImageTask) isStale() bool {
	return false
}

func (t *ImageTask) build() error {
	return nil
}
