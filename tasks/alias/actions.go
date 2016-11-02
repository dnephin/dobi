package alias

import (
	"fmt"

	"github.com/dnephin/dobi/config"
	"github.com/dnephin/dobi/tasks/task"
	"github.com/dnephin/dobi/tasks/types"
)

// GetTaskConfig returns a new TaskConfig for the action
func GetTaskConfig(name, act string, conf *config.AliasConfig) (types.TaskConfig, error) {
	switch act {
	case "", "run":
		return types.NewTaskConfig(
			task.NewDefaultName(name, "run"), conf, RunDeps(conf), NewTask), nil
	case "remove", "rm":
		return types.NewTaskConfig(
			task.NewName(name, "rm"), conf, RemoveDeps(conf), NewTask), nil
	default:
		return nil, fmt.Errorf("Invalid alias action %q for task %q", act, name)
	}
}

// NewTask creates a new Task object
func NewTask(name task.Name, conf config.Resource) types.Task {
	// TODO: cleaner way to avoid this cast?
	return &Task{name: name, config: conf.(*config.AliasConfig)}
}

// RunDeps returns the dependencies for the run action
func RunDeps(conf config.Resource) func() []string {
	return func() []string {
		return conf.Dependencies()
	}
}

// RemoveDeps returns the dependencies for the remove action
func RemoveDeps(conf config.Resource) func() []string {
	return func() []string {
		confDeps := conf.Dependencies()
		deps := []string{}
		for i := len(confDeps); i > 0; i-- {
			taskname := task.ParseName(confDeps[i-1])
			deps = append(deps, taskname.Resource()+":"+"rm")
		}
		return deps
	}
}
