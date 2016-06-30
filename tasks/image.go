package tasks

import (
	"fmt"
	"io"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/dnephin/dobi/config"
	"github.com/dnephin/dobi/logging"
	"github.com/docker/docker/pkg/jsonmessage"
	"github.com/docker/docker/pkg/term"
	docker "github.com/fsouza/go-dockerclient"
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
	return logging.Log.WithFields(log.Fields{"task": t})
}

// Repr formats the task for logging
func (t *ImageTask) Repr() string {
	return fmt.Sprintf("[image %s] %s", t.name, t.config.Image)
}

// Run builds or pulls an image if it is out of date
func (t *ImageTask) Run(ctx *ExecuteContext) error {
	t.logger().Debug("Run")

	stale, err := t.isStale(ctx)
	if !stale || err != nil {
		t.logger().Info("is fresh")
		return err
	}
	t.logger().Debug("is stale")

	t.logger().Info("Building")
	if err := t.build(ctx); err != nil {
		return err
	}
	if err = t.tag(ctx); err != nil {
		return err
	}
	ctx.setModified(t.name)
	t.logger().Info("Created")
	return nil
}

func (t *ImageTask) isStale(ctx *ExecuteContext) (bool, error) {
	if ctx.isModified(t.config.Dependencies()...) {
		return true, nil
	}

	image, err := t.GetImage(ctx)
	switch err {
	case docker.ErrNoSuchImage:
		t.logger().Debug("Image does not exist")
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
		t.logger().Debug("Image older than context")
		return true, nil
	}
	return false, nil
}

// GetImage returns the image created by this task
func (t *ImageTask) GetImage(ctx *ExecuteContext) (*docker.Image, error) {
	return ctx.client.InspectImage(t.getImageName(ctx))
}

func (t *ImageTask) getImageName(ctx *ExecuteContext) string {
	return fmt.Sprintf("%s:%s", t.config.Image, t.getCanonicalTag(ctx))
}

func (t *ImageTask) getCanonicalTag(ctx *ExecuteContext) string {
	if len(t.config.Tags) > 0 {
		return ctx.Env.GetVar(t.config.Tags[0])
	}
	return ctx.Env.Unique()
}

func (t *ImageTask) build(ctx *ExecuteContext) error {
	out := os.Stdout
	outFd, isTTY := term.GetFdInfo(out)
	rpipe, wpipe := io.Pipe()
	defer rpipe.Close()

	errChan := make(chan error)

	go func() {
		err := jsonmessage.DisplayJSONMessagesStream(rpipe, out, outFd, isTTY, nil)
		errChan <- err
	}()

	err := ctx.client.BuildImage(docker.BuildImageOptions{
		Name:           t.getImageName(ctx),
		Dockerfile:     t.config.Dockerfile,
		BuildArgs:      buildArgs(t.config.Args),
		Pull:           t.config.Pull,
		RmTmpContainer: true,
		ContextDir:     t.config.Context,
		OutputStream:   wpipe,
		RawJSONStream:  true,
		SuppressOutput: ctx.Quiet,
	})
	wpipe.Close()
	if err != nil {
		<-errChan
		return err
	}

	return <-errChan
}

func buildArgs(args map[string]string) []docker.BuildArg {
	out := []docker.BuildArg{}
	for key, value := range args {
		out = append(out, docker.BuildArg{Name: key, Value: value})
	}
	return out
}

func (t *ImageTask) tag(ctx *ExecuteContext) error {
	// The first one is already tagged in build
	if len(t.config.Tags) <= 1 {
		return nil
	}
	for _, tag := range t.config.Tags[1:] {
		err := ctx.client.TagImage(t.getImageName(ctx), docker.TagImageOptions{
			Repo:  t.config.Image,
			Tag:   ctx.Env.GetVar(tag),
			Force: true,
		})
		if err != nil {
			return err
		}
	}
	return nil
}

// Prepare the task
func (t *ImageTask) Prepare(ctx *ExecuteContext) error {
	for _, tag := range t.config.Tags {
		if _, err := ctx.Env.Resolve(tag); err != nil {
			return err
		}
	}
	return nil
}

// Stop the task
func (t *ImageTask) Stop(ctx *ExecuteContext) error {
	return nil
}
