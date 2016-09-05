package compose

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/dnephin/dobi/config"
	"github.com/dnephin/dobi/logging"
	"github.com/dnephin/dobi/tasks/common"
	"github.com/dnephin/dobi/tasks/context"
)

// Task runs a Docker Compose project
type Task struct {
	name   string
	config *config.ComposeConfig
	action action
}

type action struct {
	name string
	Run  func(ctx *context.ExecuteContext, task *Task) error
	Stop func(ctx *context.ExecuteContext, task *Task) error
}

// NewTask creates a new Task object
func NewTask(name string, conf *config.ComposeConfig, act action) *Task {
	return &Task{name: name, config: conf, action: act}
}

// Name returns the name of the task
func (t *Task) Name() common.TaskName {
	return common.NewTaskName(t.name, t.action.name)
}

func (t *Task) logger() *log.Entry {
	return logging.Log.WithFields(log.Fields{"task": t})
}

// Repr formats the task for logging
func (t *Task) Repr() string {
	return fmt.Sprintf("[compose:%s %s] %s",
		t.action.name, t.name, strings.Join(t.config.Files, ","))
}

// Run runs the action
func (t *Task) Run(ctx *context.ExecuteContext) error {
	return t.action.Run(ctx, t)
}

// Stop the task
func (t *Task) Stop(ctx *context.ExecuteContext) error {
	t.logger().Debug("Stop")
	return t.action.Stop(ctx, t)
}

// Dependencies returns the list of dependencies
func (t *Task) Dependencies() []string {
	return t.config.Dependencies()
}

// StopNothing implements the Stop interface but does nothing
func StopNothing(ctx *context.ExecuteContext, t *Task) error {
	return nil
}

func buildCommandArgs(ctx *context.ExecuteContext, conf *config.ComposeConfig) []string {
	args := []string{}
	for _, filename := range conf.Files {
		args = append(args, "-f", filename)
	}
	return append(args, "-p", conf.Project)
}

func (t *Task) execCompose(ctx *context.ExecuteContext, args ...string) error {
	if err := t.composeCommand(ctx, args...).Run(); err != nil {
		return err
	}
	t.logger().Info("Done")
	return nil
}

func (t *Task) composeCommand(ctx *context.ExecuteContext, args ...string) *exec.Cmd {
	args = append(buildCommandArgs(ctx, t.config), args...)
	cmd := exec.Command("docker-compose", args...)
	t.logger().Debugf("Args: %s", args)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd
}
