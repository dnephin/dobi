package tasks

import (
	"fmt"
	"strings"
	"time"

	"github.com/dnephin/dobi/config"
	"github.com/dnephin/dobi/execenv"
	"github.com/dnephin/dobi/logging"
	"github.com/dnephin/dobi/tasks/alias"
	"github.com/dnephin/dobi/tasks/client"
	"github.com/dnephin/dobi/tasks/compose"
	"github.com/dnephin/dobi/tasks/context"
	"github.com/dnephin/dobi/tasks/env"
	"github.com/dnephin/dobi/tasks/image"
	"github.com/dnephin/dobi/tasks/job"
	"github.com/dnephin/dobi/tasks/mount"
	"github.com/dnephin/dobi/tasks/task"
	"github.com/dnephin/dobi/tasks/types"
	log "github.com/sirupsen/logrus"
)

// TaskCollection is a collection of Task objects
type TaskCollection struct {
	tasks []types.TaskConfig
}

func (c *TaskCollection) add(task types.TaskConfig) {
	c.tasks = append(c.tasks, task)
}

// All returns all the tasks in the dependency order
func (c *TaskCollection) All() []types.TaskConfig {
	return c.tasks
}

// Get returns the TaskConfig for the Name
func (c *TaskCollection) Get(name task.Name) types.TaskConfig {
	for _, task := range c.tasks {
		if task.Name().Equal(name) {
			return task
		}
	}
	return nil
}

func newTaskCollection() *TaskCollection {
	return &TaskCollection{}
}

func collectTasks(options RunOptions) (*TaskCollection, error) {
	return collect(options, &collectionState{
		newTaskCollection(),
		task.NewStack(),
	})
}

type collectionState struct {
	tasks     *TaskCollection
	taskStack *task.Stack
}

func collect(options RunOptions, state *collectionState) (*TaskCollection, error) {
	for _, taskname := range options.Tasks {
		taskname := task.ParseName(taskname)
		resourceName := taskname.Resource()
		resource, ok := options.Config.Resources[resourceName]
		if !ok {
			return nil, fmt.Errorf("resource %q does not exist", resourceName)
		}

		taskConfig, err := buildTaskConfig(resourceName, taskname.Action(), resource)
		if err != nil {
			return nil, err
		}

		if state.taskStack.Contains(taskConfig.Name()) {
			return nil, fmt.Errorf(
				"Invalid dependency cycle: %s", strings.Join(state.taskStack.Names(), ", "))
		}
		state.taskStack.Push(taskConfig.Name())

		options.Tasks = taskConfig.Dependencies()
		if _, err := collect(options, state); err != nil {
			return nil, err
		}
		state.tasks.add(taskConfig)
		state.taskStack.Pop() // nolint: errcheck
	}
	return state.tasks, nil
}

// TODO: some way to make this a registry
func buildTaskConfig(name, action string, resource config.Resource) (types.TaskConfig, error) {
	switch conf := resource.(type) {
	case *config.ImageConfig:
		return image.GetTaskConfig(name, action, conf)
	case *config.JobConfig:
		return job.GetTaskConfig(name, action, conf)
	case *config.MountConfig:
		return mount.GetTaskConfig(name, action, conf)
	case *config.AliasConfig:
		return alias.GetTaskConfig(name, action, conf)
	case *config.EnvConfig:
		return env.GetTaskConfig(name, action, conf)
	case *config.ComposeConfig:
		return compose.GetTaskConfig(name, action, conf)
	default:
		panic(fmt.Sprintf("Unexpected config type %T", conf))
	}
}

func reversed(tasks []types.Task) []types.Task {
	reversed := []types.Task{}
	for i := len(tasks) - 1; i >= 0; i-- {
		reversed = append(reversed, tasks[i])
	}
	return reversed
}

func executeTasks(ctx *context.ExecuteContext, tasks *TaskCollection) error {
	startedTasks := []types.Task{}

	defer func() {
		logging.Log.Debug("stopping tasks")
		for _, startedTask := range reversed(startedTasks) {
			if err := startedTask.Stop(ctx); err != nil {
				logging.Log.Warnf("Failed to stop task %q: %s", startedTask.Name(), err)
			}
		}
	}()

	logging.Log.Debug("executing tasks")
	for _, taskConfig := range tasks.All() {
		resource, err := taskConfig.Resource().Resolve(ctx.Env)
		if err != nil {
			return err
		}
		ctx.Resources.Add(taskConfig.Name().Resource(), resource)

		currentTask := taskConfig.Task(resource)
		startedTasks = append(startedTasks, currentTask)
		start := time.Now()
		logging.Log.WithFields(log.Fields{"time": start, "task": currentTask}).Debug("Start")

		depsModified := hasModifiedDeps(ctx, taskConfig.Dependencies())
		modified, err := currentTask.Run(ctx, depsModified)
		if err != nil {
			return fmt.Errorf("failed to execute task %q: %s", currentTask.Name(), err)
		}
		if modified {
			ctx.SetModified(currentTask.Name())
		}
		logging.Log.WithFields(log.Fields{
			"elapsed": time.Since(start),
			"task":    currentTask,
		}).Debug("Complete")
	}
	return nil
}

func hasModifiedDeps(ctx *context.ExecuteContext, deps []string) bool {
	for _, dep := range deps {
		taskName := task.ParseName(dep)
		if ctx.IsModified(taskName) {
			return true
		}
	}
	return false
}

// RunOptions are the options supported by Run
type RunOptions struct {
	Client    client.DockerClient
	Config    *config.Config
	Tasks     []string
	Quiet     bool
	BindMount bool
}

func getNames(options RunOptions) []string {
	if len(options.Tasks) > 0 {
		return options.Tasks
	}

	if options.Config.Meta.Default != "" {
		return []string{options.Config.Meta.Default}
	}

	return options.Tasks
}

// Run one or more tasks
func Run(options RunOptions) error {
	options.Tasks = getNames(options)
	if len(options.Tasks) == 0 {
		return fmt.Errorf("no task to run, and no default task defined")
	}

	execEnv, err := execenv.NewExecEnvFromConfig(
		options.Config.Meta.ExecID,
		options.Config.Meta.Project,
		options.Config.WorkingDir,
		options.Config.Meta.Hosted,
	)
	if err != nil {
		return err
	}

	tasks, err := collectTasks(options)
	if err != nil {
		return err
	}

	ctx := context.NewExecuteContext(
		options.Config,
		options.Client,
		execEnv,
		context.NewSettings(options.Quiet, options.BindMount))
	return executeTasks(ctx, tasks)
}
