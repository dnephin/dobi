package image

import (
	"fmt"

	"github.com/dnephin/dobi/config"
	"github.com/dnephin/dobi/tasks/iface"
)

// GetTask returns a new task for the action
func GetTask(name, action string, conf *config.ImageConfig) (iface.Task, error) {
	imageAction, err := getAction(action, name)
	if err != nil {
		return nil, err
	}
	return NewTask(name, conf, imageAction), nil
}

func getAction(name string, task string) (action, error) {
	switch name {
	case "", "build":
		return action{name: "build", Run: RunBuild}, nil
	case "push":
		return action{
			name:               "push",
			Run:                RunPush,
			ActionDependencies: []string{"tag"},
		}, nil
	case "tag":
		return action{
			name:               "tag",
			Run:                RunTag,
			ActionDependencies: []string{"build"},
		}, nil
	case "remove", "rm":
		return action{name: "remove", Run: RunRemove}, nil
	default:
		return action{}, fmt.Errorf("Invalid image action %q for task %q", name, task)
	}
}
