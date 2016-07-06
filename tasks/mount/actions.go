package mount

import (
	"fmt"

	"github.com/dnephin/dobi/config"
	"github.com/dnephin/dobi/tasks/iface"
)

// GetTask returns a new task for the action
func GetTask(name, action string, conf *config.MountConfig) (iface.Task, error) {
	switch action {
	case "", "create":
		return NewCreateTask(name, conf), nil
	default:
		return nil, fmt.Errorf("Invalid mount action %q for task %q", action, name)
	}
}
