package compose

import (
	"fmt"

	"github.com/dnephin/dobi/config"
	"github.com/dnephin/dobi/tasks/context"
	"github.com/dnephin/dobi/tasks/iface"
)

// GetTask returns a new task for the action
func GetTask(name, action string, conf *config.ComposeConfig) (iface.Task, error) {
	composeAction, err := getAction(action, name)
	if err != nil {
		return nil, err
	}
	return NewTask(name, conf, composeAction), nil
}

type actionFunc func(ctx *context.ExecuteContext, task *Task) error

type action struct {
	name     string
	Run      actionFunc
	Stop     actionFunc
	withDeps bool
}

func newAction(name string, run actionFunc, stop actionFunc, withDeps bool) (action, error) {
	if stop == nil {
		stop = StopNothing
	}
	return action{name: name, Run: run, Stop: stop, withDeps: withDeps}, nil
}

func getAction(name string, task string) (action, error) {
	switch name {
	case "", "up":
		return newAction("up", RunUp, StopUp, true)
	case "remove", "rm", "down":
		return newAction("down", RunDown, nil, false)
	case "attach":
		return newAction("attach", RunUpAttached, nil, false)
	default:
		return action{}, fmt.Errorf("Invalid compose action %q for task %q", name, task)
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
