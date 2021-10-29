package alias

import (
	"fmt"

	"github.com/dnephin/dobi/config"
	"github.com/dnephin/dobi/tasks/task"
	"github.com/dnephin/dobi/tasks/types"
)

// GetTaskConfig returns a new TaskConfig for the action
func GetTaskConfig(name task.Name, conf *config.AliasConfig) (types.TaskConfig, error) {
	switch name.Action() {
	case task.Create:
		deps, err := conf.Dependencies()
		if err != nil {
			return nil, err
		}
		return types.NewTaskConfig(name, conf, deps, NewTask), nil
	case task.Remove:
		deps, err := RemoveDeps(conf)
		if err != nil {
			return nil, err
		}
		return types.NewTaskConfig(
			name, conf, deps, NewTask), nil
	default:
		return nil, fmt.Errorf("invalid alias action %q for task %q", name.Action(), name)
	}
}

// NewTask creates a new Task object
func NewTask(name task.Name, conf config.Resource) types.Task {
	// TODO: cleaner way to avoid this cast?
	return &Task{name: name, config: conf.(*config.AliasConfig)}
}

// RemoveDeps returns the dependencies for the remove action
func RemoveDeps(conf config.Resource) ([]task.Name, error) {
	confDeps, err := conf.Dependencies()
	if err != nil {
		return []task.Name{}, err
	}
	deps := []task.Name{}
	for i := len(confDeps) - 1; i >= 0; i-- {
		deps = append(deps, task.NewName(confDeps[i].Resource(), task.Remove))
	}
	return deps, nil
}
