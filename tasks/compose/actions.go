package compose

import (
	"fmt"

	"github.com/dnephin/dobi/config"
	"github.com/dnephin/dobi/tasks/context"
	"github.com/dnephin/dobi/tasks/task"
	"github.com/dnephin/dobi/tasks/types"
)

// GetTaskConfig returns a new task for the action
func GetTaskConfig(name task.Name, conf *config.ComposeConfig) (types.TaskConfig, error) {
	act, err := getAction(name, conf)
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
	deps []task.Name
}

func newAction(
	name task.Name,
	run actionFunc,
	stop actionFunc,
	deps []task.Name,
) action {
	if stop == nil {
		stop = StopNothing
	}
	return action{name: name, Run: run, Stop: stop, deps: deps}
}

func getAction(name task.Name, conf *config.ComposeConfig) (action, error) {
	switch name.Action() {
	case task.Create:
		deps, err := conf.Dependencies()
		if err != nil {
			return action{}, err
		}
		return newAction(name, RunUp, StopUp, deps), nil
	case task.Remove:
		return newAction(name, RunDown, nil, task.NoDependencies()), nil
	case task.Attach:
		deps, err := conf.Dependencies()
		if err != nil {
			return action{}, err
		}
		return newAction(name, RunUpAttached, nil, deps), nil
	case task.Detach:
		deps, err := conf.Dependencies()
		if err != nil {
			return action{}, err
		}
		return newAction(name, RunUp, nil, deps), nil
	default:
		return action{},
			fmt.Errorf("invalid compose action %q for task %q", name.Action(), name.Resource())
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
