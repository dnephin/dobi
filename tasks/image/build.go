package image

import (
	"fmt"
	"io"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/dnephin/dobi/config"
	"github.com/dnephin/dobi/logging"
	"github.com/dnephin/dobi/tasks/context"
	"github.com/dnephin/dobi/utils/fs"
	"github.com/docker/docker/pkg/jsonmessage"
	"github.com/docker/docker/pkg/term"
	docker "github.com/fsouza/go-dockerclient"
)

// BuildTask creates a Docker image
type BuildTask struct {
	name   string
	config *config.ImageConfig
}

// NewBuildTask creates a new BuildTask object
func NewBuildTask(name string, conf *config.ImageConfig) *BuildTask {
	return &BuildTask{name: name, config: conf}
}

// Name returns the name of the task
func (t *BuildTask) Name() string {
	return t.name
}

func (t *BuildTask) String() string {
	return fmt.Sprintf("image.BuildTask(name=%s, config=%s)", t.name, t.config)
}

func (t *BuildTask) logger() *log.Entry {
	return logging.Log.WithFields(log.Fields{"task": t})
}

// Repr formats the task for logging
func (t *BuildTask) Repr() string {
	return fmt.Sprintf("[image:build %s] %s", t.name, t.config.Image)
}

// Run builds or pulls an image if it is out of date
func (t *BuildTask) Run(ctx *context.ExecuteContext) error {
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
	ctx.SetModified(t.name)
	t.logger().Info("Created")
	return nil
}

func (t *BuildTask) isStale(ctx *context.ExecuteContext) (bool, error) {
	if ctx.IsModified(t.config.Dependencies()...) {
		return true, nil
	}

	image, err := GetImage(ctx, t.config)
	switch err {
	case docker.ErrNoSuchImage:
		t.logger().Debug("Image does not exist")
		return true, nil
	case nil:
	default:
		return true, err
	}

	mtime, err := fs.LastModified(t.config.Context)
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

func (t *BuildTask) build(ctx *context.ExecuteContext) error {
	out := os.Stdout
	outFd, isTTY := term.GetFdInfo(out)
	rpipe, wpipe := io.Pipe()
	defer rpipe.Close()

	errChan := make(chan error)

	go func() {
		err := jsonmessage.DisplayJSONMessagesStream(rpipe, out, outFd, isTTY, nil)
		errChan <- err
	}()

	err := ctx.Client.BuildImage(docker.BuildImageOptions{
		Name:           GetImageName(ctx, t.config),
		Dockerfile:     t.config.Dockerfile,
		BuildArgs:      buildArgs(t.config.Args),
		Pull:           t.config.PullBaseImageOnBuild,
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

func (t *BuildTask) tag(ctx *context.ExecuteContext) error {
	// The first one is already tagged in build
	if len(t.config.Tags) <= 1 {
		return nil
	}
	for _, tag := range t.config.Tags[1:] {
		err := ctx.Client.TagImage(GetImageName(ctx, t.config), docker.TagImageOptions{
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
func (t *BuildTask) Prepare(ctx *context.ExecuteContext) error {
	for _, tag := range t.config.Tags {
		if _, err := ctx.Env.Resolve(tag); err != nil {
			return err
		}
	}
	return nil
}

// Stop the task
func (t *BuildTask) Stop(ctx *context.ExecuteContext) error {
	return nil
}
