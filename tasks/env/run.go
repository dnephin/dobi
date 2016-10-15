package env

import (
	"bytes"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/dnephin/dobi/config"
	"github.com/dnephin/dobi/logging"
	"github.com/dnephin/dobi/tasks/common"
	"github.com/dnephin/dobi/tasks/context"
)

// DefaultUnixSocket to connect to the docker API
const DefaultUnixSocket = "/var/run/docker.sock"

// Task is a task which runs a command in a container to produce a
// file or set of files.
type Task struct {
	name   string
	config *config.EnvConfig
}

// NewTask creates a new Task object
func NewTask(name string, conf *config.EnvConfig) *Task {
	return &Task{name: name, config: conf}
}

// Name returns the name of the task
func (t *Task) Name() common.TaskName {
	return common.NewTaskName(t.name, "run")
}

func (t *Task) logger() *log.Entry {
	return logging.Log.WithFields(log.Fields{"task": t})
}

// Repr formats the task for logging
func (t *Task) Repr() string {
	buff := &bytes.Buffer{}
	return fmt.Sprintf("[env:run %v]%v", t.name, buff.String())
}

// Run creates the host path if it doesn't already exist
func (t *Task) Run(ctx *context.ExecuteContext) error {
	log.Println(t.config.Files)
	t.logger().Info("Done")
	return nil
}



// Dependencies returns the list of dependencies
func (t *Task) Dependencies() []string {
	return t.config.Dependencies()
}

// Stop the task
func (t *Task) Stop(ctx *context.ExecuteContext) error {
	return nil
}
