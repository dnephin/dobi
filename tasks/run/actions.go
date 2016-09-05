package run

import (
	"fmt"

	"github.com/dnephin/dobi/config"
	"github.com/dnephin/dobi/tasks/iface"
)

// GetTask returns a new task for the action
func GetTask(name, action string, conf *config.JobConfig) (iface.Task, error) {
	switch action {
	case "", "run":
		return NewTask(name, conf), nil
	case "remove", "rm":
		return NewRemoveTask(name, conf), nil
	default:
		return nil, fmt.Errorf("Invalid run action %q for task %q", name, action)
	}
}
