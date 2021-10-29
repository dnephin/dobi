package job

import (
	"fmt"

	"github.com/dnephin/dobi/config"
	"github.com/dnephin/dobi/tasks/task"
	"github.com/dnephin/dobi/tasks/types"
)

// GetTaskConfig returns a new task for the action
func GetTaskConfig(name task.Name, conf *config.JobConfig) (types.TaskConfig, error) {
	switch name.Action() {
	case task.Create:
		deps, err := conf.Dependencies()
		if err != nil {
			return nil, err
		}
		return types.NewTaskConfig(name, conf, deps, newRunTask), nil
	case task.Remove:
		return types.NewTaskConfig(name, conf, nil, newRemoveTask), nil
	case task.Capture:
		deps, err := conf.Dependencies()
		if err != nil {
			return nil, err
		}
		return types.NewTaskConfig(name, conf, deps, newCaptureTask(name.CaptureVar())), nil
	}
	return nil, fmt.Errorf("invalid run action %q for task %q", name.Action(), name.Resource())
}
