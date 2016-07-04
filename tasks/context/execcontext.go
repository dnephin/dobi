package context

import (
	docker "github.com/fsouza/go-dockerclient"

	"github.com/dnephin/dobi/config"
)

// ExecuteContext contains all the context for task execution
type ExecuteContext struct {
	modified   map[string]bool
	Resources  *config.ResourceCollection
	Client     *docker.Client
	WorkingDir string
	Env        *ExecEnv
	Quiet      bool
}

// IsModified returns true if any of the tasks named in names has been modified
// during this execution
func (ctx *ExecuteContext) IsModified(names ...string) bool {
	for _, name := range names {
		if modified, _ := ctx.modified[name]; modified {
			return true
		}
	}
	return false
}

// SetModified sets the task name as modified
func (ctx *ExecuteContext) SetModified(name string) {
	ctx.modified[name] = true
}

// NewExecuteContext craetes a new empty ExecuteContext
func NewExecuteContext(
	config *config.Config,
	client *docker.Client,
	execEnv *ExecEnv,
	quiet bool,
) *ExecuteContext {
	return &ExecuteContext{
		modified:   make(map[string]bool),
		Resources:  config.Collection,
		WorkingDir: config.WorkingDir,
		Client:     client,
		Env:        execEnv,
		Quiet:      quiet,
	}
}
