package context

import (
	"github.com/dnephin/dobi/config"
	"github.com/dnephin/dobi/execenv"
	"github.com/dnephin/dobi/logging"
	"github.com/dnephin/dobi/tasks/client"
	"github.com/dnephin/dobi/tasks/task"
	docker "github.com/fsouza/go-dockerclient"
)

// ExecuteContext contains all the context for task execution
type ExecuteContext struct {
	modified    map[task.Name]bool
	Resources   *ResourceCollection
	Client      client.DockerClient
	authConfigs *docker.AuthConfigurations
	WorkingDir  string
	Env         *execenv.ExecEnv
	Settings    Settings
}

// IsModified returns true if any of the tasks named in names has been modified
// during this execution
func (ctx *ExecuteContext) IsModified(names ...task.Name) bool {
	for _, name := range names {
		if modified := ctx.modified[name]; modified {
			return true
		}
	}
	return false
}

// SetModified sets the task name as modified
func (ctx *ExecuteContext) SetModified(name task.Name) {
	ctx.modified[name] = true
}

// GetAuthConfig returns the auth configuration for the repo
func (ctx *ExecuteContext) GetAuthConfig(repo string) docker.AuthConfiguration {
	if ctx.authConfigs == nil {
		return docker.AuthConfiguration{}
	}
	auth, ok := ctx.authConfigs.Configs[repo]
	if !ok {
		logging.Log.Warnf("Missing auth config for %q", repo)
	}
	return auth
}

// GetAuthConfigs returns all the authorization configs in the config file. This
// is used by build, because the repo isn't known until after the Dockerfile is
// parsed.
func (ctx *ExecuteContext) GetAuthConfigs() docker.AuthConfigurations {
	if ctx.authConfigs == nil {
		return docker.AuthConfigurations{}
	}
	return *ctx.authConfigs
}

// NewExecuteContext craetes a new empty ExecuteContext
func NewExecuteContext(
	config *config.Config,
	client client.DockerClient,
	execEnv *execenv.ExecEnv,
	settings Settings,
) *ExecuteContext {

	authConfigs, err := docker.NewAuthConfigurationsFromDockerCfg()
	if err != nil {
		logging.Log.Warnf("Failed to load auth config: %s", err)
	}

	return &ExecuteContext{
		modified:    make(map[task.Name]bool),
		Resources:   newResourceCollection(),
		WorkingDir:  config.WorkingDir,
		Client:      client,
		authConfigs: authConfigs,
		Env:         execEnv,
		Settings:    settings,
	}
}
