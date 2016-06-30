package tasks

import (
	docker "github.com/fsouza/go-dockerclient"
)

// ExecuteContext contains all the context for task execution
type ExecuteContext struct {
	modified map[string]bool
	tasks    *TaskCollection
	client   *docker.Client
	Env      *ExecEnv
	Quiet    bool
}

func (ctx *ExecuteContext) isModified(names ...string) bool {
	for _, name := range names {
		if modified, _ := ctx.modified[name]; modified {
			return true
		}
	}
	return false
}

func (ctx *ExecuteContext) setModified(name string) {
	ctx.modified[name] = true
}

// NewExecuteContext craetes a new empty ExecuteContext
func NewExecuteContext(
	tasks *TaskCollection,
	client *docker.Client,
	execEnv *ExecEnv,
	quiet bool,
) *ExecuteContext {
	return &ExecuteContext{
		modified: make(map[string]bool),
		tasks:    tasks,
		client:   client,
		Env:      execEnv,
		Quiet:    quiet,
	}
}
