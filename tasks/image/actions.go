package image

import (
	"fmt"

	"github.com/dnephin/dobi/config"
	"github.com/dnephin/dobi/tasks/context"
	"github.com/dnephin/dobi/tasks/iface"
)

// GetTask returns a new task for the action
func GetTask(name, action string, conf *config.ImageConfig) (iface.Task, error) {
	if action == "" {
		action = defaultAction(conf)
	}
	imageAction, err := getAction(action, name)
	if err != nil {
		return nil, err
	}
	return NewTask(name, conf, imageAction), nil
}

type runFunc func(ctx *context.ExecuteContext, task *Task) error

type action struct {
	name         string
	Run          runFunc
	Dependencies []string
}

func newAction(name string, run runFunc, deps []string) (action, error) {
	return action{name: name, Run: run, Dependencies: deps}, nil
}

func getAction(name string, task string) (action, error) {
	switch name {
	case "build":
		return newAction("build", RunBuild, nil)
	case "pull":
		return newAction("pull", RunPull, nil)
	case "push":
		return newAction("push", RunPush, []string{"tag"})
	case "tag":
		return newAction("tag", RunTag, []string{"build"})
	case "remove", "rm":
		return newAction("remove", RunRemove, nil)
	default:
		return action{}, fmt.Errorf("Invalid image action %q for task %q", name, task)
	}
}

func defaultAction(conf *config.ImageConfig) string {
	if conf.Dockerfile != "" || conf.Context != "" {
		return "build"
	}
	return "pull"
}
