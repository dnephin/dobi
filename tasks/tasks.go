package tasks

import (
	"fmt"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/dnephin/dobi/config"
	"github.com/dnephin/dobi/execenv"
	"github.com/dnephin/dobi/logging"
	"github.com/dnephin/dobi/tasks/alias"
	"github.com/dnephin/dobi/tasks/client"
	"github.com/dnephin/dobi/tasks/common"
	"github.com/dnephin/dobi/tasks/compose"
	"github.com/dnephin/dobi/tasks/context"
	"github.com/dnephin/dobi/tasks/env"
	"github.com/dnephin/dobi/tasks/iface"
	"github.com/dnephin/dobi/tasks/image"
	"github.com/dnephin/dobi/tasks/job"
	"github.com/dnephin/dobi/tasks/mount"
)

// TaskCollection is a collection of Task objects
type TaskCollection struct {
	tasks []iface.TaskConfig
}

func (c *TaskCollection) add(task iface.TaskConfig) {
	c.tasks = append(c.tasks, task)
}

func (c *TaskCollection) contains(name common.TaskName) bool {
	for _, task := range c.tasks {
		if task.Name().Name() == name.Name() {
			return true
		}
	}
	return false
}

// All returns all the tasks in the dependency order
func (c *TaskCollection) All() []iface.TaskConfig {
	return c.tasks
}

func newTaskCollection() *TaskCollection {
	return &TaskCollection{}
}

func collectTasks(options RunOptions, execEnv *execenv.ExecEnv) (*TaskCollection, error) {
	return collect(options, &collectionState{
		newTaskCollection(),
		common.NewStack(),
	})
}

type collectionState struct {
	tasks     *TaskCollection
	taskStack *common.Stack
}

func collect(options RunOptions, state *collectionState) (*TaskCollection, error) {
	for _, taskname := range options.Tasks {
		taskname := common.ParseTaskName(taskname)
		resourceName := taskname.Resource()
		resource, ok := options.Config.Resources[resourceName]
		if !ok {
			return nil, fmt.Errorf("Resource %q does not exist", resourceName)
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
		state.taskStack.Pop()
	}
	return state.tasks, nil
}

// TODO: some way to make this a registry
func buildTaskConfig(name, action string, resource config.Resource) (iface.TaskConfig, error) {
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

func reversed(tasks []iface.Task) []iface.Task {
	reversed := []iface.Task{}
	for i := len(tasks) - 1; i >= 0; i-- {
		reversed = append(reversed, tasks[i])
	}
	return reversed
}

func executeTasks(ctx *context.ExecuteContext, tasks *TaskCollection) error {
	startedTasks := []iface.Task{}

	defer func() {
		logging.Log.Debug("stopping tasks")
		for _, task := range reversed(startedTasks) {
			if err := task.Stop(ctx); err != nil {
				logging.Log.Warnf("Failed to stop task %q: %s", task.Name(), err)
			}
		}
	}()

	logging.Log.Debug("executing tasks")
	for _, taskConfig := range tasks.All() {
		resource, err := taskConfig.Resource().Resolve(ctx.Env)
		if err != nil {
			return err
		}

		task := taskConfig.Task(resource)
		startedTasks = append(startedTasks, task)
		start := time.Now()
		logging.Log.WithFields(log.Fields{
			"time": start,
			"task": task,
		}).Debug("Start")

		if err := task.Run(ctx); err != nil {
			return fmt.Errorf("Failed to execute task %q: %s", task.Name(), err)
		}
		logging.Log.WithFields(log.Fields{
			"elapsed": time.Since(start),
			"task":    task,
		}).Debug("Complete")
	}
	return nil
}

// RunOptions are the options supported by Run
type RunOptions struct {
	Client client.DockerClient
	Config *config.Config
	Tasks  []string
	Quiet  bool
}

func getTaskNames(options RunOptions) []string {
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
	options.Tasks = getTaskNames(options)
	if len(options.Tasks) == 0 {
		return fmt.Errorf("No task to run, and no default task defined.")
	}

	execEnv, err := execenv.NewExecEnvFromConfig(
		options.Config.Meta.ExecID,
		options.Config.Meta.Project,
		options.Config.WorkingDir,
	)
	if err != nil {
		return err
	}

	tasks, err := collectTasks(options, execEnv)
	if err != nil {
		return err
	}

	ctx := context.NewExecuteContext(
		options.Config,
		options.Client,
		execEnv,
		options.Quiet)
	return executeTasks(ctx, tasks)
}
