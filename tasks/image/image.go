package image

import (
	"fmt"
	"io"

	log "github.com/Sirupsen/logrus"
	"github.com/dnephin/dobi/config"
	"github.com/dnephin/dobi/logging"
	"github.com/dnephin/dobi/tasks/context"
	"github.com/docker/docker/pkg/jsonmessage"
	"github.com/docker/docker/pkg/term"
)

// Task creates a Docker image
type Task struct {
	name   string
	config *config.ImageConfig
	action action
}

type action struct {
	name string
	Run  func(ctx *context.ExecuteContext, task *Task) error
}

// NewTask creates a new Task object
func NewTask(name string, conf *config.ImageConfig, act action) *Task {
	return &Task{name: name, config: conf, action: act}
}

// Name returns the name of the task
func (t *Task) Name() string {
	return t.name
}

func (t *Task) logger() *log.Entry {
	return logging.Log.WithFields(log.Fields{"task": t})
}

// Repr formats the task for logging
func (t *Task) Repr() string {
	return fmt.Sprintf("[image:%s %s] %s", t.action.name, t.name, t.config.Image)
}

// Run builds or pulls an image if it is out of date
func (t *Task) Run(ctx *context.ExecuteContext) error {
	t.logger().Debug("Run")
	return t.action.Run(ctx, t)
}

// Stop the task
func (t *Task) Stop(ctx *context.ExecuteContext) error {
	return nil
}

// Dependencies returns the list of dependencies
func (t *Task) Dependencies() []string {
	return t.config.Dependencies()
}

// ForEachTag runs a function for each tag
func (t *Task) ForEachTag(ctx *context.ExecuteContext, each func(string) error) error {
	if len(t.config.Tags) == 0 {
		return each(GetImageName(ctx, t.config))
	}

	for _, tag := range t.config.Tags {
		if err := each(t.config.Image + ":" + tag); err != nil {
			return err
		}
	}
	return nil
}

// Stream json output to a terminal
func Stream(out io.Writer, streamer func(out io.Writer) error) error {
	outFd, isTTY := term.GetFdInfo(out)
	rpipe, wpipe := io.Pipe()
	defer rpipe.Close()

	errChan := make(chan error)

	go func() {
		err := jsonmessage.DisplayJSONMessagesStream(rpipe, out, outFd, isTTY, nil)
		errChan <- err
	}()

	err := streamer(wpipe)
	wpipe.Close()
	if err != nil {
		<-errChan
		return err
	}
	return <-errChan
}
