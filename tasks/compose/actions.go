package compose

import (
	"fmt"

	"github.com/dnephin/dobi/config"
	"github.com/dnephin/dobi/tasks/context"
	"github.com/dnephin/dobi/tasks/task"
	"github.com/dnephin/dobi/tasks/types"
)

// GetTaskConfig returns a new task for the action
func GetTaskConfig(name, action string, conf *config.ComposeConfig) (types.TaskConfig, error) {
	act, err := getAction(action, name, conf)
	if err != nil {
		return nil, err
	}
	return types.NewTaskConfig(act.name, conf, act.deps, NewTask(act.Run, act.Stop)), nil
}

type actionFunc func(*context.ExecuteContext, *Task) error

type action struct {
	name task.Name
	Run  actionFunc
	Stop actionFunc
	deps func() []string
}

func newAction(
	name task.Name,
	run actionFunc,
	stop actionFunc,
	deps func() []string,
) (action, error) {
	if stop == nil {
		stop = StopNothing
	}
	return action{name: name, Run: run, Stop: stop, deps: deps}, nil
}

func getAction(name string, resname string, conf *config.ComposeConfig) (action, error) {
	switch name {
	case "", "up":
		return newAction(
			task.NewDefaultName(resname, "up"), RunUp, StopUp, deps(conf))
	case "remove", "rm", "down":
		return newAction(task.NewName(resname, "down"), RunDown, nil, noDeps)
	case "attach":
		return newAction(
			task.NewName(resname, "attach"), RunUpAttached, nil, deps(conf))
	case "detach":
		return newAction(
			task.NewDefaultName(resname, "detach"), RunUp, nil, deps(conf))
	default:
		return action{}, fmt.Errorf("invalid compose action %q for task %q", name, resname)
	}
}

// NewTask creates a new Task object
func NewTask(run actionFunc, stop actionFunc) func(task.Name, config.Resource) types.Task {
	return func(name task.Name, res config.Resource) types.Task {
		return &Task{
			name:   name,
			config: res.(*config.ComposeConfig),
			run:    run,
			stop:   stop,
		}
	}
}

// RunUp starts the Compose project
func RunUp(_ *context.ExecuteContext, t *Task) error {
	t.logger().Info("project up")
	return t.execCompose("up", "-d")
}

// StopUp stops the project
func StopUp(_ *context.ExecuteContext, t *Task) error {
	t.logger().Info("project stop")
	return t.execCompose("stop", "-t", t.config.StopGraceString())
}

// RunDown removes all the project resources
func RunDown(_ *context.ExecuteContext, t *Task) error {
	t.logger().Info("project down")
	return t.execCompose("down")
}

func deps(conf *config.ComposeConfig) func() []string {
	return func() []string {
		return conf.Dependencies()
	}
}

func noDeps() []string {
	return []string{}
}
