package image

import (
	"fmt"

	"github.com/dnephin/dobi/config"
	"github.com/dnephin/dobi/tasks/context"
	"github.com/dnephin/dobi/tasks/task"
	"github.com/dnephin/dobi/tasks/types"
)

// GetTaskConfig returns a new TaskConfig for the action
func GetTaskConfig(name task.Name, conf *config.ImageConfig) (types.TaskConfig, error) {
	if !conf.IsBuildable() && name.Action() == task.Create {
		name = task.NewName(name.Resource(), task.Pull)
	}
	imageAction, err := getImageAction(name, conf)
	if err != nil {
		return nil, err
	}

	return types.NewTaskConfig(
		name,
		conf,
		imageAction.deps,
		NewTask(imageAction.run),
	), nil
}

type runFunc func(*context.ExecuteContext, *Task, bool) (bool, error)

type imageAction struct {
	name task.Name
	run  runFunc
	deps []task.Name
}

func newImageAction(
	name task.Name, run runFunc, deps []task.Name) imageAction {
	return imageAction{name: name, run: run, deps: deps}
}

// nolint: gocyclo
func getImageAction(name task.Name, conf *config.ImageConfig) (imageAction, error) {
	switch name.Action() {
	case task.Create:
		deps, err := getDeps(conf)
		if err != nil {
			return imageAction{}, err
		}
		return newImageAction(name, RunBuild, deps), nil
	case task.Pull:
		deps, err := getDeps(conf)
		if err != nil {
			return imageAction{}, err
		}
		return newImageAction(name, RunPull, deps), nil
	case task.Push:
		deps, err := getDepsWith(conf, task.NewName(name.Resource(), task.Tag))
		if err != nil {
			return imageAction{}, err
		}
		return newImageAction(name, RunPush, deps), nil
	case task.Tag:
		deps, err := getDepsWith(conf, task.NewName(name.Resource(), task.Create))
		if err != nil {
			return imageAction{}, err
		}
		return newImageAction(name, RunTag, deps), nil
	case task.Remove:
		return newImageAction(name, RunRemove, task.NoDependencies()), nil
	default:
		return imageAction{},
			fmt.Errorf("invalid image action %q for task %q", name.Action(), name.Resource())
	}
}

func getDepsWith(conf *config.ImageConfig, newdep task.Name) ([]task.Name, error) {
	deps, err := getDeps(conf)
	if err != nil {
		return []task.Name{}, err
	}
	return append(deps, newdep), nil
}

func getDeps(conf *config.ImageConfig) ([]task.Name, error) {
	deps, err := conf.Dependencies()
	if err != nil {
		return []task.Name{}, err
	}
	return deps, nil
}

// NewTask creates a new Task object
func NewTask(runFunc runFunc) func(task.Name, config.Resource) types.Task {
	return func(name task.Name, conf config.Resource) types.Task {
		return &Task{name: name, config: conf.(*config.ImageConfig), runFunc: runFunc}
	}
}
