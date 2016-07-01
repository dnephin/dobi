package tasks

import (
	"fmt"
	"strings"

	"github.com/dnephin/dobi/config"
	"github.com/dnephin/dobi/logging"
	"github.com/dnephin/dobi/utils/stack"
	docker "github.com/fsouza/go-dockerclient"
)

// Task is an interface implemented by all tasks
type Task interface {
	logging.LogRepresenter
	Name() string
	Prepare(*ExecuteContext) error
	Run(*ExecuteContext) error
	Stop(*ExecuteContext) error
}

// Task
type baseTask struct {
	name string
}

func (t *baseTask) Name() string {
	return t.name
}

// TaskCollection is a collection of Task objects
type TaskCollection struct {
	tasks  []Task
	mounts map[string]*MountTask
	images map[string]*ImageTask
}

func (c *TaskCollection) add(task Task) {
	c.tasks = append(c.tasks, task)
	switch typedTask := task.(type) {
	case *MountTask:
		c.mounts[task.Name()] = typedTask
	case *ImageTask:
		c.images[task.Name()] = typedTask
	}
}

func (c *TaskCollection) contains(name string) bool {
	for _, task := range c.tasks {
		if task.Name() == name {
			return true
		}
	}
	return false
}

// All returns all the tasks in the dependency order
func (c *TaskCollection) All() []Task {
	return c.tasks
}

// Reversed returns all the tasks in reversed dependency order
func (c *TaskCollection) Reversed() []Task {
	tasks := []Task{}
	for i := len(c.tasks) - 1; i >= 0; i-- {
		tasks = append(tasks, c.tasks[i])
	}
	return tasks
}

type eachMountFunc func(name string, vol *MountTask)

// EachMount iterates all the mounts in names and calls f for each
func (c *TaskCollection) EachMount(names []string, f eachMountFunc) {
	for _, name := range names {
		mount, _ := c.mounts[name]
		f(name, mount)
	}
}

func newTaskCollection() *TaskCollection {
	return &TaskCollection{
		mounts: make(map[string]*MountTask),
		images: make(map[string]*ImageTask),
	}
}

func collectTasks(options RunOptions) (*TaskCollection, error) {
	return collect(options, newTaskCollection(), stack.NewStringStack())
}

func collect(
	options RunOptions,
	tasks *TaskCollection,
	taskStack *stack.StringStack,
) (*TaskCollection, error) {
	for _, name := range options.Tasks {
		if tasks.contains(name) {
			continue
		}

		if taskStack.Contains(name) {
			return nil, fmt.Errorf(
				"Invalid dependency cycle: %s",
				strings.Join(taskStack.Items(), ", "))
		}

		resource, ok := options.Config.Resources[name]
		if !ok {
			return nil, fmt.Errorf("Resource %q does not exist", name)
		}

		task := buildTaskFromResource(taskOptions{
			name:     name,
			resource: resource,
			config:   options.Config,
		})

		taskStack.Push(name)
		options.Tasks = resource.Dependencies()
		if _, err := collect(options, tasks, taskStack); err != nil {
			return nil, err
		}
		tasks.add(task)
		taskStack.Pop()
	}
	return tasks, nil
}

type taskOptions struct {
	name     string
	client   *docker.Client
	resource config.Resource
	config   *config.Config
}

// TODO: some way to make this a registry
func buildTaskFromResource(options taskOptions) Task {
	switch conf := options.resource.(type) {
	case *config.ImageConfig:
		return NewImageTask(options, conf)
	case *config.RunConfig:
		return NewRunTask(options, conf)
	case *config.MountConfig:
		return NewMountTask(options, conf)
	case *config.AliasConfig:
		return NewAliasTask(options, conf)
	default:
		panic(fmt.Sprintf("Unexpected config type %T", conf))
	}
}

func executeTasks(ctx *ExecuteContext) error {
	logging.Log.Debug("preparing tasks")
	for _, task := range ctx.tasks.All() {
		if err := task.Prepare(ctx); err != nil {
			return fmt.Errorf("Failed to prepare task %q: %s", task.Name(), err)
		}
	}

	defer func() {
		logging.Log.Debug("stopping tasks")
		for _, task := range ctx.tasks.Reversed() {
			if err := task.Stop(ctx); err != nil {
				logging.Log.Warnf("Failed to stop task %q: %s", task.Name(), err)
			}
		}
	}()

	logging.Log.Debug("executing tasks")
	for _, task := range ctx.tasks.All() {
		if err := task.Run(ctx); err != nil {
			return fmt.Errorf("Failed to execute task %q: %s", task.Name(), err)
		}
	}
	return nil
}

// RunOptions are the options supported by Run
type RunOptions struct {
	Client *docker.Client
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

	tasks, err := collectTasks(options)
	if err != nil {
		return err
	}

	execEnv, err := NewExecEnvFromConfig(options.Config)
	if err != nil {
		return err
	}

	ctx := NewExecuteContext(tasks, options.Client, execEnv, options.Quiet)
	return executeTasks(ctx)
}
