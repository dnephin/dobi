package alias

import (
	"fmt"

	"github.com/dnephin/dobi/config"
	"github.com/dnephin/dobi/tasks/iface"
)

// GetTask returns a new task for the action
func GetTask(name, action string, conf *config.AliasConfig) (iface.Task, error) {
	switch action {
	case "", "run":
		return NewTask(name, conf), nil
	default:
		return nil, fmt.Errorf("Invalid alias action %q for task %q", action, name)
	}
}
