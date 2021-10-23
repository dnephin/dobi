package image

import (
	"fmt"
	"io"

	"github.com/dnephin/dobi/config"
	"github.com/dnephin/dobi/logging"
	"github.com/dnephin/dobi/tasks/context"
	"github.com/dnephin/dobi/tasks/task"
	"github.com/dnephin/dobi/tasks/types"
	docker "github.com/fsouza/go-dockerclient"
	"github.com/moby/moby/pkg/jsonmessage"
	"github.com/moby/term"
	log "github.com/sirupsen/logrus"
)

// Task creates a Docker image
type Task struct {
	types.NoStop
	name    task.Name
	config  *config.ImageConfig
	runFunc runFunc
}

// Name returns the name of the task
func (t *Task) Name() task.Name {
	return t.name
}

func (t *Task) logger() *log.Entry {
	return logging.ForTask(t)
}

// Repr formats the task for logging
func (t *Task) Repr() string {
	return fmt.Sprintf("%s %s", t.name.Format("image"), t.config.Image)
}

// Run builds or pulls an image if it is out of date
func (t *Task) Run(ctx *context.ExecuteContext, depsModified bool) (bool, error) {
	return t.runFunc(ctx, t, depsModified)
}

// ForEachTag runs a function for each tag
func (t *Task) ForEachTag(ctx *context.ExecuteContext, each func(string) error) error {
	if err := t.forEachLocalTag(ctx, each); err != nil {
		return err
	}

	return t.forEachRemoteTagNoFallback(each)
}

// forEachLocalTag runs a function for each local tag
func (t *Task) forEachLocalTag(ctx *context.ExecuteContext, each func(string) error) error {
	if len(t.config.Tags) == 0 {
		return t.forEachProvidedTag(each, []string{GetImageName(ctx, t.config)})
	}

	return t.forEachProvidedTag(each, t.config.Tags)
}

// ForEachRemoteTag runs a function for each remote tag, if no remote tags uses local tags
func (t *Task) ForEachRemoteTag(ctx *context.ExecuteContext, each func(string) error) error {
	if len(t.config.RemoteTags) == 0 {
		return t.forEachLocalTag(ctx, each)
	}

	return t.forEachRemoteTagNoFallback(each)
}

// forEachRemoteTagNoFallback runs a function for each remote tag
func (t *Task) forEachRemoteTagNoFallback(each func(string) error) error {
	return t.forEachProvidedTag(each, t.config.RemoteTags)
}

// forEachProvidedTag runs a function for each provided tag
func (t *Task) forEachProvidedTag(each func(string) error, tags []string) error {
	for _, tag := range tags {
		imageTag := tag
		// Create complete image name if the tag does not already have it.
		if _, hasTag := docker.ParseRepositoryTag(tag); hasTag == "" {
			imageTag = t.config.Image + ":" + tag
		}

		if err := each(imageTag); err != nil {
			return err
		}
	}
	return nil
}

// Stream json output to a terminal
func Stream(out io.Writer, streamer func(out io.Writer) error) error {
	outFd, isTTY := term.GetFdInfo(out)
	rpipe, wpipe := io.Pipe()
	defer rpipe.Close() // nolint: errcheck

	errChan := make(chan error)

	go func() {
		err := jsonmessage.DisplayJSONMessagesStream(rpipe, out, outFd, isTTY, nil)
		errChan <- err
	}()

	err := streamer(wpipe)
	wpipe.Close() // nolint: errcheck
	if err != nil {
		<-errChan
		return err
	}
	return <-errChan
}
