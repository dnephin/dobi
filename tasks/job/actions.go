package job

import (
	"fmt"

	"github.com/dnephin/dobi/config"
	"github.com/dnephin/dobi/tasks/iface"
	git "github.com/gogits/git-module"
)

// GetTask returns a new task for the action
func GetTask(name, action string, conf *config.JobConfig) (iface.Task, error) {
	repo, err := git.OpenRepository(".")
	if err != nil {
		return nil, fmt.Errorf("Failed to open git repo for action %q for task %q", name, action)
	}
	branch, err := repo.GetHEADBranch()
	if err != nil {
		return nil, fmt.Errorf("Failed to get git branch for action %q for task %q", name, action)
	}
	if conf.Branch == branch.Name || conf.Branch == "" {
		switch action {
		case "", "run":
			return NewTask(name, conf), nil
		case "remove", "rm":
			return NewRemoveTask(name, conf), nil
		default:
			return nil, fmt.Errorf("Invalid run action %q for task %q", name, action)
		}
	}

	return nil, fmt.Errorf("Task run action %q for task %q cannot run on branch %q", name, action, conf.Branch)
}
