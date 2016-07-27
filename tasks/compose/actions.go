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

func getAction(name string, task string) (action, error) {
	switch name {
	case "", "up":
		return action{name: "up", Run: RunUp, Stop: StopUp}, nil
	case "remove", "rm", "down":
		return action{name: "down", Run: StopNothing}, nil
	case "attach":
		return action{name: "attach", Run: RunUpAttached, Stop: StopNothing}, nil
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
