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
func (t *ImageTask) Run(ctx *ExecuteContext) error {
	t.logger().Info("run")

	stale, err := t.isStale(ctx)
	if !stale || err != nil {
		return err
	}

	err = t.build()
	if err != nil {
		return err
	}
	ctx.setModified(t.name)
	t.logger().Info("created")
	return nil
}

func (t *ImageTask) isStale(ctx *ExecuteContext) (bool, error) {
	if ctx.isModified(t.config.Dependencies()...) {
		return true, nil
	}

	image, err := t.getImage(ctx)
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

func (t *ImageTask) getImage(ctx *ExecuteContext) (*docker.Image, error) {
	// TODO: this should use the unique run id for the tag
	return t.client.InspectImage(t.config.Image + ":todo-unique")
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
