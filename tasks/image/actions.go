package image

import (
	"fmt"

	"github.com/dnephin/dobi/config"
	"github.com/dnephin/dobi/tasks/iface"
)

// GetTask returns a new task for the action
func GetTask(name, action string, conf *config.ImageConfig) (iface.Task, error) {
	switch action {
	case "", "build":
		return NewBuildTask(name, conf), nil
	case "push":
		return NewPushTask(name, conf), nil
	case "remove", "rm":
		return NewRemoveTask(name, conf), nil
	default:
		return nil, fmt.Errorf("Invalid image action %q for task %q", action, name)
	}
}
