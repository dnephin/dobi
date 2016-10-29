package compose

import (
	"fmt"

	"github.com/dnephin/dobi/config"
	"github.com/dnephin/dobi/tasks/common"
	"github.com/dnephin/dobi/tasks/context"
	"github.com/dnephin/dobi/tasks/iface"
)

// GetTaskConfig returns a new task for the action
func GetTaskConfig(name, action string, conf *config.ComposeConfig) (iface.TaskConfig, error) {
	act, err := getAction(action, name, conf)
	if err != nil {
		return nil, err
	}
	return iface.NewTaskConfig(act.name, conf, act.deps, NewTask(act.Run, act.Stop)), nil
}

type actionFunc func(*context.ExecuteContext, *Task) error

type action struct {
	name common.TaskName
	Run  actionFunc
	Stop actionFunc
	deps func() []string
}

func newAction(
	name common.TaskName,
	run actionFunc,
	stop actionFunc,
	deps func() []string,
) (action, error) {
	if stop == nil {
		stop = StopNothing
	}
	return action{name: name, Run: run, Stop: stop, deps: deps}, nil
}

func getAction(name string, task string, conf *config.ComposeConfig) (action, error) {
	switch name {
	case "", "up":
		return newAction(
			common.NewDefaultTaskName(task, "up"), RunUp, StopUp, deps(conf))
	case "remove", "rm", "down":
		return newAction(common.NewTaskName(task, "down"), RunDown, nil, noDeps)
	case "attach":
		return newAction(
			common.NewTaskName(task, "attach"), RunUpAttached, nil, deps(conf))
	default:
		return action{}, fmt.Errorf("Invalid compose action %q for task %q", name, task)
	}
}

// NewTask creates a new Task object
func NewTask(run actionFunc, stop actionFunc) func(common.TaskName, config.Resource) iface.Task {
	return func(name common.TaskName, res config.Resource) iface.Task {
		return &Task{
			name:   name,
			config: res.(*config.ComposeConfig),
			run:    run,
			stop:   stop,
		}
	}
}

// RunUp starts the Compose project
func RunUp(ctx *context.ExecuteContext, t *Task) error {
	t.logger().Info("project up")
	return t.execCompose(ctx, "up", "-d")
}

// StopUp stops the project
func StopUp(ctx *context.ExecuteContext, t *Task) error {
	t.logger().Info("project stop")
	return t.execCompose(ctx, "stop", "-t", t.config.StopGraceString())
}

// RunDown removes all the project resources
func RunDown(ctx *context.ExecuteContext, t *Task) error {
	t.logger().Info("project down")
	return t.execCompose(ctx, "down")
}

func deps(conf *config.ComposeConfig) func() []string {
	return func() []string {
		return conf.Dependencies()
	}
}

func noDeps() []string {
	return []string{}
}
