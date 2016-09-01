package alias

import (
	"fmt"

	"github.com/dnephin/dobi/config"
	"github.com/dnephin/dobi/tasks/common"
	"github.com/dnephin/dobi/tasks/iface"
)

// GetTask returns a new task for the action
func GetTask(name, act string, conf *config.AliasConfig) (iface.Task, error) {
	switch act {
	case "", "run":
		return NewTask(name, conf, action{name: "run", Dependencies: RunDeps}), nil
	case "remove", "rm":
		return NewTask(name, conf, action{name: "rm", Dependencies: RemoveDeps}), nil
	default:
		return nil, fmt.Errorf("Invalid alias action %q for task %q", act, name)
	}
}

// RunDeps returns the dependencies for the run action
func RunDeps(t *Task) []string {
	return t.config.Dependencies()
}

// RemoveDeps returns the dependencies for the remove action
func RemoveDeps(t *Task) []string {
	confDeps := t.config.Dependencies()
	deps := []string{}
	for i := len(confDeps); i > 0; i-- {
		taskname := common.ParseTaskName(confDeps[i-1])
		deps = append(deps, taskname.Resource()+":"+"rm")
	}
	return deps
}
