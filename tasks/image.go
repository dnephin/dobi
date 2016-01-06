package tasks

import (
	"fmt"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/dnephin/buildpipe/config"
	"github.com/fsouza/go-dockerclient"
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

	stale, err := t.isStale()
	if err != nil {
		return err
	}
	if !stale {
		return nil
	}

	return t.build()
}

func (t *ImageTask) isStale() (bool, error) {
	// TODO: this should use the unique run id for the tag
	image, err := t.client.InspectImage(t.config.Image + ":todo-unique")
	switch err {
	case docker.ErrNoSuchImage:
		return true, nil
	case nil:
	default:
		return true, err
	}

	// Images without a context can never be stale
	if t.config.Context == "" {
		return false, nil
	}

	// TODO: support .dockerignore
	mtime, err := lastModified(t.config.Context)
	if err != nil {
		t.logger().Warnf("Failed to get last modified time of context.")
		return true, err
	}
	return image.Created.Before(mtime), nil
}

func (t *ImageTask) build() error {
	return t.client.BuildImage(docker.BuildImageOptions{
		Name:           t.config.Image,
		Dockerfile:     t.config.Dockerfile,
		Pull:           t.config.Pull,
		RmTmpContainer: true,
		ContextDir:     t.config.Context,
		// TODO: support quiet, or send to loggeR?
		OutputStream: os.Stdout,
	})
}
