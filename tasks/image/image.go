package image

import (
	"fmt"
	"io"

	log "github.com/Sirupsen/logrus"
	"github.com/dnephin/dobi/config"
	"github.com/dnephin/dobi/logging"
	"github.com/dnephin/dobi/tasks/common"
	"github.com/dnephin/dobi/tasks/context"
	"github.com/docker/docker/pkg/jsonmessage"
	"github.com/docker/docker/pkg/term"
)

// Task creates a Docker image
type Task struct {
	name    common.TaskName
	config  *config.ImageConfig
	runFunc runFunc
}

// Name returns the name of the task
func (t *Task) Name() common.TaskName {
	return t.name
}

func (t *Task) logger() *log.Entry {
	return logging.Log.WithFields(log.Fields{"task": t})
}

// Repr formats the task for logging
func (t *Task) Repr() string {
	return fmt.Sprintf("[image:%s %s] %s",
		t.name.Action(), t.name.Resource(), t.config.Image)
}

// Run builds or pulls an image if it is out of date
func (t *Task) Run(ctx *context.ExecuteContext, depsModified bool) (bool, error) {
	return t.runFunc(ctx, t, depsModified)
}

// Stop the task
func (t *Task) Stop(ctx *context.ExecuteContext) error {
	return nil
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
