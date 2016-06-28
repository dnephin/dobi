package tasks

import (
	"fmt"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/dnephin/dobi/config"
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
		baseTask: baseTask{name: options.name},
		config:   conf,
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
	t.logger().Debug("image is stale")

	if err := t.build(ctx); err != nil {
		return err
	}
	if err = t.tag(ctx); err != nil {
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
		t.logger().Debug("image does not exist")
		return true, nil
	case nil:
	default:
		return true, err
	}

	// TODO: support .dockerignore
	mtime, err := lastModified(t.config.Context)
	if err != nil {
		t.logger().Warnf("Failed to get last modified time of context.")
		return true, err
	}
	if image.Created.Before(mtime) {
		t.logger().Debug("image older than context")
		return true, nil
	}
	return false, nil
}

func (t *ImageTask) getImage(ctx *ExecuteContext) (*docker.Image, error) {
	return ctx.client.InspectImage(t.getImageName(ctx))
}

func (t *ImageTask) getImageName(ctx *ExecuteContext) string {
	return fmt.Sprintf("%s:%s", t.config.Image, t.getCanonicalTag(ctx))
}

func (t *ImageTask) getCanonicalTag(ctx *ExecuteContext) string {
	if len(t.config.Tags) > 0 {
		// TODO: environment.Resolve()
		return t.config.Tags[1]
	}
	return ctx.environment.Unique()
}

func (t *ImageTask) build(ctx *ExecuteContext) error {
	return ctx.client.BuildImage(docker.BuildImageOptions{
		Name:           t.getImageName(ctx),
		Dockerfile:     t.config.Dockerfile,
		Pull:           t.config.Pull,
		RmTmpContainer: true,
		ContextDir:     t.config.Context,
		// TODO: support quiet, or send to loggeR?
		OutputStream: os.Stdout,
	})
}

func (t *ImageTask) tag(ctx *ExecuteContext) error {
	// The first one is already tagged in build
	if len(t.config.Tags) <= 1 {
		return nil
	}
	for _, tag := range t.config.Tags[1:] {
		// TODO: environment.Resolve()
		err := ctx.client.TagImage(t.getImageName(ctx), docker.TagImageOptions{
			Repo:  t.config.Image,
			Tag:   tag,
			Force: true,
		})
		if err != nil {
			return err
		}
	}
	return nil
}
