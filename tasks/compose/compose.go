package compose

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/dnephin/dobi/config"
	"github.com/dnephin/dobi/logging"
	"github.com/dnephin/dobi/tasks/context"
	"github.com/dnephin/dobi/tasks/task"
	log "github.com/sirupsen/logrus"
)

// Task runs a Docker Compose project
type Task struct {
	name   task.Name
	config *config.ComposeConfig
	run    actionFunc
	stop   actionFunc
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
	return fmt.Sprintf("[compose:%s %s] %s",
		t.name.Action(), t.name.Resource(), strings.Join(t.config.Files, ","))
}

// Run runs the action
func (t *Task) Run(ctx *context.ExecuteContext, _ bool) (bool, error) {
	return false, t.run(ctx, t)
}

// Stop the task
func (t *Task) Stop(ctx *context.ExecuteContext) error {
	t.logger().Debug("Stop")
	return t.stop(ctx, t)
}

// StopNothing implements the Stop interface but does nothing
func StopNothing(_ *context.ExecuteContext, _ *Task) error {
	return nil
}

func buildCommandArgs(conf *config.ComposeConfig) []string {
	args := []string{}
	for _, filename := range conf.Files {
		args = append(args, "-f", filename)
	}
	return append(args, "-p", conf.Project)
}

func (t *Task) execCompose(args ...string) error {
	if err := t.buildCommand(args...).Run(); err != nil {
		return err
	}
	t.logger().Info("Done")
	return nil
}

func (t *Task) buildCommand(args ...string) *exec.Cmd {
	args = append(buildCommandArgs(t.config), args...)
	cmd := exec.Command("docker-compose", args...)
	t.logger().Debugf("Args: %s", args)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd
}
